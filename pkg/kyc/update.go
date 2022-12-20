package kyc

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"

	kycmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/kyc"
	kycmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"

	reviewpb "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	reviewmwcli "github.com/NpoolPlatform/review-middleware/pkg/client/review"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func UpdateKyc(ctx context.Context, in *npool.UpdateKycRequest) (info *mwpb.Kyc, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateKyc")
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

	reviewInfo, err := reviewmwcli.GetObjectReview(
		ctx,
		kycInfo.AppID,
		constant.ServiceName,
		kycInfo.ID,
		reviewpb.ReviewObjectType_ObjectKyc,
	)
	if err != nil {
		return nil, err
	}

	newReview := true
	var reviewID *string

	if reviewInfo != nil {
		switch reviewInfo.State {
		case reviewpb.ReviewState_Wait:
			reviewID = &reviewInfo.ID
			newReview = false
		case reviewpb.ReviewState_Approved:
			return nil, fmt.Errorf("not allowed")
		}
	}

	if newReview {
		_reviewID := uuid.NewString()
		reviewID = &_reviewID
	}

	// TODO: distributed transaction
	span = commontracer.TraceInvoker(span, "kyc", "manager", "UpdateKyc")

	state := kycmgrpb.KycState_Reviewing

	kyc, err := kycmgrcli.UpdateKyc(ctx, &kycmgrpb.KycReq{
		ID:           &in.KycID,
		AppID:        &in.AppID,
		UserID:       &in.UserID,
		IDNumber:     in.IDNumber,
		FrontImg:     in.FrontImg,
		BackImg:      in.BackImg,
		SelfieImg:    in.SelfieImg,
		DocumentType: in.DocumentType,
		EntityType:   in.EntityType,
		ReviewID:     reviewID,
		State:        &state,
	})
	if err != nil {
		return nil, err
	}

	if newReview {
		span = commontracer.TraceInvoker(span, "kyc", "manager", "CreateReview")

		serviceName := constant.ServiceName
		objectType := reviewpb.ReviewObjectType_ObjectKyc

		_, err = reviewmwcli.CreateReview(ctx, &reviewpb.ReviewReq{
			ID:         reviewID,
			AppID:      &in.AppID,
			ObjectID:   &kyc.ID,
			Domain:     &serviceName,
			ObjectType: &objectType,
		})
		if err != nil {
			return nil, err
		}
	}

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "GetKyc")

	return GetKyc(ctx, kyc.ID)
}
