package kyc

import (
	"context"

	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
)

func UploadKycImage(ctx context.Context, appID, userID string, imgType kycmgrpb.KycImageType, imgBase64 string) error {
	return nil
}

func GetKycImage(ctx context.Context, appID, userID string, imgType kycmgrpb.KycImageType) (string, error) {
	return "", nil
}
