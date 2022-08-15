//nolint:nolintlint,dupl
package banappuser

//
// import (
//	"context"
//
//	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
//	"google.golang.org/grpc/codes"
//
//	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
//	banappusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banappuser"
//	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
//	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/banappuser"
//	"github.com/google/uuid"
//	"go.opentelemetry.io/otel"
//	scodes "go.opentelemetry.io/otel/codes"
// 	"google.golang.org/grpc/status"
// )
//
// func (s *Server) DeleteBanUser(ctx context.Context,
//	in *banappuser.DeleteBanUserRequest) (*banappuser.DeleteBanUserResponse, error) {
//	var err error
//
//	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteBanAppUser")
//	defer span.End()
//	defer func() {
//		if err != nil {
//			span.SetStatus(scodes.Error, err.Error())
//			span.RecordError(err)
//		}
//	}()
//
//	commontracer.TraceID(span, in.GetID())
//
//	if _, err := uuid.Parse(in.GetID()); err != nil {
//		logger.Sugar().Errorw("DeleteBanUser", "ID", in.GetID(), "err", err)
//		return &banappuser.DeleteBanUserResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
//	}
//
//	span = commontracer.TraceInvoker(span, "banappuser", "manager", "DeleteBanAppUser")
//
//	resp, err := banappusermgrcli.DeleteBanAppUser(ctx, in.GetID())
//	if err != nil {
//		logger.Sugar().Errorw("DeleteBanUser", "err", err)
//		return &banappuser.DeleteBanUserResponse{}, status.Error(codes.Internal, err.Error())
//	}
//
//	return &banappuser.DeleteBanUserResponse{
//		Info: resp,
//	}, nil
// }
