package kyc

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"

	mwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"

	reviewmgrpb "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	reviewmwcli "github.com/NpoolPlatform/review-middleware/pkg/client/review"
)

func GetKyc(ctx context.Context, id string) (*mwpb.Kyc, error) {
	info, err := mwcli.GetKyc(ctx, id)
	if err != nil {
		return nil, err
	}

	rinfo, err := reviewmwcli.GetObjectReview(
		ctx,
		info.AppID,
		constant.ServiceName,
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

func GetKycs(ctx context.Context, conds *mwpb.Conds, offset, limit int32) ([]*mwpb.Kyc, uint32, error) {
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
		constant.ServiceName,
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
