package kyc

import (
	"context"
	"fmt"

	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	kycmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	kycmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	reviewmgrpb "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	reviewmwcli "github.com/NpoolPlatform/review-middleware/pkg/client/review"

	"github.com/google/uuid"
)

type updateHandler struct {
	*Handler
	info *kycmwpb.Kyc
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

func (h *updateHandler) updateKyc(ctx context.Context) error {
	state := basetypes.KycState_Reviewing
	info, err := kycmwcli.UpdateKyc(ctx, &kycmwpb.KycReq{
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
	})
	if err != nil {
		return err
	}
	h.info = info
	return nil
}

func (h *Handler) UpdateKyc(ctx context.Context) (*kycmwpb.Kyc, error) {
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

	if err := handler.updateKyc(ctx); err != nil {
		return nil, err
	}
	if newReview {
		h.CreateKycReview(ctx)
	}
	return handler.info, nil
}
