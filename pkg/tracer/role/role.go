package user

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	trace1 "go.opentelemetry.io/otel/trace"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"
)

func traceCreate(span trace1.Span, in *npool.CreateRoleRequest, index int) trace1.Span {
	span.SetAttributes(
		attribute.String(fmt.Sprintf("AppID.%v", index), in.GetAppID()),
		attribute.String(fmt.Sprintf("UserID.%v", index), in.GetUserID()),
		attribute.Bool(fmt.Sprintf("Default.%v", index), in.GetDefault()),
		attribute.String(fmt.Sprintf("RoleName.%v", index), in.GetRoleName()),
	)
	return span
}

func TraceCreate(span trace1.Span, in *npool.CreateRoleRequest) trace1.Span {
	return traceCreate(span, in, 0)
}
