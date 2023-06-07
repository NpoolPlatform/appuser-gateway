package kyc

import (
	"context"
	"fmt"

	kyccli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	appusermwsvcname "github.com/NpoolPlatform/appuser-middleware/pkg/servicename"
	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
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
		appusermwsvcname.ServiceDomain,
		"appuser.middleware.kyc.v1.Middleware/CreateKyc",
		"appuser.middleware.kyc.v1.Middleware/DeleteKyc",
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
	infos, _, err := kyccli.GetKycs(ctx, &npool.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: handler.AppID,
		},
		UserID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *handler.UserID,
		},
	}, 0, 10)
	if err != nil {
		return nil, err
	}
	if len(infos) > 0 {
		return nil, fmt.Errorf("kyc exist")
	}

	sagaDispose := dtmcli.NewSagaDispose(dtmimp.TransOptions{
		WaitResult:     true,
		RequestTimeout: handler.RequestTimeoutSeconds,
	})

	handler.withCreateKyc(sagaDispose)
	h.WithCreateKycReview(sagaDispose)

	if err := dtmcli.WithSaga(ctx, sagaDispose); err != nil {
		return nil, err
	}

	return handler.GetKyc(ctx)
}
