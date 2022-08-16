package role

import (
	"context"

	approleusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracerrole "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/role"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/approleuser"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateRole(ctx context.Context, in *role.CreateRoleRequest) (*role.CreateRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracerrole.TraceCreate(span, in)

	err = validate(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("CreateRole", "err", err)
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "CreateAppRole")

	appRole, err := approlemgrcli.CreateAppRole(ctx, &approle.AppRoleReq{
		AppID:       &in.AppID,
		CreatedBy:   &in.UserID,
		Role:        &in.RoleName,
		Description: &in.Description,
		Default:     &in.Default,
	})
	if err != nil {
		logger.Sugar().Errorw("CreateRole", "err", err)
		return &role.CreateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	resp, err := rolemwcli.GetRole(ctx, appRole.ID)
	if err != nil {
		logger.Sugar().Errorw("CreateRole", "err", err)
		return &role.CreateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.CreateRoleResponse{
		Info: resp,
	}, nil
}

func (s *Server) CreateAppRole(ctx context.Context, in *role.CreateAppRoleRequest) (*role.CreateAppRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("TargetAppID", in.GetTargetAppID()))

	check := &role.CreateRoleRequest{
		AppID:       in.TargetAppID,
		UserID:      in.UserID,
		RoleName:    in.RoleName,
		Default:     in.Default,
		Description: in.Description,
	}
	span = tracerrole.TraceCreate(span, check)

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("CreateAppRole", "TargetAppID", in.GetTargetAppID(), "err", err)
		return &role.CreateAppRoleResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	err = validate(ctx, check)
	if err != nil {
		logger.Sugar().Errorw("CreateAppRole", "err", err)
		return &role.CreateAppRoleResponse{}, err
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "CreateAppRole")

	appRole, err := approlemgrcli.CreateAppRole(ctx, &approle.AppRoleReq{
		AppID:       &in.TargetAppID,
		CreatedBy:   &in.UserID,
		Role:        &in.RoleName,
		Description: &in.Description,
		Default:     &in.Default,
	})
	if err != nil {
		logger.Sugar().Errorw("CreateAppRole", "err", err)
		return &role.CreateAppRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	resp, err := rolemwcli.GetRole(ctx, appRole.ID)
	if err != nil {
		logger.Sugar().Errorw("CreateAppRole", "err", err)
		return &role.CreateAppRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.CreateAppRoleResponse{
		Info: resp,
	}, nil
}

func (s *Server) CreateRoleUser(ctx context.Context, in *role.CreateRoleUserRequest) (*role.CreateRoleUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateRoleUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in.GetInfo())

	err = validateRoleUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateRoleUser", "err", err)
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "CreateAppRoleUser")

	roleUser, err := approleusermgrcli.CreateAppRoleUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateRoleUser", "err", err)
		return &role.CreateRoleUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	resp, err := rolemwcli.GetRoleUser(ctx, roleUser.ID)
	if err != nil {
		logger.Sugar().Errorw("CreateRoleUser", "err", err)
		return &role.CreateRoleUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.CreateRoleUserResponse{
		Info: resp,
	}, nil
}

func (s *Server) CreateAppRoleUser(ctx context.Context, in *role.CreateAppRoleUserRequest) (*role.CreateAppRoleUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateRoleUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(
		attribute.String("TargetAppID", in.GetTargetAppID()),
		attribute.String("TargetUserID", in.GetTargetUserID()),
	)
	span = tracer.Trace(span, in.GetInfo())

	err = validateRoleUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateRoleUser", "err", err)
		return nil, err
	}

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("CreateRoleUser", "TargetAppID", in.GetTargetAppID(), "err", err)
		return &role.CreateAppRoleUserResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		logger.Sugar().Errorw("CreateRoleUser", "TargetUserID", in.GetTargetUserID(), "err", err)
		return &role.CreateAppRoleUserResponse{}, status.Error(codes.InvalidArgument, "TargetUserID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "CreateAppRoleUser")

	roleUser, err := approleusermgrcli.CreateAppRoleUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateRoleUser", "err", err)
		return &role.CreateAppRoleUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	resp, err := rolemwcli.GetRoleUser(ctx, roleUser.ID)
	if err != nil {
		logger.Sugar().Errorw("CreateRoleUser", "err", err)
		return &role.CreateAppRoleUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.CreateAppRoleUserResponse{
		Info: resp,
	}, nil
}
