package user

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	trace1 "go.opentelemetry.io/otel/trace"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
)

func trace(span trace1.Span, in *npool.SignupRequest, index int) trace1.Span {
	span.SetAttributes(
		attribute.String(fmt.Sprintf("AppID.%v", index), in.GetAppID()),
		attribute.String(fmt.Sprintf("Username.%v", index), in.GetPasswordHash()),
		attribute.String(fmt.Sprintf("PasswordHash.%v", index), in.GetPasswordHash()),
		attribute.String(fmt.Sprintf("Account.%v", index), in.GetAccount()),
		attribute.String(fmt.Sprintf("AccountType.%v", index), in.GetAccountType().String()),
		attribute.String(fmt.Sprintf("VerificationCode.%v", index), in.GetVerificationCode()),
		attribute.String(fmt.Sprintf("InvitationCode.%v", index), in.GetInvitationCode()),
	)
	return span
}

func Trace(span trace1.Span, in *npool.SignupRequest) trace1.Span {
	return trace(span, in, 0)
}
