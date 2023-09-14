package kyc

import (
	"context"
	"fmt"
	"sort"

	"github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	kycmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	kycmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	reviewtypes "github.com/NpoolPlatform/message/npool/basetypes/review/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
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
		servicename.ServiceDomain,
		info.ID,
		reviewtypes.ReviewObjectType_ObjectKyc,
	)
	if err != nil {
		return nil, err
	}
	if rinfo == nil {
		return nil, fmt.Errorf(
			"invalid review: app_id=%v, domain=%v, id=%v",
			info.AppID,
			servicename.ServiceDomain,
			info.ID,
		)
	}

	if rinfo.State == reviewtypes.ReviewState_Rejected {
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
		servicename.ServiceDomain,
		ids,
		reviewtypes.ReviewObjectType_ObjectKyc,
	)
	if err != nil {
		return nil, 0, err
	}

	// TODO: here we should only get the last one reviews of different state
	sort.Slice(rinfos, func(i, j int) bool {
		return rinfos[i].CreatedAt >= rinfos[j].CreatedAt
	})

	for _, info := range infos {
		for _, rinfo := range rinfos {
			if info.ID == rinfo.ObjectID {
				info.ReviewMessage = rinfo.Message
				break
			}
		}
	}
	return infos, total, nil
}
