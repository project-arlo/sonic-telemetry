
package gnmi_server

import (
	"google.golang.org/grpc/peer"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"common_utils"
	"github.com/golang/glog"
)
func ClientCertAuthenAndAuthor(ctx context.Context, admin_required bool) (context.Context, error) {
	rc, ctx := common_utils.GetContext(ctx)
	p, ok := peer.FromContext(ctx)
	if !ok {
	    return ctx, status.Error(codes.Unauthenticated, "no peer found")
	}
	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
	    return ctx, status.Error(codes.Unauthenticated, "unexpected peer transport credentials")
	}
	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
	    return ctx, status.Error(codes.Unauthenticated, "could not verify peer certificate")
	}

	var username string
	
	username = tlsAuth.State.VerifiedChains[0][0].Subject.CommonName

	if len(username) == 0 {
		return ctx, status.Error(codes.Unauthenticated, "invalid subject common name")
	}

	if DoesUserExist(username) == false {
		return ctx, status.Error(codes.Unauthenticated, "invalid subject common name")
	}
	if err := PopulateAuthStruct(username, &rc.Auth); err != nil {
		glog.Infof("[%s] Failed to retrieve authentication information; %v", rc.ID, err)
		return ctx, status.Errorf(codes.Unauthenticated, "")	
	}

	//Allow SET request only if user belong to admin group
	if admin_required && IsAdminGroup(username) == false {
		return ctx, status.Error(codes.Unauthenticated, "Not an admin user")
	}

	return ctx, nil
}