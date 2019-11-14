
package gnmi_server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	"crypto/rand"
	spb "proto/gnoi"
)

var (
	JwtRefreshInt time.Duration
	JwtValidInt   time.Duration
	hmacSampleSecret = make([]byte, 16)
)
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}


type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}



func generateJWT(username string, expire_dt time.Time) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expire_dt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(hmacSampleSecret)

	return tokenString
}
func GenerateJwtSecretKey() {
	rand.Read(hmacSampleSecret)
}

func tokenResp(username string) *spb.JwtToken {
	exp_tm := time.Now().Add(JwtValidInt)
	token := spb.JwtToken{AccessToken: generateJWT(username, exp_tm), Type: "Bearer", ExpiresIn: int64(JwtValidInt/time.Second)}
	return &token
}

func JwtAuthenAndAuthor(ctx context.Context, admin_required bool) (*spb.JwtToken, error) {
	var token spb.JwtToken
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unknown, "Invalid context")
	}


	if token_str, ok := md["access_token"]; ok {
		token.AccessToken = token_str[0]
	}else {
		return nil, status.Errorf(codes.Unauthenticated, "No JWT Token Provided")
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return &token, status.Errorf(codes.InvalidArgument, "Invalid JWT Signature")
			
		}
		return &token, status.Errorf(codes.InvalidArgument, "Bad Request")
	}
	if !tkn.Valid {
		return &token, status.Errorf(codes.InvalidArgument, "Invalid JWT Token")
	}
	return &token, nil
}

