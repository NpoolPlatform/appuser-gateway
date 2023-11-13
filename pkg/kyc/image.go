package kyc

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/go-service-framework/pkg/oss"
)

func (h *Handler) UploadKycImage(ctx context.Context) (string, error) {
	if h.UserID == nil {
		return "", fmt.Errorf("invalid userid")
	}
	existUser, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return "", err
	}
	if !existUser {
		return "", fmt.Errorf("invalid user")
	}

	if h.ImageType == nil {
		return "", fmt.Errorf("invalid image type")
	}
	var image *string
	switch *h.ImageType {
	case basetypes.KycImageType_FrontImg:
		image = h.FrontImg
	case basetypes.KycImageType_BackImg:
		image = h.BackImg
	case basetypes.KycImageType_SelfieImg:
		image = h.SelfieImg
	default:
		return "", fmt.Errorf("invalid image type")
	}
	if image == nil || *image == "" {
		return "", fmt.Errorf("invalid image")
	}
	key := fmt.Sprintf("kyc/%v/%v/%v", h.AppID, *h.UserID, *h.ImageType)
	return key, oss.PutObject(ctx, key, []byte(*image), true)
}

//nolint:gocyclo
func (h *Handler) GetKycImage(ctx context.Context) (string, error) {
	if h.UserID == nil {
		return "", fmt.Errorf("invalid userid")
	}
	if h.ImageType == nil {
		return "", fmt.Errorf("invalid image type")
	}
	switch *h.ImageType {
	case basetypes.KycImageType_FrontImg:
	case basetypes.KycImageType_BackImg:
	case basetypes.KycImageType_SelfieImg:
	default:
		return "", fmt.Errorf("invalid image type")
	}

	key := fmt.Sprintf("kyc/%v/%v/%v", h.AppID, *h.UserID, *h.ImageType)
	imgBase64, err := oss.GetObject(ctx, key, true)
	if err == nil && imgBase64 != nil {
		return string(imgBase64), nil
	}

	switch *h.ImageType {
	case basetypes.KycImageType_FrontImg:
		key = fmt.Sprintf("kyc/%v/%v/front", h.AppID, *h.UserID)
	case basetypes.KycImageType_BackImg:
		key = fmt.Sprintf("kyc/%v/%v/back", h.AppID, *h.UserID)
	case basetypes.KycImageType_SelfieImg:
		key = fmt.Sprintf("kyc/%v/%v/handing", h.AppID, *h.UserID)
	default:
		return "", fmt.Errorf("invalid image type")
	}
	imgBase64, err = oss.GetObject(ctx, key, true)
	if err == nil && imgBase64 != nil {
		return string(imgBase64), nil
	}

	switch *h.ImageType {
	case basetypes.KycImageType_FrontImg:
		key = fmt.Sprintf("kyc/%v/%v/Front", h.AppID, *h.UserID)
	case basetypes.KycImageType_BackImg:
		key = fmt.Sprintf("kyc/%v/%v/Back", h.AppID, *h.UserID)
	case basetypes.KycImageType_SelfieImg:
		key = fmt.Sprintf("kyc/%v/%v/Handing", h.AppID, *h.UserID)
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
