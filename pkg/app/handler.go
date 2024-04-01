package app

import (
	"context"
	"fmt"
	"time"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	appusertypes "github.com/NpoolPlatform/message/npool/basetypes/appuser/v1"
	inspiretypes "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID                       *uint32
	EntID                    *string
	NewEntID                 *string
	EntIDs                   []string
	CreatedBy                *string
	Name                     *string
	Logo                     *string
	Description              *string
	SignupMethods            []basetypes.SignMethod
	ExtSigninMethods         []basetypes.SignMethod
	RecaptchaMethod          *basetypes.RecaptchaMethod
	KycEnable                *bool
	SigninVerifyEnable       *bool
	InvitationCodeMust       *bool
	CreateInvitationCodeWhen *basetypes.CreateInvitationCodeWhen
	MaxTypedCouponsPerOrder  *uint32
	Maintaining              *bool
	CouponWithdrawEnable     *bool
	CommitButtonTargets      []string
	Banned                   *bool
	BanMessage               *string
	ResetUserMethod          *appusertypes.ResetUserMethod
	SettleMode               *inspiretypes.SettleMode
	SettleAmountType         *inspiretypes.SettleAmountType
	SettleInterval           *inspiretypes.SettleInterval
	CommissionType           *inspiretypes.CommissionType
	SettleBenefit            *bool
	StartAt                  *uint32
	RequestTimeoutSeconds    *int64
	Offset                   int32
	Limit                    int32
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	requestTimeoutSeconds := int64(10) //nolint
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

func WithNewEntID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid newentid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.NewEntID = id
		return nil
	}
}

func WithCreatedBy(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid createdby")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.CreatedBy = id
		return nil
	}
}

func WithName(name *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if name == nil {
			if must {
				return fmt.Errorf("invalid name")
			}
			return nil
		}
		const leastNameLen = 3
		if len(*name) < leastNameLen {
			return fmt.Errorf("invalid name")
		}
		h.Name = name
		return nil
	}
}

func WithLogo(logo *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if logo == nil {
			if must {
				return fmt.Errorf("invalid logo")
			}
			return nil
		}
		const leastLogoLen = 5
		if len(*logo) < leastLogoLen {
			return fmt.Errorf("invalid logo")
		}
		h.Logo = logo
		return nil
	}
}

func WithDescription(description *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Description = description
		return nil
	}
}

func WithSignupMethods(methods []basetypes.SignMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		for _, method := range methods {
			switch method {
			case basetypes.SignMethod_Username:
				return fmt.Errorf("username signup not implemented")
			case basetypes.SignMethod_Mobile:
			case basetypes.SignMethod_Email:
			default:
				return fmt.Errorf("signup method %v invalid", method)
			}
		}
		h.SignupMethods = methods
		return nil
	}
}

func WithExtSigninMethods(methods []basetypes.SignMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		for _, method := range methods {
			switch method {
			case basetypes.SignMethod_Twitter:
				fallthrough //nolint
			case basetypes.SignMethod_Github:
				fallthrough //nolint
			case basetypes.SignMethod_Facebook:
				fallthrough //nolint
			case basetypes.SignMethod_Linkedin:
				fallthrough //nolint
			case basetypes.SignMethod_Wechat:
				fallthrough //nolint
			case basetypes.SignMethod_Google:
				return fmt.Errorf("%v signin not implemented", method)
			default:
				return fmt.Errorf("ext signin method %v invalid", method)
			}
		}
		h.ExtSigninMethods = methods
		return nil
	}
}

func WithRecaptchaMethod(method *basetypes.RecaptchaMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if method == nil {
			if must {
				return fmt.Errorf("invalid recaptchamethod")
			}
			return nil
		}
		switch *method {
		case basetypes.RecaptchaMethod_GoogleRecaptchaV3:
		case basetypes.RecaptchaMethod_NoRecaptcha:
		default:
			return fmt.Errorf("invalid recaptchamethod")
		}
		h.RecaptchaMethod = method
		return nil
	}
}

func WithKycEnable(enable *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.KycEnable = enable
		return nil
	}
}

func WithSigninVerifyEnable(enable *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.SigninVerifyEnable = enable
		return nil
	}
}

func WithInvitationCodeMust(must *bool, _must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.InvitationCodeMust = must
		return nil
	}
}

