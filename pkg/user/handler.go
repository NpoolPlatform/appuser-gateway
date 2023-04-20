package user

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	AppID                 string
	App                   *appmwpb.App
	UserID                string
	User                  *usermwpb.User
	Account               *string
	NewAccount            *string
	PasswordHash          *string
	OldPasswordHash       *string
	AccountType           *basetypes.SignMethod
	NewAccountType        *basetypes.SignMethod
	VerificationCode      string
	NewVerificationCode   string
	InvitationCode        *string
	EmailAddress          *string
	PhoneNO               *string
	RequestTimeoutSeconds int64
	ManMachineSpec        string
	Metadata              *Metadata
	Token                 *string
	Username              *string
	AddressFields         []string
	Gender                *string
	PostalCode            *string
	Age                   *uint32
	Birthday              *uint32
	Avatar                *string
	Organization          *string
	FirstName             *string
	LastName              *string
	IDNumber              *string
	SigninVerifyType      *basetypes.SignMethod
	KolConfirmed          *bool
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

func WithAppID(id string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		app, err := appmwcli.GetApp(ctx, id)
		if err != nil {
			return err
		}
		if app == nil {
			return fmt.Errorf("invalid app")
		}
		if _, err := uuid.Parse(id); err != nil {
			return err
		}
		h.AppID = id
		h.App = app
		return nil
	}
}

func WithPasswordHash(pwdHash string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if pwdHash == "" {
			return fmt.Errorf("invalid password")
		}
		h.PasswordHash = &pwdHash
		return nil
	}
}

func validateEmailAddress(emailAddress string) error {
	if _, err := mail.ParseAddress(emailAddress); err != nil {
		return err
	}
	return nil
}

func validatePhoneNO(phoneNO string) error {
	re := regexp.MustCompile(
		`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[` +
			`\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?)` +
			`{0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)` +
			`[\-\.\ \\\/]?(\d+))?$`,
	)
	if !re.MatchString(phoneNO) {
		return fmt.Errorf("invalid phone no")
	}

	return nil
}

func WithAccount(account string, accountType basetypes.SignMethod) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if account == "" {
			return fmt.Errorf("invalid account")
		}

		var err error

		switch accountType {
		case basetypes.SignMethod_Mobile:
			h.PhoneNO = &account
			err = validatePhoneNO(account)
		case basetypes.SignMethod_Email:
			h.EmailAddress = &account
			err = validateEmailAddress(account)
		default:
			return fmt.Errorf("invalid account type")
		}

		if err != nil {
			return err
		}

		h.AccountType = &accountType
		h.Account = &account
		return nil
	}
}

func WithVerificationCode(code string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == "" {
			return fmt.Errorf("invalid verification code")
		}
		h.VerificationCode = code
		return nil
	}
}

func WithInvitationCode(code *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == nil {
			return nil
		}
		if *code == "" {
			return fmt.Errorf("invalid invitation code")
		}
		h.InvitationCode = code
		return nil
	}
}

func WithRequestTimeoutSeconds(seconds int64) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.RequestTimeoutSeconds = seconds
		return nil
	}
}

func WithManMachineSpec(manMachineSpec string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ManMachineSpec = manMachineSpec
		return nil
	}
}
