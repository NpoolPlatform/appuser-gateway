package user

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	AppID                 string
	App                   *appmwpb.App
	UserID                *string
	User                  *usermwpb.User
	TargetUserID          *string
	TargetUser            *usermwpb.User
	CheckInvitation       *bool
	Account               *string
	NewAccount            *string
	PasswordHash          *string
	OldPasswordHash       *string
	AccountType           *basetypes.SignMethod
	NewAccountType        *basetypes.SignMethod
	VerificationCode      *string
	NewVerificationCode   *string
	InvitationCode        *string
	EmailAddress          *string
	PhoneNO               *string
	RequestTimeoutSeconds int64
	ManMachineSpec        *string
	EnvironmentSpec       *string
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
	SelectedLangID        *string
	Kol                   *bool
	GoogleSecret          *string
	GoogleAuthVerified    *bool
	Banned                *bool
	BanMessage            *string
	RecoveryCode          *string
	ShouldUpdateCache     bool
	Offset                int32
	Limit                 int32
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	handler := &Handler{
		ShouldUpdateCache: true,
	}

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
			return fmt.Errorf("get app error: %v", err)
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

func WithUserID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.UserID = id
		return nil
	}
}

func WithPasswordHash(pwdHash *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if pwdHash == nil {
			return nil
		}
		if *pwdHash == "" {
			return fmt.Errorf("invalid password")
		}
		h.PasswordHash = pwdHash
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

//nolint:gocyclo
func WithAccount(account *string, accountType *basetypes.SignMethod) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if accountType == nil {
			return nil
		}

		switch *accountType {
		case basetypes.SignMethod_Mobile:
			fallthrough //nolint
		case basetypes.SignMethod_Email:
			fallthrough //nolint
		case basetypes.SignMethod_Github:
			fallthrough //nolint
		case basetypes.SignMethod_Google:
			fallthrough //nolint
		case basetypes.SignMethod_Facebook:
			fallthrough //nolint
		case basetypes.SignMethod_Twitter:
			fallthrough //nolint
		case basetypes.SignMethod_Linkedin:
			fallthrough //nolint
		case basetypes.SignMethod_Wechat:
			if account == nil {
				return fmt.Errorf("invalid account")
			}
		}

		var err error

		switch *accountType {
		case basetypes.SignMethod_Mobile:
			h.PhoneNO = account
			err = validatePhoneNO(*account)
		case basetypes.SignMethod_Email:
			h.EmailAddress = account
			err = validateEmailAddress(*account)
		case basetypes.SignMethod_Github:
		case basetypes.SignMethod_Google:
		case basetypes.SignMethod_Facebook:
		case basetypes.SignMethod_Twitter:
		case basetypes.SignMethod_Linkedin:
		case basetypes.SignMethod_Wechat:
		default:
			return fmt.Errorf("invalid account type")
		}

		if err != nil {
			return err
		}

		h.AccountType = accountType
		h.Account = account
		return nil
	}
}

func WithEmailAddress(emailAddress *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if emailAddress == nil {
			return nil
		}
		if err := validateEmailAddress(*emailAddress); err != nil {
			return err
		}
		h.EmailAddress = emailAddress
		return nil
	}
}

func WithPhoneNO(phoneNO *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if phoneNO == nil {
			return nil
		}
		if err := validatePhoneNO(*phoneNO); err != nil {
			return err
		}
		h.PhoneNO = phoneNO
		return nil
	}
}

func WithVerificationCode(code *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == nil {
			return nil
		}
		if *code == "" {
			return fmt.Errorf("invalid verification code")
		}
		h.VerificationCode = code
		return nil
	}
}

func WithNewVerificationCode(code *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == nil {
			return nil
		}
		if *code == "" {
			return fmt.Errorf("invalid new verification code")
		}
		h.NewVerificationCode = code
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
		h.ManMachineSpec = &manMachineSpec
		return nil
	}
}

func WithEnvironmentSpec(envSpec string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.EnvironmentSpec = &envSpec
		return nil
	}
}

