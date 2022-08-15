package roleuser

//
// import (
//	"context"
//
//	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/approleuser"
//	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
//	"google.golang.org/grpc"
//)
//
// type Server struct {
//	approleuser.UnimplementedGatewayServer
//}
//
// func Register(server grpc.ServiceRegistrar) {
//	approleuser.RegisterGatewayServer(server, &Server{})
//}
//
// func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
//	return approleuser.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
//}
