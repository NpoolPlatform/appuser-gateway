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
	ID                    *uint32
	AppID                 *string
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
	ThirdPartyID          *string
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

func WithID(id *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid id")
			}
		}
		h.ID = id
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
		app, err := appmwcli.GetApp(ctx, *id)
		if err != nil {
			return err
		}
		if app == nil {
			return fmt.Errorf("invalid app")
		}
		h.AppID = id
		h.App = app
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
		h.UserID = id
		return nil
	}
}

func WithPasswordHash(pwdHash *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if pwdHash == nil {
			if must {
				return fmt.Errorf("invalid passwordhash")
			}
			return nil
		}
		if *pwdHash == "" {
			return fmt.Errorf("invalid passwordhash")
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
		return fmt.Errorf("invalid phoneno")
	}

	return nil
}

//nolint:gocyclo
func WithAccount(account *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if account == nil {
			if must {
				return fmt.Errorf("invalid account")
			}
		}
		if *account == "" {
			return fmt.Errorf("invalid account")
		}

		var accountType basetypes.SignMethod
		if err := validatePhoneNO(*account); err == nil {
			h.PhoneNO = account
			accountType = basetypes.SignMethod_Mobile
		} else if err := validateEmailAddress(*account); err == nil {
			accountType = basetypes.SignMethod_Email
			h.EmailAddress = account
		} else {
			return fmt.Errorf("invalid account")
		}

		if h.AccountType != nil && accountType != *h.AccountType {
			return fmt.Errorf("invalid accounttype")
		}

		h.AccountType = &accountType
		h.Account = account
		return nil
	}
}

func WithAccountType(accountType *basetypes.SignMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if accountType == nil {
			if must {
				return fmt.Errorf("invalid accounttype")
			}
			return nil
		}
		if h.AccountType != nil && *accountType != *h.AccountType {
			return fmt.Errorf("invalid accounttype")
		}
		switch *accountType {
		case basetypes.SignMethod_Mobile:
		case basetypes.SignMethod_Email:
		default:
			return fmt.Errorf("invalid accounttype")
		}
		h.AccountType = accountType
		return nil
	}
}

func WithEmailAddress(emailAddress *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if emailAddress == nil {
			if must {
				return fmt.Errorf("invalid emailaddress")
			}
			return nil
		}
		if err := validateEmailAddress(*emailAddress); err != nil {
			return err
		}
		h.EmailAddress = emailAddress
		return nil
	}
}

func WithPhoneNO(phoneNO *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if phoneNO == nil {
			if must {
				return fmt.Errorf("invalid phoneno")
			}
			return nil
		}
		if err := validatePhoneNO(*phoneNO); err != nil {
			return err
		}
		h.PhoneNO = phoneNO
		return nil
	}
}

func WithVerificationCode(code *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == nil {
			if must {
				return fmt.Errorf("invalid verificationcode")
			}
			return nil
		}
		if *code == "" {
			return fmt.Errorf("invalid verificationcode")
		}
		h.VerificationCode = code
		return nil
	}
}

func WithNewVerificationCode(code *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == nil {
			if must {
				return fmt.Errorf("invalid newverificationcode")
			}
			return nil
		}
		if *code == "" {
			return fmt.Errorf("invalid newverificationcode")
		}
		h.NewVerificationCode = code
		return nil
	}
}

func WithInvitationCode(code *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == nil {
			if must {
				return fmt.Errorf("invalid invitationcode")
			}
			return nil
		}
		if *code == "" {
			return fmt.Errorf("invalid invitationcode")
		}
		h.InvitationCode = code
		return nil
	}
}

func WithRequestTimeoutSeconds(seconds int64, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.RequestTimeoutSeconds = seconds
		return nil
	}
}

func WithManMachineSpec(manMachineSpec string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ManMachineSpec = &manMachineSpec
		return nil
	}
}

func WithEnvironmentSpec(envSpec string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.EnvironmentSpec = &envSpec
		return nil
	}
}

func WithToken(token string, must bool) func(context.Context, *Handler) error {
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

func WithNewAccount(account *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if account == nil {
			if must {
				return fmt.Errorf("invalid account")
			}
		}
		if *account == "" {
			return fmt.Errorf("invalid account")
		}

		var accountType basetypes.SignMethod
		if err := validatePhoneNO(*account); err == nil {
			accountType = basetypes.SignMethod_Mobile
		} else if err := validateEmailAddress(*account); err == nil {
			accountType = basetypes.SignMethod_Email
		} else {
			return fmt.Errorf("invalid account")
		}

		if h.NewAccountType != nil && accountType != *h.NewAccountType {
			return fmt.Errorf("invalid accounttype")
		}

		h.NewAccountType = &accountType
		h.NewAccount = account
		return nil
	}
}

