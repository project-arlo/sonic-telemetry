package gnmi

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"github.com/Azure/sonic-telemetry/common_utils"
	log "github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	gnoi_system_pb "github.com/openconfig/gnoi/system"
	sdc "github.com/Azure/sonic-telemetry/sonic_data_client"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
	spb "github.com/Azure/sonic-telemetry/proto/gnoi"
	"bytes"
)

var (
	supportedEncodings = []gnmipb.Encoding{gnmipb.Encoding_JSON, gnmipb.Encoding_JSON_IETF}
)

// Server manages a single gNMI Server implementation. Each client that connects
// via Subscribe or Get will receive a stream of updates based on the requested
// path. Set request is processed by server too.
type Server struct {
	s       *grpc.Server
	lis     net.Listener
	config  *Config
	cMu     sync.Mutex
	clients map[string]*Client
}
type AuthTypes map[string]bool
// Config is a collection of values for Server
type Config struct {
	// Port for the Server to listen on. If 0 or unset the Server will pick a port
	// for this Server.
	Port int64
	UserAuth AuthTypes
}

func (i AuthTypes) String() string {
	if i["none"] {
		return ""
	}
	b := new(bytes.Buffer)
	for key, value := range i {
		if value {
			fmt.Fprintf(b, "%s ", key)
		}
	}
	return b.String()
}

func (i AuthTypes) Any() bool {
	if i["none"] {
		return false
	}
	for _, value := range i {
		if value {
			return true
		}
	}
	return false
}

func (i AuthTypes) Enabled(mode string) bool {
	if i["none"] {
		return false
	}
	if value, exist := i[mode]; exist && value {
		return true
	}
	return false
}

func (i AuthTypes) Set(mode string) error {
	modes := strings.Split(mode, ",")
	for _, m := range modes {
		m = strings.Trim(m, " ")
		if m == "none" || m == "" {
			i["none"] = true
			return nil
		}

		if _, exist := i[m]; !exist {
			return fmt.Errorf("Expecting one or more of 'cert', 'password' or 'jwt'")
		}
		i[m] = true
	}
	return nil
}

