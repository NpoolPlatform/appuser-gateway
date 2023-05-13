package kyc

import (
	"context"
	"fmt"

	kycmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	npool "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
	info *npool.Kyc
}

func (h *createHandler) createKyc(ctx context.Context) error {
	reviewID := uuid.NewString()
	if h.ReviewID == nil {
		h.ReviewID = &reviewID
	}
	id := uuid.NewString()
	if h.ID == nil {
		h.ID = &id
	}
	state := basetypes.KycState_Reviewing

	if h.FrontImg == nil || h.SelfieImg == nil {
		return fmt.Errorf("invalid image")
	}

	info, err := kycmwcli.CreateKyc(ctx, &npool.KycReq{
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
	})
	if err != nil {
		return err
	}

	h.info = info
	return nil
}

func (h *Handler) CreateKyc(ctx context.Context) (*npool.Kyc, error) {
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.createKyc(ctx); err != nil {
		return nil, err
	}
	h.CreateKycReview(ctx)
	return handler.info, nil
}