func WithNewAccountType(accountType *basetypes.SignMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if accountType == nil {
			if must {
				return fmt.Errorf("invalid newaccounttype")
			}
			return nil
		}
		if h.NewAccountType != nil && *accountType != *h.NewAccountType {
			return fmt.Errorf("invalid newaccounttype")
		}
		switch *accountType {
		case basetypes.SignMethod_Mobile:
		case basetypes.SignMethod_Email:
		default:
			return fmt.Errorf("invalid newaccounttype")
		}
		h.NewAccountType = accountType
		return nil
	}
}

func WithOldPasswordHash(pwdHash *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if pwdHash == nil {
			if must {
				return fmt.Errorf("invalid oldpasswordhash")
			}
			return nil
		}
		if *pwdHash == "" {
			return fmt.Errorf("invalid oldpasswordhash")
		}
		h.OldPasswordHash = pwdHash
		return nil
	}
}

func WithRecoveryCode(code *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == nil {
			if must {
				return fmt.Errorf("invalid recoverycode")
			}
			return nil
		}
		if *code == "" {
			return fmt.Errorf("invalid recoverycode")
		}
		h.RecoveryCode = code
		return nil
	}
}

func WithTargetUserID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid targetuserid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.TargetUserID = id
		return nil
	}
}

func WithCheckInvitation(check bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.CheckInvitation = &check
		return nil
	}
}

func WithKol(kol *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Kol = kol
		return nil
	}
}

func WithUsername(username *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if username == nil {
			if must {
				return fmt.Errorf("invalid username")
			}
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

func WithAddressFields(addressFields []string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.AddressFields = addressFields
		return nil
	}
}

func WithGender(gender *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if gender == nil {
			if must {
				return fmt.Errorf("invalid gender")
			}
			return nil
		}
		if *gender == "" {
			return fmt.Errorf("invalid gender")
		}
		h.Gender = gender
		return nil
	}
}

func WithPostalCode(postalCode *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if postalCode == nil {
			if must {
				return fmt.Errorf("invalid postalcode")
			}
			return nil
		}
		if *postalCode == "" {
			return fmt.Errorf("invalid postalcode")
		}
		h.PostalCode = postalCode
		return nil
	}
}

func WithAge(age *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Age = age
		return nil
	}
}

func WithBirthday(birthday *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Birthday = birthday
		return nil
	}
}

func WithAvatar(avatar *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if avatar == nil {
			if must {
				return fmt.Errorf("invalid avatar")
			}
			return nil
		}
		if *avatar == "" {
			return fmt.Errorf("invalid avatar")
		}
		h.Avatar = avatar
		return nil
	}
}

func WithOrganization(organization *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if organization == nil {
			if must {
				return fmt.Errorf("invalid organization")
			}
			return nil
		}
		if *organization == "" {
			return fmt.Errorf("invalid organization")
		}
		h.Organization = organization
		return nil
	}
}

func WithFirstName(firstName *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if firstName == nil {
			if must {
				return fmt.Errorf("invalid firstname")
			}
			return nil
		}
		if *firstName == "" {
			return fmt.Errorf("invalid firstname")
		}
		h.FirstName = firstName
		return nil
	}
}

func WithLastName(lastName *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if lastName == nil {
			if must {
				return fmt.Errorf("invalid lastname")
			}
			return nil
		}
		if *lastName == "" {
			return fmt.Errorf("invalid lastname")
		}
		h.LastName = lastName
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
		if *idNumber == "" {
			return fmt.Errorf("invalid idnumber")
		}
		h.IDNumber = idNumber
		return nil
	}
}

func WithSigninVerifyType(verifyType *basetypes.SignMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if verifyType == nil {
			if must {
				return fmt.Errorf("invalid signinverifytype")
			}
			return nil
		}
		switch *verifyType {
		case basetypes.SignMethod_Email:
		case basetypes.SignMethod_Mobile:
		case basetypes.SignMethod_Google:
		default:
			return fmt.Errorf("invalid signverifytype")
		}
		h.SigninVerifyType = verifyType
		return nil
	}
}

func WithKolConfirmed(confirmed *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.KolConfirmed = confirmed
		return nil
	}
}

func WithSelectedLangID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid selectedlangid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.SelectedLangID = id
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

func WithShouldUpdateCache(update bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ShouldUpdateCache = update
		return nil
	}
}

func WithThirdPartyID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid thirdpartyid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ThirdPartyID = id
		return nil
	}
}
