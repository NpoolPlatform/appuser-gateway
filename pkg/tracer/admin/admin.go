package admin

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	trace1 "go.opentelemetry.io/otel/trace"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
)

func trace(span trace1.Span, in *npool.CreateGenesisUserRequest, index int) trace1.Span {
	span.SetAttributes(
		attribute.String(fmt.Sprintf("ID.%v", index), in.GetTargetAppID()),
		attribute.String(fmt.Sprintf("Description.%v", index), in.GetPasswordHash()),
		attribute.String(fmt.Sprintf("CreatedBy.%v", index), in.GetEmailAddress()),
	)
	return span
}

func Trace(span trace1.Span, in *npool.CreateGenesisUserRequest) trace1.Span {
	return trace(span, in, 0)
}
