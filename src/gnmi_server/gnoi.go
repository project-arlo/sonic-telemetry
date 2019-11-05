package gnmi_server

import (
	"context"
	gnoi_system_pb "github.com/openconfig/gnoi/system"
	log "github.com/golang/glog"
	"time"
)

func (srv *Server) Reboot(context.Context, *gnoi_system_pb.RebootRequest) (*gnoi_system_pb.RebootResponse, error) {
	log.V(1).Info("gNOI: Reboot")
	return nil, nil
}
func (srv *Server) RebootStatus(context.Context, *gnoi_system_pb.RebootStatusRequest) (*gnoi_system_pb.RebootStatusResponse, error) {
	log.V(1).Info("gNOI: RebootStatus")
	return nil, nil
}
func (srv *Server) CancelReboot(context.Context, *gnoi_system_pb.CancelRebootRequest) (*gnoi_system_pb.CancelRebootResponse, error) {
	log.V(1).Info("gNOI: CancelReboot")
	return nil, nil
}
func (srv *Server) Ping(*gnoi_system_pb.PingRequest, gnoi_system_pb.System_PingServer) error {
	log.V(1).Info("gNOI: Ping")
	return nil
}
func (srv *Server) Traceroute(*gnoi_system_pb.TracerouteRequest, gnoi_system_pb.System_TracerouteServer) error {
	log.V(1).Info("gNOI: Traceroute")
	return nil
}
func (srv *Server) SetPackage(gnoi_system_pb.System_SetPackageServer) error {
	log.V(1).Info("gNOI: SetPackage")
	return nil
}
func (srv *Server) SwitchControlProcessor(context.Context, *gnoi_system_pb.SwitchControlProcessorRequest) (*gnoi_system_pb.SwitchControlProcessorResponse, error) {
	log.V(1).Info("gNOI: SwitchControlProcessor")
	return nil, nil
}
func (srv *Server) Time(ctx context.Context, req *gnoi_system_pb.TimeRequest) (*gnoi_system_pb.TimeResponse, error) {
	if srv.config.UserAuth { 
		err := PAMAuthenAndAuthor(ctx, false)
		if err != nil {
			return nil, err
		}
	}
	log.V(1).Info("gNOI: Time")
	var tm gnoi_system_pb.TimeResponse
	tm.Time = uint64(time.Now().UnixNano())
	return &tm, nil
}