package kyc

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	kycmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	kycmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	reviewmgrpb "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	reviewmwcli "github.com/NpoolPlatform/review-middleware/pkg/client/review"
)

func (h *Handler) GetKyc(ctx context.Context) (*kycmwpb.Kyc, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := kycmwcli.GetKyc(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	rinfo, err := reviewmwcli.GetObjectReview(
		ctx,
		info.AppID,
		servicename.ServiceName,
		info.ID,
		reviewmgrpb.ReviewObjectType_ObjectKyc,
	)
	if err != nil {
		return nil, err
	}
	if rinfo == nil {
		return nil, fmt.Errorf("invalid review")
	}

	if rinfo.State == reviewmgrpb.ReviewState_Rejected {
		info.ReviewMessage = rinfo.GetMessage()
	}

	return info, nil
}

func (h *Handler) GetKycs(ctx context.Context) ([]*kycmwpb.Kyc, uint32, error) {
	conds := &kycmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}

	infos, total, err := kycmwcli.GetKycs(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	ids := []string{}
	for _, info := range infos {
		ids = append(ids, info.ID)
	}

	rinfos, err := reviewmwcli.GetObjectReviews(
		ctx,
		infos[0].AppID,
		servicename.ServiceName,
		ids,
		reviewmgrpb.ReviewObjectType_ObjectKyc,
	)
	if err != nil {
		return nil, 0, err
	}

	for _, rinfo := range rinfos {
		for _, info := range infos {
			if info.ID == rinfo.ObjectID && rinfo.State == reviewmgrpb.ReviewState_Rejected {
				info.ReviewMessage = rinfo.Message
				break
			}
		}
	}
	return infos, total, nil
}
