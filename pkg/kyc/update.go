package kyc

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	kycmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/kyc"
	kycmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewmgrpb "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	reviewmgrcli "github.com/NpoolPlatform/review-service/pkg/client"
)

func UpdateKyc(ctx context.Context, in *npool.UpdateKycRequest) (info *mwpb.Kyc, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateKyc")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "GetKyc")

	kycInfo, err := kycmwcli.GetKyc(ctx, in.KycID)
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "kyc", "review-service", "GetReview")

	reviewInfo, err := reviewmgrcli.GetReview(ctx, kycInfo.ReviewID)
	if err != nil {
		return nil, err
	}

	if reviewInfo.State != reviewmgrpb.ReviewState_Wait.String() {
		span = commontracer.TraceInvoker(span, "kyc", "manager", "UpdateKyc")

		kyc, err := kycmgrcli.UpdateKyc(ctx, &kycmgrpb.KycReq{
			AppID:     &in.AppID,
			UserID:    &in.UserID,
			IDNumber:  in.IDNumber,
			FrontImg:  in.FrontImg,
			BackImg:   in.BackImg,
			SelfieImg: in.SelfieImg,
		})
		if err != nil {
			return nil, err
		}
		return kycmwcli.GetKyc(ctx, kyc.ID)
	}

	// TODO: distributed transaction
	reviewID := uuid.NewString()

	span = commontracer.TraceInvoker(span, "kyc", "manager", "UpdateKyc")

	kyc, err := kycmgrcli.UpdateKyc(ctx, &kycmgrpb.KycReq{
		AppID:     &in.AppID,
		UserID:    &in.UserID,
		IDNumber:  in.IDNumber,
		FrontImg:  in.FrontImg,
		BackImg:   in.BackImg,
		SelfieImg: in.SelfieImg,
		ReviewID:  &reviewID,
	})
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "kyc", "manager", "CreateReview")

	_, err = reviewmgrcli.CreateReview(ctx, &reviewpb.Review{
		ID:         reviewID,
		ObjectType: reviewmgrpb.ReviewObjectType_ObjectKyc.String(),
		AppID:      in.AppID,
		ObjectID:   kyc.ID,
		Domain:     constant.ServiceName,
	})
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "GetKyc")

	return kycmwcli.GetKyc(ctx, kyc.ID)
}
