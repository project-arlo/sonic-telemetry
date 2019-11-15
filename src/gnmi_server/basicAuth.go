package gnmi_server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/golang/glog"
	"common_utils"
)

func BasicAuthenAndAuthor(ctx context.Context, admin_required bool) (context.Context, error) {
	rc, ctx := common_utils.GetContext(ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unknown, "Invalid context")
	}
	
	var username string
	var passwd string
	if username_a, ok := md["username"]; ok {
		username = username_a[0]
	}else {
		return ctx, status.Errorf(codes.Unauthenticated, "No Username Provided")
	}

	if passwd_a, ok := md["password"]; ok {
		passwd = passwd_a[0]
	}else {
		return ctx, status.Errorf(codes.Unauthenticated, "No Password Provided")
	}
	if err := PopulateAuthStruct(username, &rc.Auth); err != nil {
		glog.Infof("[%s] Failed to retrieve authentication information; %v", rc.ID, err)
		return ctx, status.Errorf(codes.Unauthenticated, "")	
	}
	auth_success, _ := UserPwAuth(username, passwd)
	if auth_success == false {
		return ctx, status.Errorf(codes.PermissionDenied, "Invalid Password")	
	}

	//Allow SET request only if user belong to admin group
	if admin_required && IsAdminGroup(username) == false {
		return ctx, status.Errorf(codes.PermissionDenied, "Admin user required for this operation")
	}
	
	return ctx, nil
}