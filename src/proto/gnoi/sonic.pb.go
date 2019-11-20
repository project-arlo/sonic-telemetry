// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sonic.proto

/*
Package gnoi_sonic is a generated protocol buffer package.

It is generated from these files:
	sonic.proto

It has these top-level messages:
	TechsupportRequest
	TechsupportResponse
	SumRequest
	SumResponse
	SaveConfigRequest
	SaveConfigResponse
	ReloadConfigRequest
	ReloadConfigResponse
	LoadMgmtConfigRequest
	LoadMgmtConfigResponse
	LoadMinigraphRequest
	LoadMinigraphResponse
	JwtToken
	AuthenticateRequest
	AuthenticateResponse
	RefreshRequest
	RefreshResponse
*/
package gnoi_sonic

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TechsupportRequest struct {
	Input *TechsupportRequest_Input `protobuf:"bytes,1,opt,name=input" json:"input,omitempty"`
}

func (m *TechsupportRequest) Reset()                    { *m = TechsupportRequest{} }
func (m *TechsupportRequest) String() string            { return proto.CompactTextString(m) }
func (*TechsupportRequest) ProtoMessage()               {}
func (*TechsupportRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TechsupportRequest) GetInput() *TechsupportRequest_Input {
	if m != nil {
		return m.Input
	}
	return nil
}

type TechsupportRequest_Input struct {
	Date string `protobuf:"bytes,1,opt,name=date" json:"date,omitempty"`
}

func (m *TechsupportRequest_Input) Reset()                    { *m = TechsupportRequest_Input{} }
func (m *TechsupportRequest_Input) String() string            { return proto.CompactTextString(m) }
func (*TechsupportRequest_Input) ProtoMessage()               {}
func (*TechsupportRequest_Input) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *TechsupportRequest_Input) GetDate() string {
	if m != nil {
		return m.Date
	}
	return ""
}

type TechsupportResponse struct {
	Output *TechsupportResponse_Output `protobuf:"bytes,1,opt,name=output" json:"output,omitempty"`
}

func (m *TechsupportResponse) Reset()                    { *m = TechsupportResponse{} }
func (m *TechsupportResponse) String() string            { return proto.CompactTextString(m) }
func (*TechsupportResponse) ProtoMessage()               {}
func (*TechsupportResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TechsupportResponse) GetOutput() *TechsupportResponse_Output {
	if m != nil {
		return m.Output
	}
	return nil
}

type TechsupportResponse_Output struct {
	OutputFilename string `protobuf:"bytes,1,opt,name=output_filename,json=outputFilename" json:"output_filename,omitempty"`
}

func (m *TechsupportResponse_Output) Reset()                    { *m = TechsupportResponse_Output{} }
func (m *TechsupportResponse_Output) String() string            { return proto.CompactTextString(m) }
func (*TechsupportResponse_Output) ProtoMessage()               {}
func (*TechsupportResponse_Output) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *TechsupportResponse_Output) GetOutputFilename() string {
	if m != nil {
		return m.OutputFilename
	}
	return ""
}

type SumRequest struct {
	Input *SumRequest_Input `protobuf:"bytes,1,opt,name=input" json:"input,omitempty"`
}

func (m *SumRequest) Reset()                    { *m = SumRequest{} }
func (m *SumRequest) String() string            { return proto.CompactTextString(m) }
func (*SumRequest) ProtoMessage()               {}
func (*SumRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SumRequest) GetInput() *SumRequest_Input {
	if m != nil {
		return m.Input
	}
	return nil
}

type SumRequest_Input struct {
	Left  int32 `protobuf:"varint,1,opt,name=left" json:"left,omitempty"`
	Right int32 `protobuf:"varint,2,opt,name=right" json:"right,omitempty"`
}

func (m *SumRequest_Input) Reset()                    { *m = SumRequest_Input{} }
func (m *SumRequest_Input) String() string            { return proto.CompactTextString(m) }
func (*SumRequest_Input) ProtoMessage()               {}
func (*SumRequest_Input) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

func (m *SumRequest_Input) GetLeft() int32 {
	if m != nil {
		return m.Left
	}
	return 0
}

func (m *SumRequest_Input) GetRight() int32 {
	if m != nil {
		return m.Right
	}
	return 0
}

type SumResponse struct {
	Output *SumResponse_Output `protobuf:"bytes,1,opt,name=output" json:"output,omitempty"`
}

func (m *SumResponse) Reset()                    { *m = SumResponse{} }
func (m *SumResponse) String() string            { return proto.CompactTextString(m) }
func (*SumResponse) ProtoMessage()               {}
func (*SumResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SumResponse) GetOutput() *SumResponse_Output {
	if m != nil {
		return m.Output
	}
	return nil
}

