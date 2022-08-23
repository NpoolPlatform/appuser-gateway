package kyc

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	kycmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/kyc"
	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewmgrpb "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	reviewmgrcli "github.com/NpoolPlatform/review-service/pkg/client"
)

func CreateKyc(
	ctx context.Context,
	appID, userID, frontImg, selfieImg string,
	idNumber, backImg *string,
	documentType kycmgrpb.KycDocumentType,
	entityType kycmgrpb.KycEntityType,
) (
	info *mwpb.Kyc, err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateKyc")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	reviewID := uuid.NewString()

	span = commontracer.TraceInvoker(span, "kyc", "manager", "CreateKyc")

	kycInfo, err := kycmgrcli.CreateKyc(ctx, &kycmgrpb.KycReq{
		AppID:        &appID,
		UserID:       &userID,
		DocumentType: &documentType,
		IDNumber:     idNumber,
		FrontImg:     &frontImg,
		BackImg:      backImg,
		SelfieImg:    &selfieImg,
		EntityType:   &entityType,
		ReviewID:     &reviewID,
	})
	if err != nil {
		return nil, err
	}
	// TODO: distributed transaction

	span = commontracer.TraceInvoker(span, "kyc", "review-service", "CreateReview")

	_, err = reviewmgrcli.CreateReview(ctx, &reviewpb.Review{
		ID:         reviewID,
		ObjectType: reviewmgrpb.ReviewObjectType_ObjectKyc.String(),
		AppID:      appID,
		ObjectID:   kycInfo.ID,
		Domain:     constant.ServiceName,
	})
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "GetKyc")

	return GetKyc(ctx, kycInfo.ID)
}
