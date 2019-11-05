package gnmi_server

import (
	"context"
	lcmpb "proto/gnoi/life_cycle_mgmt"
	log "github.com/golang/glog"
	"fmt"
	transutil "transl_utils"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"encoding/json"
	"strings"
)

func latestImage() string {
	
	jsresp, err:= transutil.TranslProcessAction("/sonic-image-mgmt:image-list", []byte(""))
	if err != nil {
		return ""
	}
	images := strings.Split(string(jsresp), "\n")
	return images[len(images)-1]
}

func (srv *Server) Activate(ctx context.Context, req *lcmpb.ActivateRequest) (*lcmpb.ActivateResponse, error) {
	log.V(1).Info("gNOI: LCM Activate")

	var resp lcmpb.ActivateResponse
	reqstr := fmt.Sprintf("{\"image-default:input\": {\"imagename\": \"%s\"}", latestImage())
	jsresp, err:= transutil.TranslProcessAction("/sonic-image-mgmt:image-default", []byte(reqstr))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	
	var jobj map[string]map[string]interface{}
	err = json.Unmarshal(jsresp, &jobj)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.ResponseReason = jobj["image-install:output"]["status-detail"].(string)
	resp.ResponseCode = lcmpb.ActivateResponse_ActivateRespCode(jobj["image-install:output"]["status"].(int32))
	return &resp, nil
}
func (srv *Server) CancelUpgrade(ctx context.Context, req *lcmpb.CancelUpgradeRequest) (*lcmpb.CancelUpgradeResponse, error) {
	log.V(1).Info("gNOI: LCM CancelUpgrade")
	return nil, nil
}
func (srv *Server) GetUpgradeStatus(ctx context.Context, req *lcmpb.GetUpgradeStatusRequest) (*lcmpb.GetUpgradeStatusResponse, error) {
	log.V(1).Info("gNOI: LCM GetUpgradeStatus")
	return nil, nil
}
func (srv *Server) DownloadAndInstall(ctx context.Context, req *lcmpb.DownloadAndInstallRequest) (*lcmpb.DownloadAndInstallResponse, error) {
	log.V(1).Info("gNOI: LCM DownloadAndInstall")
	var resp lcmpb.DownloadAndInstallResponse
	reqstr := fmt.Sprintf("{\"image-install:input\": {\"installer-arg-type\": \"url\", \"filename\": \"%s\"}", req.ImageUrl)
	jsresp, err:= transutil.TranslProcessAction("/sonic-image-mgmt:image-install", []byte(reqstr))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	
	var jobj map[string]map[string]interface{}
	err = json.Unmarshal(jsresp, &jobj)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.ResponseReason = jobj["image-install:output"]["status-detail"].(string)
	resp.ResponseCode = lcmpb.DownloadAndInstallResponse_DiRespCode(jobj["image-install:output"]["status"].(int32))
	return &resp, nil
}