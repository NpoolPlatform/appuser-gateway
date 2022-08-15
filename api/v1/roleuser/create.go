package roleuser

//
// import (
//	"context"
//
//	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
//	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
//	approleusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
//	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/approleuser"
//	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
//	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/approleuser"
//	"go.opentelemetry.io/otel"
//	"go.opentelemetry.io/otel/attribute"
//	scodes "go.opentelemetry.io/otel/codes"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//)
//
// func (s *Server) CreateRoleUser(ctx context.Context, in *approleuser.CreateRoleUserRequest) (*approleuser.CreateRoleUserResponse, error) {
//	var err error
//
//	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateRoleUser")
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
//	err = validate(ctx, in.GetInfo())
//	if err != nil {
//		logger.Sugar().Errorw("CreateRoleUser", "err", err)
//		return nil, err
//	}
//
//	span = commontracer.TraceInvoker(span, "role", "manager", "CreateAppRoleUser")
//
//	resp, err := approleusermgrcli.CreateAppRoleUser(ctx, in.GetInfo())
//	if err != nil {
//		logger.Sugar().Errorw("CreateRoleUser", "err", err)
//		return &approleuser.CreateRoleUserResponse{}, status.Error(codes.Internal, err.Error())
//	}
//
//	return &approleuser.CreateRoleUserResponse{
//		Info: resp,
//	}, nil
//}
//
// func (s *Server) CreateAppRoleUser(ctx context.Context,
//	in *approleuser.CreateAppRoleUserRequest) (*approleuser.CreateAppRoleUserResponse, error) {
//	var err error
//
//	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppUserRoleUser")
//	defer span.End()
//	defer func() {
//		if err != nil {
//			span.SetStatus(scodes.Error, err.Error())
//			span.RecordError(err)
//		}
//	}()
//
//	span.SetAttributes(attribute.String("TargetAppID", in.GetTargetAppID()))
//	span.SetAttributes(attribute.String("TargetUserID", in.GetTargetUserID()))
//	span = tracer.Trace(span, in.GetInfo())
//
//	info := in.GetInfo()
//	appID := in.GetTargetAppID()
//	userID := in.GetTargetUserID()
//	info.UserID = &appID
//	info.AppID = &userID
//
//	err = validate(ctx, in.GetInfo())
//	if err != nil {
//		logger.Sugar().Errorw("CreateAppRoleUser", "err", err)
//		return nil, err
//	}
//
//	span = commontracer.TraceInvoker(span, "role", "manager", "CreateAppRoleUser")
//
//	resp, err := approleusermgrcli.CreateAppRoleUser(ctx, in.GetInfo())
//	if err != nil {
//		logger.Sugar().Errorw("CreateAppRoleUser", "err", err)
//		return &approleuser.CreateAppRoleUserResponse{}, status.Error(codes.Internal, err.Error())
//	}
//
//	return &approleuser.CreateAppRoleUserResponse{
//		Info: resp,
//	}, nil
//}
