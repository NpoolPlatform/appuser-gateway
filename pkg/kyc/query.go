package kyc

import (
	"context"

	mwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/kyc"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"
)

func GetKyc(ctx context.Context, id string) (*mwpb.Kyc, error) {
	info, err := mwcli.GetKyc(ctx, id)
	if err != nil {
		return nil, err
	}
	// TODO: review state wait for review migration remove ToUpper
	return info, nil
}

func GetKycs(ctx context.Context, conds *mwpb.Conds, offset, limit int32) ([]*mwpb.Kyc, uint32, error) {
	infos, total, err := mwcli.GetKycs(ctx, conds, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return infos, total, nil
}
