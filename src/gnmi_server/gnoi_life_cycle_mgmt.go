package gnmi_server

import (
	"context"
	lcmpb "proto/gnoi/life_cycle_mgmt"
)

func (srv *Server) Activate(ctx context.Context, req *lcmpb.ActivateRequest) (*lcmpb.ActivateResponse, error) {
	return nil, nil
}
func (srv *Server) CancelUpgrade(ctx context.Context, req *lcmpb.CancelUpgradeRequest) (*lcmpb.CancelUpgradeResponse, error) {
	return nil, nil
}
func (srv *Server) GetUpgradeStatus(ctx context.Context, req *lcmpb.GetUpgradeStatusRequest) (*lcmpb.GetUpgradeStatusResponse, error) {
	return nil, nil
}
func (srv *Server) DownloadAndInstall(ctx context.Context, req *lcmpb.DownloadAndInstallRequest) (*lcmpb.DownloadAndInstallResponse, error) {
	return nil, nil
}