func (i AuthTypes) Unset(mode string) error {
	modes := strings.Split(mode, ",")
	for _, m := range modes {
		m = strings.Trim(m, " ")
		if _, exist := i[m]; !exist {
			return fmt.Errorf("Expecting one or more of 'cert', 'password' or 'jwt'")
		}
		i[m] = false
	}
	return nil
}
// New returns an initialized Server.
func NewServer(config *Config, opts []grpc.ServerOption) (*Server, error) {
	if config == nil {
		return nil, errors.New("config not provided")
	}

	s := grpc.NewServer(opts...)
	reflection.Register(s)

	srv := &Server{
		s:       s,
		config:  config,
		clients: map[string]*Client{},
	}
	var err error
	if srv.config.Port < 0 {
		srv.config.Port = 0
	}
	srv.lis, err = net.Listen("tcp", fmt.Sprintf(":%d", srv.config.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to open listener port %d: %v", srv.config.Port, err)
	}
	gnmipb.RegisterGNMIServer(srv.s, srv)
	gnoi_system_pb.RegisterSystemServer(srv.s, srv)
	spb.RegisterSonicServiceServer(srv.s, srv)
	log.V(1).Infof("Created Server on %s", srv.Address())
	return srv, nil
}

// Serve will start the Server serving and block until closed.
func (srv *Server) Serve() error {
	s := srv.s
	if s == nil {
		return fmt.Errorf("Serve() failed: not initialized")
	}
	return srv.s.Serve(srv.lis)
}

// Address returns the port the Server is listening to.
func (srv *Server) Address() string {
	addr := srv.lis.Addr().String()
	return strings.Replace(addr, "[::]", "localhost", 1)
}

// Port returns the port the Server is listening to.
func (srv *Server) Port() int64 {
	return srv.config.Port
}

func authenticate(UserAuth AuthTypes, ctx context.Context) (context.Context,error) {
	var err error
	success := false
	rc, ctx := common_utils.GetContext(ctx)
	if !UserAuth.Any() {
		//No Auth enabled
		rc.Auth.AuthEnabled = false
		return ctx, nil
	}
	rc.Auth.AuthEnabled = true
	if UserAuth.Enabled("password") {
		ctx, err = BasicAuthenAndAuthor(ctx)
		if err == nil {
			success = true
		}
	}
	if !success && UserAuth.Enabled("jwt") {
		_,ctx,err = JwtAuthenAndAuthor(ctx)
		if err == nil {
			success = true
		}
	}
	if !success && UserAuth.Enabled("cert") {
		ctx,err = ClientCertAuthenAndAuthor(ctx)
		if err == nil {
			success = true
		}
	}

	if !success {
		return ctx,status.Error(codes.Unauthenticated, "Unauthenticated")
	} 

	return ctx,nil
}

// Subscribe implements the gNMI Subscribe RPC.
func (srv *Server) Subscribe(stream gnmipb.GNMI_SubscribeServer) error {

	ctx := stream.Context()
	ctx, err := authenticate(srv.config.UserAuth, ctx)
	if err != nil {
		return err
	}
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return grpc.Errorf(codes.InvalidArgument, "failed to get peer from ctx")
		//return fmt.Errorf("failed to get peer from ctx")
	}
	if pr.Addr == net.Addr(nil) {
		return grpc.Errorf(codes.InvalidArgument, "failed to get peer address")
	}

	/* TODO: authorize the user
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Set request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	*/

	c := NewClient(pr.Addr)

	srv.cMu.Lock()
	if oc, ok := srv.clients[c.String()]; ok {
		log.V(2).Infof("Delete duplicate client %s", oc)
		oc.Close()
		delete(srv.clients, c.String())
	}
	srv.clients[c.String()] = c
	srv.cMu.Unlock()

	err = c.Run(stream)
	srv.cMu.Lock()
	delete(srv.clients, c.String())
	srv.cMu.Unlock()

	log.Flush()
	return err
}

// checkEncodingAndModel checks whether encoding and models are supported by the server. Return error if anything is unsupported.
func (s *Server) checkEncodingAndModel(encoding gnmipb.Encoding, models []*gnmipb.ModelData) error {
	hasSupportedEncoding := false
	for _, supportedEncoding := range supportedEncodings {
		if encoding == supportedEncoding {
			hasSupportedEncoding = true
			break
		}
	}
	if !hasSupportedEncoding {
		return fmt.Errorf("unsupported encoding: %s", gnmipb.Encoding_name[int32(encoding)])
	}

	return nil
}

// Get implements the Get RPC in gNMI spec.
func (s *Server) Get(ctx context.Context, req *gnmipb.GetRequest) (*gnmipb.GetResponse, error) {
	ctx, err := authenticate(s.config.UserAuth, ctx)
	if err != nil {
		return nil, err
	}

	if req.GetType() != gnmipb.GetRequest_ALL {
		return nil, status.Errorf(codes.Unimplemented, "unsupported request type: %s", gnmipb.GetRequest_DataType_name[int32(req.GetType())])
	}

	if err = s.checkEncodingAndModel(req.GetEncoding(), req.GetUseModels()); err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	var target string
	prefix := req.GetPrefix()
	paths := req.GetPath()
	target = prefix.GetTarget()
	log.V(5).Infof("GetRequest paths: %v", paths)

	var dc sdc.Client

	if target == "OTHERS" {
		dc, err = sdc.NewNonDbClient(paths, prefix)
	} else if isTargetDb(target) == true {
		dc, err = sdc.NewDbClient(paths, prefix)
	} else {
		/* If no prefix target is specified create new Transl Data Client . */
		dc, err = sdc.NewTranslClient(prefix, paths, ctx)
	}

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	notifications := make([]*gnmipb.Notification, len(paths))
	spbValues, err := dc.Get(nil)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	for index, spbValue := range spbValues {
		update := &gnmipb.Update{
			Path: spbValue.GetPath(),
			Val:  spbValue.GetVal(),
		}

		notifications[index] = &gnmipb.Notification{
			Timestamp: spbValue.GetTimestamp(),
			Prefix:    prefix,
			Update:    []*gnmipb.Update{update},
		}
		index++
	}
	return &gnmipb.GetResponse{Notification: notifications}, nil
}

