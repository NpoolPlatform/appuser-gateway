package kyc

import (
	"context"
	"fmt"

	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"

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

	switch rinfo.State {
	case "wait":
		info.State = kycmgrpb.KycState_Reviewing
	case "rejected":
		info.State = kycmgrpb.KycState_Rejected
		info.ReviewMessage = rinfo.GetMessage()
	case "approved":
		info.State = kycmgrpb.KycState_Approved
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
		switch info.State {
		case "wait":
			infos[key].State = kycmgrpb.KycState_Reviewing
		case "rejected":
			infos[key].State = kycmgrpb.KycState_Rejected
			infos[key].ReviewMessage = info.GetMessage()
		case "approved":
			infos[key].State = kycmgrpb.KycState_Approved
		}
	}
	return infos, total, nil
}
