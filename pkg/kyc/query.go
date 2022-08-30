package kyc

import (
	"context"
	"fmt"

	mwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	reviewmgrcli "github.com/NpoolPlatform/review-service/pkg/client"
)

func GetKyc(ctx context.Context, id string) (*mwpb.Kyc, error) {
	info, err := mwcli.GetKyc(ctx, id)
	if err != nil {
		return nil, err
	}

	rinfo, err := reviewmgrcli.GetReview(ctx, info.GetReviewID())
	if err != nil {
		return nil, err
	}
	if rinfo == nil {
		return nil, fmt.Errorf("invalid review")
	}

	if rinfo.State == "rejected" {
		info.ReviewMessage = rinfo.GetMessage()
	}

	return info, nil
}

func GetKycs(ctx context.Context, conds *mwpb.Conds, offset, limit int32) ([]*mwpb.Kyc, uint32, error) {
	infos, total, err := mwcli.GetKycs(ctx, conds, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	for key, val := range infos {
		info, err := reviewmgrcli.GetReview(ctx, val.GetReviewID())
		if err != nil {
			return nil, 0, err
		}
		if info == nil {
			continue
		}

		if info.State == "rejected" {
			infos[key].ReviewMessage = info.GetMessage()
		}
	}
	return infos, total, nil
}
