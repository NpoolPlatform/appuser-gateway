//nolint:nolintlint,dupl
package banapp

//
// import (
//	"context"
//
//	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
//	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/banapp"
//	"google.golang.org/grpc/codes"
//
//	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
//	banappusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banapp"
//	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
//	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/banapp"
//	banappcrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banapp"
//	"github.com/google/uuid"
//	"go.opentelemetry.io/otel"
//	scodes "go.opentelemetry.io/otel/codes"
//	"google.golang.org/grpc/status"
// )
//
// func (s *Server) UpdateBanApp(ctx context.Context, in *banapp.UpdateBanAppRequest) (*banapp.UpdateBanAppResponse, error) {
//	var err error
//
//	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateBanApp")
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
//		logger.Sugar().Errorw("UpdateBanApp", "err", err)
//		return nil, err
//	}
//
//	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
//		logger.Sugar().Errorw("UpdateBanApp", "ID", in.GetInfo().GetID(), "err", err)
//		return &banapp.UpdateBanAppResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
//	}
//
//	if in.GetInfo().GetMessage() == "" {
//		logger.Sugar().Errorw("UpdateBanApp", "Message", in.GetInfo().GetMessage(), "err", "Message is empty")
//		return &banapp.UpdateBanAppResponse{}, status.Error(codes.InvalidArgument, "Message is empty")
//	}
//
//	span = commontracer.TraceInvoker(span, "banapp", "manager", "UpdateBanApp")
//
//	resp, err := banappusermgrcli.UpdateBanApp(ctx, &banappcrud.BanAppReq{
//		Message: in.GetInfo().Message,
//	})
//	if err != nil {
//		logger.Sugar().Errorw("UpdateBanApp", "err", err)
//		return &banapp.UpdateBanAppResponse{}, status.Error(codes.Internal, err.Error())
//	}
//
//	return &banapp.UpdateBanAppResponse{
//		Info: resp,
//	}, nil
// }
