package kyc

import (
	"context"

	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	reviewmgrpb "github.com/NpoolPlatform/message/npool/review/mw/v2/review"
	reviewsvcname "github.com/NpoolPlatform/review-middleware/pkg/servicename"
)

func (h *Handler) CreateKycReview(ctx context.Context) {
	serviceName := servicename.ServiceDomain
	objectType := reviewmgrpb.ReviewObjectType_ObjectKyc

	logger.Sugar().Infow(
		"CreateKycReview",
		"AppID", h.AppID,
		"UserID", h.UserID,
		"ObjectID", *h.ID,
	)

	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		req := &reviewmgrpb.ReviewReq{
			ID:         h.ReviewID,
			AppID:      &h.AppID,
			ObjectID:   h.ID,
			Domain:     &serviceName,
			ObjectType: &objectType,
		}
		return publisher.Update(
			basetypes.MsgID_CreateReviewReq.String(),
			nil,
			nil,
			nil,
			req,
		)
	}); err != nil {
		logger.Sugar().Errorw(
			"CreateKycReview",
			"AppID", h.AppID,
			"UserID", h.UserID,
			"ObjectID", *h.ID,
			"Error", err,
		)
	}
}

func (h *Handler) WithCreateKycReview(dispose *dtmcli.SagaDispose) {
	serviceDomain := servicename.ServiceDomain
	objectType := reviewmgrpb.ReviewObjectType_ObjectKyc

	logger.Sugar().Infow(
		"withCreateKycReview",
		"ReviewID", *h.ReviewID,
		"AppID", h.AppID,
		"ObjectID", *h.ID,
		"ServiceDomain", serviceDomain,
		"ObjectType", objectType,
	)

	req := &reviewmgrpb.ReviewReq{
		ID:         h.ReviewID,
		AppID:      &h.AppID,
		ObjectID:   h.ID,
		Domain:     &serviceDomain,
		ObjectType: &objectType,
	}

	dispose.Add(
		reviewsvcname.ServiceDomain,
		"review.middleware.review.v2.Middleware.CreateReview",
		"review.middleware.review.v2.Middleware.DeleteReview",
		&reviewmgrpb.CreateReviewRequest{
			Info: req,
		},
	)
}