type SumResponse_Output struct {
	Result int32 `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
}

func (m *SumResponse_Output) Reset()                    { *m = SumResponse_Output{} }
func (m *SumResponse_Output) String() string            { return proto.CompactTextString(m) }
func (*SumResponse_Output) ProtoMessage()               {}
func (*SumResponse_Output) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

func (m *SumResponse_Output) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

type SaveConfigRequest struct {
	Input *SaveConfigRequest_Input `protobuf:"bytes,1,opt,name=input" json:"input,omitempty"`
}

func (m *SaveConfigRequest) Reset()                    { *m = SaveConfigRequest{} }
func (m *SaveConfigRequest) String() string            { return proto.CompactTextString(m) }
func (*SaveConfigRequest) ProtoMessage()               {}
func (*SaveConfigRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SaveConfigRequest) GetInput() *SaveConfigRequest_Input {
	if m != nil {
		return m.Input
	}
	return nil
}

type SaveConfigRequest_Input struct {
	FilePath string `protobuf:"bytes,1,opt,name=FilePath" json:"FilePath,omitempty"`
}

func (m *SaveConfigRequest_Input) Reset()                    { *m = SaveConfigRequest_Input{} }
func (m *SaveConfigRequest_Input) String() string            { return proto.CompactTextString(m) }
func (*SaveConfigRequest_Input) ProtoMessage()               {}
func (*SaveConfigRequest_Input) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4, 0} }

func (m *SaveConfigRequest_Input) GetFilePath() string {
	if m != nil {
		return m.FilePath
	}
	return ""
}

type SaveConfigResponse struct {
	Output *SaveConfigResponse_Output `protobuf:"bytes,1,opt,name=output" json:"output,omitempty"`
}

func (m *SaveConfigResponse) Reset()                    { *m = SaveConfigResponse{} }
func (m *SaveConfigResponse) String() string            { return proto.CompactTextString(m) }
func (*SaveConfigResponse) ProtoMessage()               {}
func (*SaveConfigResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SaveConfigResponse) GetOutput() *SaveConfigResponse_Output {
	if m != nil {
		return m.Output
	}
	return nil
}

type SaveConfigResponse_Output struct {
	Status string `protobuf:"bytes,1,opt,name=Status" json:"Status,omitempty"`
}

func (m *SaveConfigResponse_Output) Reset()                    { *m = SaveConfigResponse_Output{} }
func (m *SaveConfigResponse_Output) String() string            { return proto.CompactTextString(m) }
func (*SaveConfigResponse_Output) ProtoMessage()               {}
func (*SaveConfigResponse_Output) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

func (m *SaveConfigResponse_Output) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type ReloadConfigRequest struct {
	Input *ReloadConfigRequest_Input `protobuf:"bytes,1,opt,name=input" json:"input,omitempty"`
}

func (m *ReloadConfigRequest) Reset()                    { *m = ReloadConfigRequest{} }
func (m *ReloadConfigRequest) String() string            { return proto.CompactTextString(m) }
func (*ReloadConfigRequest) ProtoMessage()               {}
func (*ReloadConfigRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ReloadConfigRequest) GetInput() *ReloadConfigRequest_Input {
	if m != nil {
		return m.Input
	}
	return nil
}

type ReloadConfigRequest_Input struct {
	FilePath string `protobuf:"bytes,1,opt,name=FilePath" json:"FilePath,omitempty"`
}

func (m *ReloadConfigRequest_Input) Reset()                    { *m = ReloadConfigRequest_Input{} }
func (m *ReloadConfigRequest_Input) String() string            { return proto.CompactTextString(m) }
func (*ReloadConfigRequest_Input) ProtoMessage()               {}
func (*ReloadConfigRequest_Input) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6, 0} }

func (m *ReloadConfigRequest_Input) GetFilePath() string {
	if m != nil {
		return m.FilePath
	}
	return ""
}

type ReloadConfigResponse struct {
	Output *ReloadConfigResponse_Output `protobuf:"bytes,1,opt,name=output" json:"output,omitempty"`
}

func (m *ReloadConfigResponse) Reset()                    { *m = ReloadConfigResponse{} }
func (m *ReloadConfigResponse) String() string            { return proto.CompactTextString(m) }
func (*ReloadConfigResponse) ProtoMessage()               {}
func (*ReloadConfigResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ReloadConfigResponse) GetOutput() *ReloadConfigResponse_Output {
	if m != nil {
		return m.Output
	}
	return nil
}

type ReloadConfigResponse_Output struct {
	Status string `protobuf:"bytes,1,opt,name=Status" json:"Status,omitempty"`
}

func (m *ReloadConfigResponse_Output) Reset()                    { *m = ReloadConfigResponse_Output{} }
func (m *ReloadConfigResponse_Output) String() string            { return proto.CompactTextString(m) }
func (*ReloadConfigResponse_Output) ProtoMessage()               {}
func (*ReloadConfigResponse_Output) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7, 0} }

func (m *ReloadConfigResponse_Output) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type LoadMgmtConfigRequest struct {
	Input *LoadMgmtConfigRequest_Input `protobuf:"bytes,1,opt,name=input" json:"input,omitempty"`
}

func (m *LoadMgmtConfigRequest) Reset()                    { *m = LoadMgmtConfigRequest{} }
func (m *LoadMgmtConfigRequest) String() string            { return proto.CompactTextString(m) }
func (*LoadMgmtConfigRequest) ProtoMessage()               {}
func (*LoadMgmtConfigRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *LoadMgmtConfigRequest) GetInput() *LoadMgmtConfigRequest_Input {
	if m != nil {
		return m.Input
	}
	return nil
}

type LoadMgmtConfigRequest_Input struct {
	FilePath string `protobuf:"bytes,1,opt,name=FilePath" json:"FilePath,omitempty"`
}

func (m *LoadMgmtConfigRequest_Input) Reset()                    { *m = LoadMgmtConfigRequest_Input{} }
func (m *LoadMgmtConfigRequest_Input) String() string            { return proto.CompactTextString(m) }
func (*LoadMgmtConfigRequest_Input) ProtoMessage()               {}
func (*LoadMgmtConfigRequest_Input) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8, 0} }

func (m *LoadMgmtConfigRequest_Input) GetFilePath() string {
	if m != nil {
		return m.FilePath
	}
	return ""
}

type LoadMgmtConfigResponse struct {
	Output *LoadMgmtConfigResponse_Output `protobuf:"bytes,1,opt,name=output" json:"output,omitempty"`
}

func (m *LoadMgmtConfigResponse) Reset()                    { *m = LoadMgmtConfigResponse{} }
func (m *LoadMgmtConfigResponse) String() string            { return proto.CompactTextString(m) }
func (*LoadMgmtConfigResponse) ProtoMessage()               {}
func (*LoadMgmtConfigResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *LoadMgmtConfigResponse) GetOutput() *LoadMgmtConfigResponse_Output {
	if m != nil {
		return m.Output
	}
	return nil
}

type LoadMgmtConfigResponse_Output struct {
	Status string `protobuf:"bytes,1,opt,name=Status" json:"Status,omitempty"`
}

func (m *LoadMgmtConfigResponse_Output) Reset()         { *m = LoadMgmtConfigResponse_Output{} }
func (m *LoadMgmtConfigResponse_Output) String() string { return proto.CompactTextString(m) }
func (*LoadMgmtConfigResponse_Output) ProtoMessage()    {}
func (*LoadMgmtConfigResponse_Output) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{9, 0}
}

func (m *LoadMgmtConfigResponse_Output) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type LoadMinigraphRequest struct {
	Input *LoadMinigraphRequest_Input `protobuf:"bytes,1,opt,name=input" json:"input,omitempty"`
}

func (m *LoadMinigraphRequest) Reset()                    { *m = LoadMinigraphRequest{} }
func (m *LoadMinigraphRequest) String() string            { return proto.CompactTextString(m) }
func (*LoadMinigraphRequest) ProtoMessage()               {}
func (*LoadMinigraphRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *LoadMinigraphRequest) GetInput() *LoadMinigraphRequest_Input {
	if m != nil {
		return m.Input
	}
	return nil
}

type LoadMinigraphRequest_Input struct {
	FilePath string `protobuf:"bytes,1,opt,name=FilePath" json:"FilePath,omitempty"`
}

func (m *LoadMinigraphRequest_Input) Reset()                    { *m = LoadMinigraphRequest_Input{} }
func (m *LoadMinigraphRequest_Input) String() string            { return proto.CompactTextString(m) }
func (*LoadMinigraphRequest_Input) ProtoMessage()               {}
func (*LoadMinigraphRequest_Input) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10, 0} }

func (m *LoadMinigraphRequest_Input) GetFilePath() string {
	if m != nil {
		return m.FilePath
	}
	return ""
}

type LoadMinigraphResponse struct {
	Output *LoadMinigraphResponse_Output `protobuf:"bytes,1,opt,name=output" json:"output,omitempty"`
}

func (m *LoadMinigraphResponse) Reset()                    { *m = LoadMinigraphResponse{} }
func (m *LoadMinigraphResponse) String() string            { return proto.CompactTextString(m) }
func (*LoadMinigraphResponse) ProtoMessage()               {}
func (*LoadMinigraphResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *LoadMinigraphResponse) GetOutput() *LoadMinigraphResponse_Output {
	if m != nil {
		return m.Output
	}
	return nil
}

type LoadMinigraphResponse_Output struct {
	Status string `protobuf:"bytes,1,opt,name=Status" json:"Status,omitempty"`
}

func (m *LoadMinigraphResponse_Output) Reset()         { *m = LoadMinigraphResponse_Output{} }
func (m *LoadMinigraphResponse_Output) String() string { return proto.CompactTextString(m) }
func (*LoadMinigraphResponse_Output) ProtoMessage()    {}
func (*LoadMinigraphResponse_Output) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{11, 0}
}

func (m *LoadMinigraphResponse_Output) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type JwtToken struct {
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token,omitempty"`
	Type        string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	ExpiresIn   int64  `protobuf:"varint,3,opt,name=expires_in,json=expiresIn" json:"expires_in,omitempty"`
}

