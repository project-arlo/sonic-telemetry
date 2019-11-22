package gnmi_server

import (
	"context"
	gnoi_system_pb "github.com/openconfig/gnoi/system"
	log "github.com/golang/glog"
	"time"
	spb "proto/gnoi"
	transutil "transl_utils"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"encoding/json"
	"strings"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
)

func (srv *Server) Reboot(ctx context.Context, req *gnoi_system_pb.RebootRequest) (*gnoi_system_pb.RebootResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Reboot")
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (srv *Server) RebootStatus(ctx context.Context, req *gnoi_system_pb.RebootStatusRequest) (*gnoi_system_pb.RebootStatusResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: RebootStatus")
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (srv *Server) CancelReboot(ctx context.Context, req *gnoi_system_pb.CancelRebootRequest) (*gnoi_system_pb.CancelRebootResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: CancelReboot")
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (srv *Server) Ping(req *gnoi_system_pb.PingRequest, rs gnoi_system_pb.System_PingServer) error {
	ctx := rs.Context()
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return err
	}
	log.V(1).Info("gNOI: Ping")
	return status.Errorf(codes.Unimplemented, "")
}
func (srv *Server) Traceroute(req *gnoi_system_pb.TracerouteRequest, rs gnoi_system_pb.System_TracerouteServer) error {
	ctx := rs.Context()
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return err
	}
	log.V(1).Info("gNOI: Traceroute")
	return status.Errorf(codes.Unimplemented, "")
}
func (srv *Server) SetPackage(rs gnoi_system_pb.System_SetPackageServer) error {
	ctx := rs.Context()
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return err
	}
	log.V(1).Info("gNOI: SetPackage")
	return status.Errorf(codes.Unimplemented, "")
}
func (srv *Server) SwitchControlProcessor(ctx context.Context, req *gnoi_system_pb.SwitchControlProcessorRequest) (*gnoi_system_pb.SwitchControlProcessorResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: SwitchControlProcessor")
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (srv *Server) Time(ctx context.Context, req *gnoi_system_pb.TimeRequest) (*gnoi_system_pb.TimeResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Time")
	var tm gnoi_system_pb.TimeResponse
	tm.Time = uint64(time.Now().UnixNano())
	return &tm, nil
}

func (srv *Server) CopyConfig(ctx context.Context, req *spb.CopyConfigRequest) (*spb.CopyConfigResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic CopyConfig")
	var jobj map[string]map[string]interface{}
	resp := &spb.CopyConfigResponse{
		Output: &spb.SonicOutput {

		},
	}
	reqstr := fmt.Sprintf("{\"sonic-config-mgmt:input\": {\"source\": \"%s\", \"overwrite\": %t, \"destination\": \"%s\"}}", req.Input.Source, req.Input.Overwrite, req.Input.Destination)
	jsresp, err:= transutil.TranslProcessAction("/sonic-config-mgmt:copy", []byte(reqstr), ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	jsresp = []byte(strings.Replace(string(jsresp), "sonic-config-mgmt:output", "output", 1))
	err = json.Unmarshal(jsresp, &jobj)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.Output.Status = jobj["output"]["status"].(int32)
	resp.Output.StatusDetail = jobj["output"]["status-detail"].(string)
	
	return resp, nil
}

func (srv *Server) ShowTechsupport(ctx context.Context, req *spb.TechsupportRequest) (*spb.TechsupportResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic ShowTechsupport")
	var jobj map[string]map[string]string
	resp := &spb.TechsupportResponse{
		Output: &spb.TechsupportResponse_Output {

		},
	}
	reqstr := fmt.Sprintf("{\"sonic-show-techsupport-info:input\": {\"date\": \"%s\"}}", req.Input.Date)
	jsresp, err:= transutil.TranslProcessAction("/sonic-show-techsupport:sonic-show-techsupport-info", []byte(reqstr), ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	jsresp = []byte(strings.Replace(string(jsresp), "sonic-show-techsupport:output", "output", 1))
	err = json.Unmarshal(jsresp, &jobj)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.Output.OutputFilename = jobj["output"]["output-filename"]
	
	return resp, nil
}

func (srv *Server) ImageInstall(ctx context.Context, req *spb.ImageInstallRequest) (*spb.ImageInstallResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic ImageInstall")
	var jobj map[string]map[string]interface{}
	resp := &spb.ImageInstallResponse{
		Output: &spb.SonicOutput {

		},
	}
	reqstr := fmt.Sprintf("{\"sonic-image-mgmt:input\": {\"imagename\": \"%s\"}}", req.Input.Imagename)
	jsresp, err:= transutil.TranslProcessAction("/sonic-image-mgmt:image-install", []byte(reqstr), ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	jsresp = []byte(strings.Replace(string(jsresp), "sonic-image-mgmt:output", "output", 1))
	err = json.Unmarshal(jsresp, &jobj)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.Output.Status = jobj["output"]["status"].(int32)
	resp.Output.StatusDetail = jobj["output"]["status-detail"].(string)
	
	return resp, nil
}

func (srv *Server) ImageRemove(ctx context.Context, req *spb.ImageRemoveRequest) (*spb.ImageRemoveResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic ImageRemove")
	var jobj map[string]map[string]interface{}
	resp := &spb.ImageRemoveResponse{
		Output: &spb.SonicOutput {

		},
	}
	reqstr := fmt.Sprintf("{\"sonic-image-mgmt:input\": {\"imagename\": \"%s\"}}", req.Input.Imagename)
	jsresp, err:= transutil.TranslProcessAction("/sonic-image-mgmt:image-remove", []byte(reqstr), ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	jsresp = []byte(strings.Replace(string(jsresp), "sonic-image-mgmt:output", "output", 1))
	err = json.Unmarshal(jsresp, &jobj)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.Output.Status = jobj["output"]["status"].(int32)
	resp.Output.StatusDetail = jobj["output"]["status-detail"].(string)
	
	return resp, nil
}

func (srv *Server) ImageDefault(ctx context.Context, req *spb.ImageDefaultRequest) (*spb.ImageDefaultResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic ImageDefault")
	var jobj map[string]map[string]interface{}
	resp := &spb.ImageDefaultResponse{
		Output: &spb.SonicOutput {

		},
	}
	reqstr := fmt.Sprintf("{\"sonic-image-mgmt:input\": {\"imagename\": \"%s\"}}", req.Input.Imagename)
	jsresp, err:= transutil.TranslProcessAction("/sonic-image-mgmt:image-default", []byte(reqstr), ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	jsresp = []byte(strings.Replace(string(jsresp), "sonic-image-mgmt:output", "output", 1))
	err = json.Unmarshal(jsresp, &jobj)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.Output.Status = jobj["output"]["status"].(int32)
	resp.Output.StatusDetail = jobj["output"]["status-detail"].(string)
	
	return resp, nil
}



func (srv *Server) Sum(ctx context.Context, req *spb.SumRequest) (*spb.SumResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic Sum")
	var resp spb.SumResponse
	reqstr := fmt.Sprintf("{\"sonic-tests:input\": {\"left\": %d, \"right\": %d}}", req.Input.Left, req.Input.Right)
	jsresp, err:= transutil.TranslProcessAction("/sonic-tests:sum", []byte(reqstr), ctx)
	
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	jsresp = []byte(strings.Replace(string(jsresp), "sonic-tests:output", "output", 1))
	err = json.Unmarshal(jsresp, &resp)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &resp, nil
}



func (srv *Server) Authenticate(ctx context.Context, req *spb.AuthenticateRequest) (*spb.AuthenticateResponse, error) {
	// Can't enforce normal authentication here.. maybe only enforce client cert auth if enabled?
	// ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	// if err != nil {
	// 	return nil, err
	// }
	log.V(1).Info("gNOI: Sonic Authenticate")


	if !srv.config.UserAuth.Enabled("jwt") {
		return nil, status.Errorf(codes.Unimplemented, "")
	}
	auth_success, _ := UserPwAuth(req.Username, req.Password)
	if  auth_success {
		return &spb.AuthenticateResponse{Token: tokenResp(req.Username)}, nil
	}
	return nil, status.Errorf(codes.PermissionDenied, "Invalid Username or Password")

}
func (srv *Server) Refresh(ctx context.Context, req *spb.RefreshRequest) (*spb.RefreshResponse, error) {
	ctx,err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic Refresh")

	if !srv.config.UserAuth.Enabled("jwt") {
		return nil, status.Errorf(codes.Unimplemented, "")
	}

	token, ctx, err := JwtAuthenAndAuthor(ctx, false)
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	jwt.ParseWithClaims(token.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > JwtRefreshInt {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid JWT Token")
	}
	
	return &spb.RefreshResponse{Token: tokenResp(claims.Username)}, nil

}