package kyc

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"
	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
)

func CreateKyc(
	ctx context.Context,
	appID, userID, frontImg, selfieImg string,
	idNumber, backImg *string,
	documentType kycmgrpb.KycDocumentType,
	entityType kycmgrpb.KycEntityType,
) (
	*npool.Kyc, error,
) {
	return nil, nil
}
