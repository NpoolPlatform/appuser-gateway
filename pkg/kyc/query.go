package kyc

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	mwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	kycmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"

	reviewmgrpb "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	reviewmwcli "github.com/NpoolPlatform/review-middleware/pkg/client/review"
)

func (h *Handler) GetKyc(ctx context.Context) (*kycmwpb.Kyc, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := mwcli.GetKyc(ctx, *h.ID)
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

func GetKyc(ctx context.Context, id string) (*kycmwpb.Kyc, error) {
	info, err := mwcli.GetKyc(ctx, id)
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

func GetKycs(ctx context.Context, conds *kycmwpb.Conds, offset, limit int32) ([]*kycmwpb.Kyc, uint32, error) {
	infos, total, err := mwcli.GetKycs(ctx, conds, offset, limit)
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
