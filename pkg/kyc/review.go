package kyc

import (
	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	reviewtypes "github.com/NpoolPlatform/message/npool/basetypes/review/v1"
	reviewmwpb "github.com/NpoolPlatform/message/npool/review/mw/v2/review"
	reviewsvcname "github.com/NpoolPlatform/review-middleware/pkg/servicename"
)

func (h *Handler) WithCreateKycReview(dispose *dtmcli.SagaDispose) {
	serviceDomain := servicename.ServiceDomain
	objectType := reviewtypes.ReviewObjectType_ObjectKyc

	logger.Sugar().Infow(
		"withCreateKycReview",
		"ReviewID", *h.ReviewID,
		"AppID", h.AppID,
		"ObjectID", *h.ID,
		"ServiceDomain", serviceDomain,
		"ObjectType", objectType,
	)

	req := &reviewmwpb.ReviewReq{
		ID:         h.ReviewID,
		AppID:      &h.AppID,
		ObjectID:   h.ID,
		Domain:     &serviceDomain,
		ObjectType: &objectType,
	}

	dispose.Add(
		reviewsvcname.ServiceDomain,
		"review.middleware.review.v2.Middleware/CreateReview",
		"review.middleware.review.v2.Middleware/DeleteReview",
		&reviewmwpb.CreateReviewRequest{
			Info: req,
		},
	)
}