func WithCreateInvitationCodeWhen(when *basetypes.CreateInvitationCodeWhen, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if when == nil {
			if must {
				return fmt.Errorf("invalid createinvitationcodewhen")
			}
			return nil
		}
		switch *when {
		case basetypes.CreateInvitationCodeWhen_Registration:
		case basetypes.CreateInvitationCodeWhen_SetToKol:
		case basetypes.CreateInvitationCodeWhen_HasPaidOrder:
		default:
			return fmt.Errorf("invalid createinvitationcodewhen")
		}
		h.CreateInvitationCodeWhen = when
		return nil
	}
}

func WithMaxTypedCouponsPerOrder(max *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.MaxTypedCouponsPerOrder = max
		return nil
	}
}

func WithMaintaining(maintaining *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Maintaining = maintaining
		return nil
	}
}

func WithCouponWithdrawEnable(enable *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.CouponWithdrawEnable = enable
		return nil
	}
}

func WithCommitButtonTargets(targets []string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.CommitButtonTargets = targets
		return nil
	}
}

func WithResetUserMethod(method *appusertypes.ResetUserMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if method == nil {
			if must {
				return fmt.Errorf("invalid reset user method")
			}
			return nil
		}
		switch *method {
		case appusertypes.ResetUserMethod_Normal:
		case appusertypes.ResetUserMethod_Link:
		default:
			return fmt.Errorf("invalid reset method %v", method)
		}
		h.ResetUserMethod = method
		return nil
	}
}

func WithCommissionType(commissionType *inspiretypes.CommissionType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if commissionType == nil {
			if must {
				return fmt.Errorf("invalid commissiontype")
			}
			return nil
		}
		switch *commissionType {
		case inspiretypes.CommissionType_LegacyCommission:
		case inspiretypes.CommissionType_LayeredCommission:
		case inspiretypes.CommissionType_DirectCommission:
		default:
			return fmt.Errorf("invalid commissiontype")
		}
		h.CommissionType = commissionType
		return nil
	}
}

func WithSettleAmountType(settleAmountType *inspiretypes.SettleAmountType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if settleAmountType == nil {
			if must {
				return fmt.Errorf("invalid settleamounttype")
			}
			return nil
		}
		switch *settleAmountType {
		case inspiretypes.SettleAmountType_SettleByPercent:
		case inspiretypes.SettleAmountType_SettleByAmount:
		default:
			return fmt.Errorf("invalid settleamounttype")
		}
		h.SettleAmountType = settleAmountType
		return nil
	}
}

func WithSettleMode(settleMode *inspiretypes.SettleMode, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if settleMode == nil {
			if must {
				return fmt.Errorf("invalid settlemode")
			}
			return nil
		}
		switch *settleMode {
		case inspiretypes.SettleMode_SettleWithPaymentAmount:
		case inspiretypes.SettleMode_SettleWithGoodValue:
		default:
			return fmt.Errorf("invalid settlemode")
		}
		h.SettleMode = settleMode
		return nil
	}
}

func WithSettleInterval(settleInterval *inspiretypes.SettleInterval, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if settleInterval == nil {
			if must {
				return fmt.Errorf("invalid settleinterval")
			}
			return nil
		}
		switch *settleInterval {
		case inspiretypes.SettleInterval_SettleAggregate:
		case inspiretypes.SettleInterval_SettleYearly:
		case inspiretypes.SettleInterval_SettleMonthly:
		case inspiretypes.SettleInterval_SettleEveryOrder:
		default:
			return fmt.Errorf("invalid settleinterval")
		}
		h.SettleInterval = settleInterval
		return nil
	}
}

func WithStartAt(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid startat")
			}
			return nil
		}
		if *value == 0 {
			*value = uint32(time.Now().Unix())
		}
		h.StartAt = value
		return nil
	}
}

func WithSettleBenefit(value *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid settlebenefit")
			}
			return nil
		}
		h.SettleBenefit = value
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
		if limit <= 0 {
			limit = constant.DefaultRowLimit
		}
		h.Limit = limit
		return nil
	}
}

func WithEntIDs(ids []string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.EntIDs = ids
		return nil
	}
}

func WithBanned(banned *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Banned = banned
		return nil
	}
}

func WithBanMessage(message *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.BanMessage = message
		return nil
	}
}
