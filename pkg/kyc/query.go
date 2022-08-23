package kyc

import (
	"context"
	"strings"

	mwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
	reviewmgrpb "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	reviewmgrcli "github.com/NpoolPlatform/review-service/pkg/client"
)

func GetKyc(ctx context.Context, id string) (*mwpb.Kyc, error) {
	info, err := mwcli.GetKyc(ctx, id)
	if err != nil {
		return nil, err
	}

	reviewInfo, err := reviewmgrcli.GetReview(ctx, info.GetReviewID())
	if err != nil {
		return nil, err
	}
	// TODO: review state wait for review migration remove ToUpper
	info.ReviewState = reviewmgrpb.ReviewState(reviewmgrpb.ReviewState_value[strings.ToUpper(reviewInfo.State[:1])+reviewInfo.State[1:]])
	info.ReviewMessage = reviewInfo.GetMessage()
	return info, nil
}

func GetKycs(ctx context.Context, conds *mwpb.Conds, offset, limit int32) ([]*mwpb.Kyc, uint32, error) {
	infos, total, err := mwcli.GetKycs(ctx, conds, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	for key, val := range infos {
		reviewInfo, err := reviewmgrcli.GetReview(ctx, val.GetReviewID())
		if err != nil {
			return nil, 0, err
		}
		// TODO: review state wait for review migration remove ToUpper
		infos[key].ReviewState = reviewmgrpb.ReviewState(
			reviewmgrpb.ReviewState_value[strings.ToUpper(reviewInfo.State[:1])+reviewInfo.State[1:]],
		)
		infos[key].ReviewMessage = reviewInfo.GetMessage()
	}
	return infos, total, nil
}
