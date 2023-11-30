package kyc

import (
	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	reviewtypes "github.com/NpoolPlatform/message/npool/basetypes/review/v1"
	reviewmwpb "github.com/NpoolPlatform/message/npool/review/mw/v2/review"
	reviewsvcname "github.com/NpoolPlatform/review-middleware/pkg/servicename"
)

func (h *Handler) WithCreateKycReview(dispose *dtmcli.SagaDispose) {
	serviceDomain := servicename.ServiceDomain
	objectType := reviewtypes.ReviewObjectType_ObjectKyc

	req := &reviewmwpb.ReviewReq{
		EntID:      h.ReviewID,
		AppID:      h.AppID,
		ObjectID:   h.EntID,
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
