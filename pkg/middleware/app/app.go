package app

//import (
//	"context"
//	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
//	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
//	npool "github.com/NpoolPlatform/message/npool/appusermgrv2/app"
//	"github.com/google/uuid"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//)
//
//func checkAppInfo(info *appcrud.AppReq) error {
//
//	if _, err := uuid.Parse(info.GetCreatedBy()); err != nil {
//		logger.Sugar().Error("CreatedBy is invalid")
//		return status.Error(codes.InvalidArgument, "CreatedBy is invalid")
//	}
//
//	if info.Name == nil {
//		logger.Sugar().Error("Name is empty")
//		return status.Error(codes.InvalidArgument, "Name is empty")
//	}
//
//	if info.GetLogo() == "" {
//		logger.Sugar().Error("Logo is empty")
//		return status.Error(codes.InvalidArgument, "Logo is empty")
//	}
//
//	return nil
//}
//func CreateApp(ctx context.Context, in *npool.AppReq) (*npool.App,error){
//	resp,err := grpc.CreateAppV2(ctx,in)
//	if err != nil {
//		return nil,err
//	}
//	return resp, err
//}
