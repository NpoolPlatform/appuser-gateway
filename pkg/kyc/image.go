package kyc

import (
	"context"
	"fmt"

	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"

	"github.com/NpoolPlatform/go-service-framework/pkg/oss"
)

func UploadKycImage(ctx context.Context, appID, userID string, imgType kycmgrpb.KycImageType, imgBase64 string) error {
	key := fmt.Sprintf("kyc/%v/%v/%v", appID, userID, imgType)
	return oss.PutObject(ctx, key, []byte(imgBase64), true)
}

func GetKycImage(ctx context.Context, appID, userID string, imgType kycmgrpb.KycImageType) (string, error) {
	key := fmt.Sprintf("kyc/%v/%v/%v", appID, userID, imgType)
	imgBase64, err := oss.GetObject(ctx, key, true)
	if err != nil {
		return "", err
	}
	if imgBase64 != nil {
		return string(imgBase64), nil
	}

	switch imgType {
	case kycmgrpb.KycImageType_FrontImg:
		key = fmt.Sprintf("kyc/%v/%v/front", appID, userID)
	case kycmgrpb.KycImageType_BackImg:
		key = fmt.Sprintf("kyc/%v/%v/back", appID, userID)
	case kycmgrpb.KycImageType_SelfieImg:
		key = fmt.Sprintf("kyc/%v/%v/handing", appID, userID)
	default:
		return "", fmt.Errorf("invalid image type")
	}

	imgBase64, err = oss.GetObject(ctx, key, true)
	if err != nil {
		return "", err
	}
	if imgBase64 == nil {
		return "", fmt.Errorf("no image")
	}

	return string(imgBase64), nil
}