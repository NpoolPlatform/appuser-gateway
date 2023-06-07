package kyc

import (
	"context"
	"fmt"

	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	npool "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	reviewsvcname "github.com/NpoolPlatform/review-middleware/pkg/servicename"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) withCreateKyc(dispose *dtmcli.SagaDispose) {
	reviewID := uuid.NewString()
	if h.ReviewID == nil {
		h.ReviewID = &reviewID
	}
	id := uuid.NewString()
	if h.ID == nil {
		h.ID = &id
	}
	state := basetypes.KycState_Reviewing

	req := &npool.KycReq{
		ID:           h.ID,
		AppID:        &h.AppID,
		UserID:       h.UserID,
		DocumentType: h.DocumentType,
		IDNumber:     h.IDNumber,
		FrontImg:     h.FrontImg,
		BackImg:      h.BackImg,
		SelfieImg:    h.SelfieImg,
		EntityType:   h.EntityType,
		ReviewID:     h.ReviewID,
		State:        &state,
	}
	dispose.Add(
		reviewsvcname.ServiceDomain,
		"appuser.middleware.kyc.v1.Middleware.CreateKyc",
		"appuser.middleware.kyc.v1.Middleware.DeleteKyc",
		&npool.CreateKycRequest{
			Info: req,
		},
	)
}

func (h *Handler) CreateKyc(ctx context.Context) (*npool.Kyc, error) {
	handler := &createHandler{
		Handler: h,
	}
	if h.FrontImg == nil || h.SelfieImg == nil {
		return nil, fmt.Errorf("invalid image")
	}

	sagaDispose := dtmcli.NewSagaDispose(dtmimp.TransOptions{
		WaitResult:     true,
		RequestTimeout: handler.RequestTimeoutSeconds,
	})

	handler.withCreateKyc(sagaDispose)
	h.WithCreateKycReview(sagaDispose)

	return handler.GetKyc(ctx)
}
