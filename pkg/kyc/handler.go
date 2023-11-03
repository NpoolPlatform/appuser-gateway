package kyc

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID                    *uint32
	EntID                 *string
	AppID                 *string
	UserID                *string
	DocumentType          *basetypes.KycDocumentType
	IDNumber              *string
	FrontImg              *string
	BackImg               *string
	SelfieImg             *string
	EntityType            *basetypes.KycEntityType
	ReviewID              *string
	State                 *basetypes.KycState
	ImageType             *basetypes.KycImageType
	RequestTimeoutSeconds *int64
	Offset                int32
	Limit                 int32
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	requestTimeoutSeconds := int64(10)
	handler := &Handler{
		RequestTimeoutSeconds: &requestTimeoutSeconds,
	}
	for _, opt := range options {
		if err := opt(ctx, handler); err != nil {
			return nil, err
		}
	}
	return handler, nil
}

func WithID(id *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid id")
			}
			return nil
		}
		h.ID = id
		return nil
	}
}

func WithEntID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid entid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.EntID = id
		return nil
	}
}

func WithAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid appid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		exist, err := appmwcli.ExistApp(ctx, *id)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("app not exist")
		}
		h.AppID = id
		return nil
	}
}

func WithUserID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid userid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		// Here shoud check app/user exist at low level
		h.UserID = id
		return nil
	}
}

func WithDocumentType(docType *basetypes.KycDocumentType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if docType == nil {
			if must {
				return fmt.Errorf("invalid documenttype")
			}
			return nil
		}
		switch *docType {
		case basetypes.KycDocumentType_IDCard:
		case basetypes.KycDocumentType_DriverLicense:
		case basetypes.KycDocumentType_Passport:
		default:
			return fmt.Errorf("invalid documenttype")
		}
		h.DocumentType = docType
		return nil
	}
}

func WithIDNumber(idNumber *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if idNumber == nil {
			if must {
				return fmt.Errorf("invalid idnumber")
			}
			return nil
		}
		const leastIDNumberLen = 8
		if len(*idNumber) < leastIDNumberLen {
			return fmt.Errorf("invalid id number")
		}
		h.IDNumber = idNumber
		return nil
	}
}

func WithImage(imgType *basetypes.KycImageType, img *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if imgType == nil {
			if must {
				return fmt.Errorf("invalid imagetype")
			}
			return nil
		}
		if img != nil && *img == "" {
			if must {
				return fmt.Errorf("invalid image")
			}
			return nil
		}
		switch *imgType {
		case basetypes.KycImageType_FrontImg:
			h.FrontImg = img
		case basetypes.KycImageType_BackImg:
			h.BackImg = img
		case basetypes.KycImageType_SelfieImg:
			h.SelfieImg = img
		default:
			return fmt.Errorf("invalid imagetype")
		}
		h.ImageType = imgType
		return nil
	}
}

func WithBackImg(img *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.BackImg = img
		return nil
	}
}

func WithSelfieImg(img *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.SelfieImg = img
		return nil
	}
}

func WithEntityType(entType *basetypes.KycEntityType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if entType == nil {
			if must {
				return fmt.Errorf("invalid entitytype")
			}
			return nil
		}
		switch *entType {
		case basetypes.KycEntityType_Individual:
		case basetypes.KycEntityType_Organization:
		default:
			return fmt.Errorf("invalid entitytype")
		}
		h.EntityType = entType
		return nil
	}
}

func WithReviewID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid reviewid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ReviewID = id
		return nil
	}
}

func WithState(state *basetypes.KycState, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if state == nil {
			if must {
				return fmt.Errorf("invalid state")
			}
			return nil
		}
		switch *state {
		case basetypes.KycState_Approved:
		case basetypes.KycState_Reviewing:
		case basetypes.KycState_Rejected:
		default:
			return fmt.Errorf("invalid state")
		}
		h.State = state
		return nil
	}
}

func WithRequestTimeoutSeconds(seconds *int64, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.RequestTimeoutSeconds = seconds
		return nil
	}
}

func WithOffset(offset int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Offset = offset
		return nil
	}
}

func WithLimit(limit int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if limit == 0 {
			limit = constant.DefaultRowLimit
		}
		h.Limit = limit
		return nil
	}
}
