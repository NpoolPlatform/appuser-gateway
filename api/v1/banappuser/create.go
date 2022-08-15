//nolint:nolintlint,dupl
package banappuser

//
// import (
//	"context"
//
//	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
//	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/banappuser"
//	"google.golang.org/grpc/codes"
//
//	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
//	banappusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banappuser"
//	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
//	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/banappuser"
//	"go.opentelemetry.io/otel"
//	scodes "go.opentelemetry.io/otel/codes"
//	"google.golang.org/grpc/status"
// )
//
// func (s *Server) CreateBanUser(ctx context.Context,
//	in *banappuser.CreateBanUserRequest) (*banappuser.CreateBanUserResponse, error) {
//	var err error
//
//	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateBanAppUser")
//	defer span.End()
//	defer func() {
//		if err != nil {
//			span.SetStatus(scodes.Error, err.Error())
//			span.RecordError(err)
//		}
//	}()
//
//	span = tracer.Trace(span, in.GetInfo())
//
//	err = validate(in.GetInfo())
//	if err != nil {
//		return nil, err
//	}
//
//	span = commontracer.TraceInvoker(span, "banappuser", "middleware", "CreateBanAppUser")
//
//	resp, err := banappusermgrcli.CreateBanAppUser(ctx, in.GetInfo())
//	if err != nil {
//		logger.Sugar().Errorw("CreateBanUser", "err", err)
//		return &banappuser.CreateBanUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
//	}
//
//	return &banappuser.CreateBanUserResponse{
//		Info: resp,
//	}, nil
// }
//
// func (s *Server) CreateAppBanUser(ctx context.Context,
//	in *banappuser.CreateAppBanUserRequest) (*banappuser.CreateAppBanUserResponse, error) {
//	var err error
//
//	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppBanUser")
//	defer span.End()
//	defer func() {
//		if err != nil {
//			span.SetStatus(scodes.Error, err.Error())
//			span.RecordError(err)
//		}
//	}()
//
//	targetAppID := in.GetTargetAppID()
//	targetUserID := in.GetTargetUserID()
//	in.Info.AppID = &targetAppID
//	in.Info.UserID = &targetUserID
//
//	span = tracer.Trace(span, in.GetInfo())
//
//	err = validate(in.GetInfo())
//	if err != nil {
//		logger.Sugar().Errorw("CreateAppBanUser", "err", err)
//		return nil, err
//	}
//
//	span = commontracer.TraceInvoker(span, "banappuser", "middleware", "CreateBanAppUser")
//
//	resp, err := banappusermgrcli.CreateBanAppUser(ctx, in.GetInfo())
//	if err != nil {
//		logger.Sugar().Errorw("CreateAppBanUser", "err", err)
//		return &banappuser.CreateAppBanUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
//	}
//
//	return &banappuser.CreateAppBanUserResponse{
//		Info: resp,
//	}, nil
// }
