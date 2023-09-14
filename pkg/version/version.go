package version

import (
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	logger "github.com/NpoolPlatform/go-service-framework/pkg/logger"
	cv "github.com/NpoolPlatform/go-service-framework/pkg/version"
)

func Version() (*basetypes.VersionResponse, error) {
	info, err := cv.GetVersion()
	if err != nil {
		logger.Sugar().Errorf("get service version error: %+w", err)
		return nil, fmt.Errorf("get service version error: %w", err)
	}
	return &basetypes.VersionResponse{
		Info: info,
	}, nil
}