func WithToken(token string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Token = &token
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

func WithNewAccount(account *string, accountType *basetypes.SignMethod) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if accountType == nil {
			return nil
		}

		switch *accountType {
		case basetypes.SignMethod_Mobile:
			fallthrough //nolint
		case basetypes.SignMethod_Email:
			if account == nil {
				return fmt.Errorf("invalid account")
			}
		}

		var err error

		switch *accountType {
		case basetypes.SignMethod_Mobile:
			err = validatePhoneNO(*account)
		case basetypes.SignMethod_Email:
			err = validateEmailAddress(*account)
		case basetypes.SignMethod_Google:
		default:
			return fmt.Errorf("invalid account type")
		}

		if err != nil {
			return err
		}

		h.NewAccountType = accountType
		h.NewAccount = account
		return nil
	}
}

func WithOldPasswordHash(pwdHash *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if pwdHash == nil {
			return nil
		}
		if *pwdHash == "" {
			return fmt.Errorf("invalid old password")
		}
		h.OldPasswordHash = pwdHash
		return nil
	}
}

func WithRecoveryCode(code *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == nil {
			return nil
		}
		if *code == "" {
			return fmt.Errorf("invalid code")
		}
		h.RecoveryCode = code
		return nil
	}
}

func WithTargetUserID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.TargetUserID = id
		return nil
	}
}

func WithCheckInvitation(check bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.CheckInvitation = &check
		return nil
	}
}

func WithKol(kol *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Kol = kol
		return nil
	}
}

func WithUsername(username *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if username == nil {
			return nil
		}
		re := regexp.MustCompile("^[a-zA-Z0-9\u3040-\u31ff][[a-zA-Z0-9_\\-\\.\u3040-\u31ff]{3,32}$") //nolint
		if !re.MatchString(*username) {
			return fmt.Errorf("invalid username")
		}
		h.Username = username
		return nil
	}
}

func WithAddressFields(addressFields []string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.AddressFields = addressFields
		return nil
	}
}

func WithGender(gender *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if gender == nil {
			return nil
		}
		if *gender == "" {
			return fmt.Errorf("invalid gender")
		}
		h.Gender = gender
		return nil
	}
}

func WithPostalCode(postalCode *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if postalCode == nil {
			return nil
		}
		if *postalCode == "" {
			return fmt.Errorf("invalid postalCode")
		}
		h.PostalCode = postalCode
		return nil
	}
}

func WithAge(age *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if age == nil {
			return nil
		}
		h.Age = age
		return nil
	}
}

func WithBirthday(birthday *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if birthday == nil {
			return nil
		}
		h.Birthday = birthday
		return nil
	}
}

func WithAvatar(avatar *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if avatar == nil {
			return nil
		}
		if *avatar == "" {
			return fmt.Errorf("invalid avatar")
		}
		h.Avatar = avatar
		return nil
	}
}

func WithOrganization(organization *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if organization == nil {
			return nil
		}
		if *organization == "" {
			return fmt.Errorf("invalid organization")
		}
		h.Organization = organization
		return nil
	}
}

func WithFirstName(firstName *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if firstName == nil {
			return nil
		}
		if *firstName == "" {
			return fmt.Errorf("invalid firstname")
		}
		h.FirstName = firstName
		return nil
	}
}

func WithLastName(lastName *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if lastName == nil {
			return nil
		}
		if *lastName == "" {
			return fmt.Errorf("invalid lastname")
		}
		h.LastName = lastName
		return nil
	}
}

func WithIDNumber(idNumber *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if idNumber == nil {
			return nil
		}
		if *idNumber == "" {
			return fmt.Errorf("invalid idnumber")
		}
		h.IDNumber = idNumber
		return nil
	}
}

func WithSigninVerifyType(verifyType *basetypes.SignMethod) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if verifyType == nil {
			return nil
		}
		switch *verifyType {
		case basetypes.SignMethod_Email:
		case basetypes.SignMethod_Mobile:
		case basetypes.SignMethod_Google:
		default:
			return fmt.Errorf("invalid sign verify type")
		}
		h.SigninVerifyType = verifyType
		return nil
	}
}

func WithKolConfirmed(confirmed *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.KolConfirmed = confirmed
		return nil
	}
}

func WithSelectedLangID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.SelectedLangID = id
		return nil
	}
}

func WithBanned(banned *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Banned = banned
		return nil
	}
}

func WithBanMessage(message *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.BanMessage = message
		return nil
	}
}

func WithShouldUpdateCache(update bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ShouldUpdateCache = update
		return nil
	}
}