func (m *JwtToken) Reset()                    { *m = JwtToken{} }
func (m *JwtToken) String() string            { return proto.CompactTextString(m) }
func (*JwtToken) ProtoMessage()               {}
func (*JwtToken) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *JwtToken) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *JwtToken) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *JwtToken) GetExpiresIn() int64 {
	if m != nil {
		return m.ExpiresIn
	}
	return 0
}

type AuthenticateRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *AuthenticateRequest) Reset()                    { *m = AuthenticateRequest{} }
func (m *AuthenticateRequest) String() string            { return proto.CompactTextString(m) }
func (*AuthenticateRequest) ProtoMessage()               {}
func (*AuthenticateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *AuthenticateRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *AuthenticateRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type AuthenticateResponse struct {
	Token *JwtToken `protobuf:"bytes,1,opt,name=Token" json:"Token,omitempty"`
}

func (m *AuthenticateResponse) Reset()                    { *m = AuthenticateResponse{} }
func (m *AuthenticateResponse) String() string            { return proto.CompactTextString(m) }
func (*AuthenticateResponse) ProtoMessage()               {}
func (*AuthenticateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *AuthenticateResponse) GetToken() *JwtToken {
	if m != nil {
		return m.Token
	}
	return nil
}

type RefreshRequest struct {
}

func (m *RefreshRequest) Reset()                    { *m = RefreshRequest{} }
func (m *RefreshRequest) String() string            { return proto.CompactTextString(m) }
func (*RefreshRequest) ProtoMessage()               {}
func (*RefreshRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type RefreshResponse struct {
	Token *JwtToken `protobuf:"bytes,1,opt,name=Token" json:"Token,omitempty"`
}

func (m *RefreshResponse) Reset()                    { *m = RefreshResponse{} }
func (m *RefreshResponse) String() string            { return proto.CompactTextString(m) }
func (*RefreshResponse) ProtoMessage()               {}
func (*RefreshResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *RefreshResponse) GetToken() *JwtToken {
	if m != nil {
		return m.Token
	}
	return nil
}

func init() {
	proto.RegisterType((*TechsupportRequest)(nil), "gnoi.sonic.TechsupportRequest")
	proto.RegisterType((*TechsupportRequest_Input)(nil), "gnoi.sonic.TechsupportRequest.Input")
	proto.RegisterType((*TechsupportResponse)(nil), "gnoi.sonic.TechsupportResponse")
	proto.RegisterType((*TechsupportResponse_Output)(nil), "gnoi.sonic.TechsupportResponse.Output")
	proto.RegisterType((*SumRequest)(nil), "gnoi.sonic.SumRequest")
	proto.RegisterType((*SumRequest_Input)(nil), "gnoi.sonic.SumRequest.Input")
	proto.RegisterType((*SumResponse)(nil), "gnoi.sonic.SumResponse")
	proto.RegisterType((*SumResponse_Output)(nil), "gnoi.sonic.SumResponse.Output")
	proto.RegisterType((*SaveConfigRequest)(nil), "gnoi.sonic.SaveConfigRequest")
	proto.RegisterType((*SaveConfigRequest_Input)(nil), "gnoi.sonic.SaveConfigRequest.Input")
	proto.RegisterType((*SaveConfigResponse)(nil), "gnoi.sonic.SaveConfigResponse")
	proto.RegisterType((*SaveConfigResponse_Output)(nil), "gnoi.sonic.SaveConfigResponse.Output")
	proto.RegisterType((*ReloadConfigRequest)(nil), "gnoi.sonic.ReloadConfigRequest")
	proto.RegisterType((*ReloadConfigRequest_Input)(nil), "gnoi.sonic.ReloadConfigRequest.Input")
	proto.RegisterType((*ReloadConfigResponse)(nil), "gnoi.sonic.ReloadConfigResponse")
	proto.RegisterType((*ReloadConfigResponse_Output)(nil), "gnoi.sonic.ReloadConfigResponse.Output")
	proto.RegisterType((*LoadMgmtConfigRequest)(nil), "gnoi.sonic.LoadMgmtConfigRequest")
	proto.RegisterType((*LoadMgmtConfigRequest_Input)(nil), "gnoi.sonic.LoadMgmtConfigRequest.Input")
	proto.RegisterType((*LoadMgmtConfigResponse)(nil), "gnoi.sonic.LoadMgmtConfigResponse")
	proto.RegisterType((*LoadMgmtConfigResponse_Output)(nil), "gnoi.sonic.LoadMgmtConfigResponse.Output")
	proto.RegisterType((*LoadMinigraphRequest)(nil), "gnoi.sonic.LoadMinigraphRequest")
	proto.RegisterType((*LoadMinigraphRequest_Input)(nil), "gnoi.sonic.LoadMinigraphRequest.Input")
	proto.RegisterType((*LoadMinigraphResponse)(nil), "gnoi.sonic.LoadMinigraphResponse")
	proto.RegisterType((*LoadMinigraphResponse_Output)(nil), "gnoi.sonic.LoadMinigraphResponse.Output")
	proto.RegisterType((*JwtToken)(nil), "gnoi.sonic.JwtToken")
	proto.RegisterType((*AuthenticateRequest)(nil), "gnoi.sonic.AuthenticateRequest")
	proto.RegisterType((*AuthenticateResponse)(nil), "gnoi.sonic.AuthenticateResponse")
	proto.RegisterType((*RefreshRequest)(nil), "gnoi.sonic.RefreshRequest")
	proto.RegisterType((*RefreshResponse)(nil), "gnoi.sonic.RefreshResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SonicService service

type SonicServiceClient interface {
	ShowTechsupport(ctx context.Context, in *TechsupportRequest, opts ...grpc.CallOption) (*TechsupportResponse, error)
	Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumResponse, error)
	SaveConfig(ctx context.Context, in *SaveConfigRequest, opts ...grpc.CallOption) (*SaveConfigResponse, error)
	ReloadConfig(ctx context.Context, in *ReloadConfigRequest, opts ...grpc.CallOption) (*ReloadConfigResponse, error)
	LoadMgmtConfig(ctx context.Context, in *LoadMgmtConfigRequest, opts ...grpc.CallOption) (*LoadMgmtConfigResponse, error)
	LoadMinigraph(ctx context.Context, in *LoadMinigraphRequest, opts ...grpc.CallOption) (*LoadMinigraphResponse, error)
	Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error)
	Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error)
}

type sonicServiceClient struct {
	cc *grpc.ClientConn
}

func NewSonicServiceClient(cc *grpc.ClientConn) SonicServiceClient {
	return &sonicServiceClient{cc}
}

func (c *sonicServiceClient) ShowTechsupport(ctx context.Context, in *TechsupportRequest, opts ...grpc.CallOption) (*TechsupportResponse, error) {
	out := new(TechsupportResponse)
	err := grpc.Invoke(ctx, "/gnoi.sonic.SonicService/ShowTechsupport", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sonicServiceClient) Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumResponse, error) {
	out := new(SumResponse)
	err := grpc.Invoke(ctx, "/gnoi.sonic.SonicService/Sum", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sonicServiceClient) SaveConfig(ctx context.Context, in *SaveConfigRequest, opts ...grpc.CallOption) (*SaveConfigResponse, error) {
	out := new(SaveConfigResponse)
	err := grpc.Invoke(ctx, "/gnoi.sonic.SonicService/SaveConfig", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sonicServiceClient) ReloadConfig(ctx context.Context, in *ReloadConfigRequest, opts ...grpc.CallOption) (*ReloadConfigResponse, error) {
	out := new(ReloadConfigResponse)
	err := grpc.Invoke(ctx, "/gnoi.sonic.SonicService/ReloadConfig", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sonicServiceClient) LoadMgmtConfig(ctx context.Context, in *LoadMgmtConfigRequest, opts ...grpc.CallOption) (*LoadMgmtConfigResponse, error) {
	out := new(LoadMgmtConfigResponse)
	err := grpc.Invoke(ctx, "/gnoi.sonic.SonicService/LoadMgmtConfig", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sonicServiceClient) LoadMinigraph(ctx context.Context, in *LoadMinigraphRequest, opts ...grpc.CallOption) (*LoadMinigraphResponse, error) {
	out := new(LoadMinigraphResponse)
	err := grpc.Invoke(ctx, "/gnoi.sonic.SonicService/LoadMinigraph", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sonicServiceClient) Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	out := new(AuthenticateResponse)
	err := grpc.Invoke(ctx, "/gnoi.sonic.SonicService/Authenticate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sonicServiceClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error) {
	out := new(RefreshResponse)
	err := grpc.Invoke(ctx, "/gnoi.sonic.SonicService/Refresh", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SonicService service

type SonicServiceServer interface {
	ShowTechsupport(context.Context, *TechsupportRequest) (*TechsupportResponse, error)
	Sum(context.Context, *SumRequest) (*SumResponse, error)
	SaveConfig(context.Context, *SaveConfigRequest) (*SaveConfigResponse, error)
	ReloadConfig(context.Context, *ReloadConfigRequest) (*ReloadConfigResponse, error)
	LoadMgmtConfig(context.Context, *LoadMgmtConfigRequest) (*LoadMgmtConfigResponse, error)
	LoadMinigraph(context.Context, *LoadMinigraphRequest) (*LoadMinigraphResponse, error)
	Authenticate(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error)
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
}

func RegisterSonicServiceServer(s *grpc.Server, srv SonicServiceServer) {
	s.RegisterService(&_SonicService_serviceDesc, srv)
}

func _SonicService_ShowTechsupport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TechsupportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SonicServiceServer).ShowTechsupport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.sonic.SonicService/ShowTechsupport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SonicServiceServer).ShowTechsupport(ctx, req.(*TechsupportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SonicService_Sum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SonicServiceServer).Sum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.sonic.SonicService/Sum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SonicServiceServer).Sum(ctx, req.(*SumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SonicService_SaveConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SonicServiceServer).SaveConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.sonic.SonicService/SaveConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SonicServiceServer).SaveConfig(ctx, req.(*SaveConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SonicService_ReloadConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReloadConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SonicServiceServer).ReloadConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.sonic.SonicService/ReloadConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SonicServiceServer).ReloadConfig(ctx, req.(*ReloadConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SonicService_LoadMgmtConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadMgmtConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SonicServiceServer).LoadMgmtConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.sonic.SonicService/LoadMgmtConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SonicServiceServer).LoadMgmtConfig(ctx, req.(*LoadMgmtConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SonicService_LoadMinigraph_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadMinigraphRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SonicServiceServer).LoadMinigraph(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.sonic.SonicService/LoadMinigraph",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SonicServiceServer).LoadMinigraph(ctx, req.(*LoadMinigraphRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SonicService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SonicServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.sonic.SonicService/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SonicServiceServer).Authenticate(ctx, req.(*AuthenticateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SonicService_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SonicServiceServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.sonic.SonicService/Refresh",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SonicServiceServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SonicService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gnoi.sonic.SonicService",
	HandlerType: (*SonicServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShowTechsupport",
			Handler:    _SonicService_ShowTechsupport_Handler,
		},
		{
			MethodName: "Sum",
			Handler:    _SonicService_Sum_Handler,
		},
		{
			MethodName: "SaveConfig",
			Handler:    _SonicService_SaveConfig_Handler,
		},
		{
			MethodName: "ReloadConfig",
			Handler:    _SonicService_ReloadConfig_Handler,
		},
		{
			MethodName: "LoadMgmtConfig",
			Handler:    _SonicService_LoadMgmtConfig_Handler,
		},
		{
			MethodName: "LoadMinigraph",
			Handler:    _SonicService_LoadMinigraph_Handler,
		},
		{
			MethodName: "Authenticate",
			Handler:    _SonicService_Authenticate_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _SonicService_Refresh_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sonic.proto",
}

func init() { proto.RegisterFile("sonic.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 718 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x96, 0xdf, 0x6f, 0xd3, 0x30,
	0x10, 0xc7, 0x5b, 0xb6, 0x96, 0xed, 0x3a, 0x36, 0xb8, 0x95, 0x31, 0x65, 0x6c, 0xeb, 0x3c, 0xd8,
	0x06, 0x0f, 0x95, 0x56, 0x24, 0xc4, 0xaf, 0x01, 0x03, 0x84, 0x34, 0xc4, 0x04, 0x4a, 0x2a, 0x24,
	0x9e, 0x46, 0x68, 0xdd, 0x36, 0xa2, 0x8d, 0x43, 0xec, 0xac, 0x4c, 0x88, 0x77, 0x1e, 0xf9, 0x93,
	0x51, 0x1d, 0xa7, 0x8d, 0xb3, 0xb4, 0x81, 0xbc, 0xf9, 0xce, 0x77, 0xf7, 0xcd, 0xe7, 0x6a, 0x9f,
	0x0b, 0x15, 0xce, 0x5c, 0xa7, 0x55, 0xf7, 0x7c, 0x26, 0x18, 0x42, 0xd7, 0x65, 0x4e, 0x5d, 0x7a,
	0xc8, 0x00, 0xb0, 0x49, 0x5b, 0x3d, 0x1e, 0x78, 0x1e, 0xf3, 0x85, 0x49, 0xbf, 0x07, 0x94, 0x0b,
	0x7c, 0x02, 0x25, 0xc7, 0xf5, 0x02, 0xb1, 0x5e, 0xac, 0x15, 0x0f, 0x2a, 0x8d, 0x3b, 0xf5, 0x49,
	0x46, 0xfd, 0x72, 0x78, 0xfd, 0x64, 0x14, 0x6b, 0x86, 0x29, 0xc6, 0x06, 0x94, 0xa4, 0x8d, 0x08,
	0xf3, 0x6d, 0x5b, 0x50, 0x59, 0x63, 0xd1, 0x94, 0x6b, 0xf2, 0xbb, 0x08, 0xab, 0x5a, 0x01, 0xee,
	0x31, 0x97, 0x53, 0x7c, 0x0e, 0x65, 0x16, 0x88, 0x89, 0xe2, 0xde, 0x54, 0xc5, 0x30, 0xa1, 0xfe,
	0x41, 0x46, 0x9b, 0x2a, 0xcb, 0x38, 0x84, 0x72, 0xe8, 0xc1, 0x7d, 0x58, 0x09, 0x7d, 0x67, 0x1d,
	0xa7, 0x4f, 0x5d, 0x7b, 0x10, 0x7d, 0xc0, 0x72, 0xe8, 0x7e, 0xab, 0xbc, 0x84, 0x03, 0x58, 0xc1,
	0x20, 0x22, 0x6e, 0xe8, 0xc4, 0xb7, 0xe3, 0xfa, 0x93, 0x30, 0x9d, 0xf4, 0x30, 0x46, 0xda, 0xa7,
	0x9d, 0x30, 0xb7, 0x64, 0xca, 0x35, 0x56, 0xa1, 0xe4, 0x3b, 0xdd, 0x9e, 0x58, 0xbf, 0x22, 0x9d,
	0xa1, 0x41, 0xba, 0x50, 0x91, 0xd5, 0x14, 0xf6, 0xc3, 0x04, 0xf6, 0xd6, 0x25, 0xd9, 0x74, 0xdc,
	0xda, 0x18, 0x77, 0x0d, 0xca, 0x3e, 0xe5, 0x41, 0x3f, 0x12, 0x57, 0x16, 0xe1, 0x70, 0xc3, 0xb2,
	0xcf, 0xe9, 0x6b, 0xe6, 0x76, 0x9c, 0x6e, 0x04, 0xf9, 0x58, 0x87, 0xdc, 0xd5, 0xd4, 0x92, 0xd1,
	0x3a, 0xeb, 0x6e, 0xc4, 0x6a, 0xc0, 0xc2, 0xa8, 0x85, 0x1f, 0x6d, 0xd1, 0x53, 0x8d, 0x1d, 0xdb,
	0x24, 0x00, 0x8c, 0x97, 0x51, 0x90, 0x47, 0x09, 0xc8, 0xbb, 0xd3, 0x64, 0xff, 0x81, 0xd5, 0x12,
	0xb6, 0x08, 0xb8, 0x12, 0x56, 0x16, 0x19, 0xc2, 0xaa, 0x49, 0xfb, 0xcc, 0x6e, 0xeb, 0xb4, 0x4f,
	0x75, 0x5a, 0x4d, 0x36, 0x25, 0x3e, 0x07, 0xef, 0x05, 0x54, 0xf5, 0x42, 0x8a, 0xf8, 0x45, 0x82,
	0x78, 0x7f, 0xba, 0x74, 0x5e, 0xe6, 0x9f, 0x70, 0xf3, 0x3d, 0xb3, 0xdb, 0xa7, 0xdd, 0x81, 0xd0,
	0xa9, 0x8f, 0x74, 0x6a, 0x4d, 0x3a, 0x35, 0x23, 0x07, 0xf7, 0x2f, 0x58, 0x4b, 0x96, 0x52, 0xe4,
	0xc7, 0x09, 0xf2, 0x7b, 0xb3, 0xe4, 0xf3, 0xb2, 0x5f, 0x40, 0x55, 0x96, 0x72, 0x5c, 0xa7, 0xeb,
	0xdb, 0x5e, 0x2f, 0x42, 0x7f, 0xa6, 0xa3, 0xef, 0x5d, 0xd2, 0x4e, 0x24, 0xe4, 0x20, 0x8f, 0xda,
	0x3e, 0xa9, 0xa4, 0xc0, 0x5f, 0x26, 0xc0, 0x0f, 0x66, 0x88, 0xe7, 0xe5, 0xfe, 0x02, 0x0b, 0xef,
	0x86, 0xa2, 0xc9, 0xbe, 0x51, 0x17, 0x77, 0x60, 0xc9, 0x6e, 0xb5, 0x28, 0xe7, 0x67, 0x62, 0x64,
	0xab, 0xc8, 0x4a, 0xe8, 0x0b, 0x43, 0x10, 0xe6, 0xc5, 0x85, 0x47, 0xe5, 0x00, 0x5a, 0x34, 0xe5,
	0x1a, 0x37, 0x01, 0xe8, 0x0f, 0xcf, 0xf1, 0x29, 0x3f, 0x73, 0xdc, 0xf5, 0xb9, 0x5a, 0xf1, 0x60,
	0xce, 0x5c, 0x54, 0x9e, 0x13, 0x97, 0x9c, 0xc2, 0xea, 0x71, 0x20, 0x7a, 0xd4, 0x15, 0x4e, 0xcb,
	0x16, 0x34, 0x6a, 0xac, 0x01, 0x0b, 0x01, 0xa7, 0x7e, 0x6c, 0x98, 0x8e, 0xed, 0xd1, 0x9e, 0x67,
	0x73, 0x3e, 0x64, 0x7e, 0x5b, 0x29, 0x8d, 0x6d, 0xf2, 0x0a, 0xaa, 0x7a, 0x39, 0xd5, 0xac, 0xfb,
	0x50, 0x6a, 0x8e, 0xbf, 0xba, 0xd2, 0xa8, 0xc6, 0x7b, 0x15, 0x11, 0x9a, 0x61, 0x08, 0xb9, 0x0e,
	0xcb, 0x26, 0xed, 0xf8, 0x94, 0x47, 0xbf, 0x1a, 0x39, 0x82, 0x95, 0xb1, 0xe7, 0xff, 0x0b, 0x36,
	0xfe, 0x94, 0x60, 0xc9, 0x1a, 0xed, 0x58, 0xd4, 0x3f, 0x77, 0x5a, 0x14, 0x9b, 0xb0, 0x62, 0xf5,
	0xd8, 0x30, 0xf6, 0xca, 0xe0, 0xd6, 0xec, 0x07, 0xcf, 0xd8, 0xce, 0x78, 0x9e, 0x48, 0x01, 0x1f,
	0xc1, 0x9c, 0x15, 0x0c, 0x70, 0x2d, 0xfd, 0x21, 0x31, 0x6e, 0x4d, 0x99, 0xf4, 0xa4, 0x80, 0xa7,
	0x00, 0x93, 0xa9, 0x88, 0x9b, 0x33, 0x87, 0xb4, 0xb1, 0x35, 0x7b, 0x98, 0x92, 0x02, 0x5a, 0xb0,
	0x14, 0x1f, 0x39, 0xb8, 0x9d, 0x31, 0x07, 0x8d, 0x5a, 0xd6, 0xb4, 0x22, 0x05, 0xfc, 0x0c, 0xcb,
	0xfa, 0x6d, 0xc6, 0x9d, 0xcc, 0x41, 0x63, 0x90, 0xec, 0x61, 0x40, 0x0a, 0xf8, 0x09, 0xae, 0x69,
	0xf7, 0x05, 0x6b, 0x59, 0xf7, 0xd8, 0xd8, 0xc9, 0xbc, 0x6c, 0x61, 0x1f, 0xe2, 0x87, 0x51, 0xef,
	0x43, 0xca, 0xa9, 0xd7, 0xfb, 0x90, 0x76, 0x8e, 0x49, 0x01, 0xdf, 0xc0, 0x55, 0x75, 0x16, 0xd1,
	0xd0, 0xdb, 0x16, 0x3f, 0xb2, 0xc6, 0x46, 0xea, 0x5e, 0x54, 0xe5, 0x6b, 0x59, 0xfe, 0x2f, 0x7b,
	0xf0, 0x37, 0x00, 0x00, 0xff, 0xff, 0xf7, 0xae, 0x9d, 0x03, 0xa6, 0x09, 0x00, 0x00,
}
