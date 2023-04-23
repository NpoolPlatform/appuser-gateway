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
	ID           *string
	AppID        string
	UserID       string
	DocumentType *basetypes.KycDocumentType
	IDNumber     *string
	FrontImg     *string
	BackImg      *string
	SelfieImg    *string
	EntityType   *basetypes.KycEntityType
	ReviewID     *string
	State        *basetypes.KycState
	Offset       int32
	Limit        int32
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	handler := &Handler{}
	for _, opt := range options {
		if err := opt(ctx, handler); err != nil {
			return nil, err
		}
	}
	return handler, nil
}

func WithID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ID = id
		return nil
	}
}

func WithAppID(id string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _, err := uuid.Parse(id); err != nil {
			return err
		}
		exist, err := appmwcli.ExistApp(ctx, id)
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

func WithUserID(id string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _, err := uuid.Parse(id); err != nil {
			return err
		}
		// Here shoud check app/user exist at low level
		h.UserID = id
		return nil
	}
}

func WithDocumentType(docType *basetypes.KycDocumentType) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if docType == nil {
			return nil
		}
		switch *docType {
		case basetypes.KycDocumentType_IDCard:
		case basetypes.KycDocumentType_DriverLicense:
		case basetypes.KycDocumentType_Passport:
		default:
			return fmt.Errorf("invalid document type")
		}
		h.DocumentType = docType
		return nil
	}
}

func WithIDNumber(idNumber *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if idNumber == nil {
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

func WithFrontImg(img *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.FrontImg = img
		return nil
	}
}

func WithBackImg(img *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.BackImg = img
		return nil
	}
}

func WithSelfieImg(img *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.SelfieImg = img
		return nil
	}
}

func WithEntityType(entType *basetypes.KycEntityType) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if entType == nil {
			return nil
		}
		switch *entType {
		case basetypes.KycEntityType_Individual:
		case basetypes.KycEntityType_Organization:
		default:
			return fmt.Errorf("invalid entity type")
		}
		h.EntityType = entType
		return nil
	}
}

func WithReviewID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ReviewID = id
		return nil
	}
}

func WithState(state *basetypes.KycState) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if state == nil {
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
