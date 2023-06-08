package kyc

import (
	"context"
	"fmt"

	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	appusermwsvcname "github.com/NpoolPlatform/appuser-middleware/pkg/servicename"
	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	npool "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	reviewmgrpb "github.com/NpoolPlatform/message/npool/review/mw/v2/review"
	reviewmwcli "github.com/NpoolPlatform/review-middleware/pkg/client/review"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"

	"github.com/google/uuid"
)

type updateHandler struct {
	*Handler
	info *npool.Kyc
}

func (h *updateHandler) checkReview(ctx context.Context) (bool, error) {
	info, err := reviewmwcli.GetObjectReview(
		ctx,
		h.info.AppID,
		servicename.ServiceDomain,
		*h.ID,
		reviewmgrpb.ReviewObjectType_ObjectKyc,
	)
	if err != nil {
		return false, err
	}
	if info == nil {
		return true, nil
	}

	switch info.State {
	case reviewmgrpb.ReviewState_Wait:
		h.ReviewID = &info.ID
		return false, nil
	case reviewmgrpb.ReviewState_Approved:
		return false, fmt.Errorf("not allowed")
	}

	return true, nil
}

func (h *updateHandler) withUpdateKyc(dispose *dtmcli.SagaDispose) {
	state := basetypes.KycState_Reviewing
	req := &npool.KycReq{
		ID:           h.ID,
		AppID:        &h.info.AppID,
		UserID:       &h.info.UserID,
		IDNumber:     h.IDNumber,
		FrontImg:     h.FrontImg,
		BackImg:      h.BackImg,
		SelfieImg:    h.SelfieImg,
		DocumentType: h.DocumentType,
		EntityType:   h.EntityType,
		ReviewID:     h.ReviewID,
		State:        &state,
	}

	dispose.Add(
		appusermwsvcname.ServiceDomain,
		"appuser.middleware.kyc.v1.Middleware/UpdateKyc",
		"appuser.middleware.kyc.v1.Middleware/DeleteKyc",
		&npool.UpdateKycRequest{
			Info: req,
		},
	)
}

func (h *Handler) UpdateKyc(ctx context.Context) (*npool.Kyc, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	info, err := h.GetKyc(ctx)
	if err != nil {
		return nil, err
	}

	h.AppID = info.AppID

	handler := &updateHandler{
		Handler: h,
		info:    info,
	}
	newReview, err := handler.checkReview(ctx)
	if err != nil {
		return nil, err
	}

	if newReview {
		id := uuid.NewString()
		h.ReviewID = &id
	}

	sagaDispose := dtmcli.NewSagaDispose(dtmimp.TransOptions{
		WaitResult:     true,
		RequestTimeout: handler.RequestTimeoutSeconds,
	})

	handler.withUpdateKyc(sagaDispose)
	if newReview {
		h.WithCreateKycReview(sagaDispose)
	}

	if err := dtmcli.WithSaga(ctx, sagaDispose); err != nil {
		return nil, err
	}

	return h.GetKyc(ctx)
}
