package kyc

import (
	"context"

	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"

	kycmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/kyc"
	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	reviewpb "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	reviewmwcli "github.com/NpoolPlatform/review-middleware/pkg/client/review"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"github.com/google/uuid"
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
	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "CreateKyc")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	reviewID := uuid.NewString()

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "CreateKyc")

	state := kycmgrpb.KycState_Reviewing
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
		State:        &state,
	})
	if err != nil {
		return nil, err
	}
	// TODO: distributed transaction

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "CreateReview")

	serviceName := servicename.ServiceDomain
	objectType := reviewpb.ReviewObjectType_ObjectKyc

	_, err = reviewmwcli.CreateReview(ctx, &reviewpb.ReviewReq{
		ID:         &reviewID,
		AppID:      &appID,
		ObjectID:   &kycInfo.ID,
		Domain:     &serviceName,
		ObjectType: &objectType,
	})
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "GetKyc")

	return GetKyc(ctx, kycInfo.ID)
}