// Set method is not implemented. Refer to gnxi for examples with openconfig integration
func (s *Server) Set(ctx context.Context,req *gnmipb.SetRequest) (*gnmipb.SetResponse, error) {
	ctx, err := authenticate(s.config.UserAuth, ctx)
	if err != nil {
		return nil, err
	}
	var results []*gnmipb.UpdateResult
	

	/* Fetch the prefix. */
	prefix := req.GetPrefix()

           /* Create Transl client. */
	dc, _ := sdc.NewTranslClient(prefix, nil, ctx)

	/* DELETE */
	for _, path := range req.GetDelete() {
		log.V(2).Infof("Delete path: %v", path)

		err := dc.Set(path, nil, sdc.DELETE)

		if err != nil {
			return nil, err
		}

		res := gnmipb.UpdateResult{
						Path: path,
      						Op:   gnmipb.UpdateResult_DELETE,
 			    		          }

		/* Add to Set response results. */
 			results = append(results, &res)

	}

	/* REPLACE */
	for _, path := range req.GetReplace(){
		log.V(2).Infof("Replace path: %v ", path)

		err = dc.Set(path.GetPath(), path.GetVal(), sdc.REPLACE)

		if err != nil {
			return nil, err
		}
		res := gnmipb.UpdateResult{
						Path: path.GetPath(),
      						Op:   gnmipb.UpdateResult_REPLACE,
				                  }
		/* Add to Set response results. */
 			results = append(results, &res)
	}


	/* UPDATE */
	for _, path := range req.GetUpdate(){
		log.V(2).Infof("Update path: %v ", path)

		err = dc.Set(path.GetPath(), path.GetVal(), sdc.UPDATE)

		if err != nil {
			return nil, err
		}

		res := gnmipb.UpdateResult{
						Path: path.GetPath(),
      						Op:   gnmipb.UpdateResult_UPDATE,
 					          }
		/* Add to Set response results. */
 			results = append(results, &res)
	}



	return &gnmipb.SetResponse{
 					Prefix:   req.GetPrefix(),
		  			Response: results,
				  }, nil

}

// Capabilities method is not implemented. Refer to gnxi for examples with openconfig integration
func (s *Server) Capabilities(ctx context.Context, req *gnmipb.CapabilityRequest) (*gnmipb.CapabilityResponse, error) {
	ctx, err := authenticate(s.config.UserAuth, ctx)
	if err != nil {
		return nil, err
	}
	dc, _ := sdc.NewTranslClient(nil , nil, ctx)

		/* Fetch the client capabitlities. */
		supportedModels := dc.Capabilities()
		suppModels := make([]*gnmipb.ModelData, len(supportedModels))

		for index, model := range supportedModels {
			suppModels[index] = &gnmipb.ModelData{
						    	     	Name: model.Name, 
								Organization: model.Organization, 
								Version: model.Version,
			}
		}

	return &gnmipb.CapabilityResponse{SupportedModels: suppModels, 
				 	  SupportedEncodings: supportedEncodings,
					  GNMIVersion: "0.7.0"}, nil
}

func  isTargetDb ( target string) (bool) {
	isDbClient := false 
	dbTargetSupported := []string { "APPL_DB", "ASIC_DB" , "COUNTERS_DB", "LOGLEVEL_DB", "CONFIG_DB", "PFC_WD_DB", "FLEX_COUNTER_DB", "STATE_DB"}

	    for _, name := range  dbTargetSupported {
		    if  target ==  name {
			    isDbClient = true
				    return isDbClient
		    }
	    }

	    return isDbClient
}


