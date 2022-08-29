//nolint:nolintlint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	kycent "github.com/NpoolPlatform/kyc-management/pkg/db/ent"
	kycconstant "github.com/NpoolPlatform/kyc-management/pkg/message/const"
	"github.com/google/uuid"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/NpoolPlatform/appuser-manager/pkg/db"
	"github.com/NpoolPlatform/appuser-manager/pkg/db/ent"
	authingent "github.com/NpoolPlatform/authing-gateway/pkg/db/ent"
	authconstant "github.com/NpoolPlatform/authing-gateway/pkg/message/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	kycpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
	reviewent "github.com/NpoolPlatform/review-service/pkg/db/ent"
	reviewtb "github.com/NpoolPlatform/review-service/pkg/db/ent/review"
	reviewconstant "github.com/NpoolPlatform/review-service/pkg/message/const"
)

func Migrate(ctx context.Context) error {
	err := migrationAuthingGateway(ctx)
	if err != nil {
		return err
	}
	return migrationKyc(ctx)
}

const (
	keyUsername = "username"
	keyPassword = "password"
	keyDBName   = "database_name"
	maxOpen     = 10
	maxIdle     = 10
	MaxLife     = 3
)

func dsn(hostname string) (string, error) {
	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	dbname := config.GetStringValueWithNameSpace(hostname, keyDBName)

	svc, err := config.PeekService(constant.MysqlServiceName)
	if err != nil {
		logger.Sugar().Warnw("dsb", "error", err)
		return "", err
	}

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true",
		username, password,
		svc.Address,
		svc.Port,
		dbname,
	), nil
}

func open(hostname string) (conn *sql.DB, err error) {
	hdsn, err := dsn(hostname)
	if err != nil {
		return nil, err
	}

	conn, err = sql.Open("mysql", hdsn)
	if err != nil {
		return nil, err
	}

	// https://github.com/go-sql-driver/mysql
	// See "Important settings" section.

	conn.SetConnMaxLifetime(time.Minute * MaxLife)
	conn.SetMaxOpenConns(maxOpen)
	conn.SetMaxIdleConns(maxIdle)

	return conn, nil
}

func migrationAuthingGateway(ctx context.Context) (err error) {
	cli, err := db.Client()
	if err != nil {
		return err
	}

	authInfos, err := cli.Auth.Query().Limit(1).All(ctx)
	if err != nil {
		return err
	}

	if len(authInfos) > 0 {
		return nil
	}

	logger.Sugar().Infow("Migrate", "Start", "...")
	defer func() {
		logger.Sugar().Infow("Migrate", "Done", "...", "error", err)
	}()

	auth, err := open(authconstant.ServiceName)
	if err != nil {
		return err
	}

	defer auth.Close()

	authCli := authingent.NewClient(authingent.Driver(entsql.OpenDB(dialect.MySQL, auth)))

	appAuths, err := authCli.
		AppAuth.
		Query().
		All(ctx)
	if err != nil {
		return err
	}

	appRoleAuths, err := authCli.
		AppRoleAuth.
		Query().
		All(ctx)
	if err != nil {
		return err
	}

	appUserAuths, err := authCli.
		AppUserAuth.
		Query().
		All(ctx)
	if err != nil {
		return err
	}

	auths := []*ent.Auth{}

	for _, val := range appAuths {
		auths = append(auths, &ent.Auth{
			AppID:     val.AppID,
			Resource:  val.Resource,
			Method:    val.Method,
			CreatedAt: val.CreateAt,
			UpdatedAt: val.UpdateAt,
		})
	}

	for _, val := range appRoleAuths {
		auths = append(auths, &ent.Auth{
			AppID:     val.AppID,
			RoleID:    val.RoleID,
			Resource:  val.Resource,
			Method:    val.Method,
			CreatedAt: val.CreateAt,
			UpdatedAt: val.UpdateAt,
		})
	}

	for _, val := range appUserAuths {
		auths = append(auths, &ent.Auth{
			AppID:     val.AppID,
			UserID:    val.UserID,
			Resource:  val.Resource,
			Method:    val.Method,
			CreatedAt: val.CreateAt,
			UpdatedAt: val.UpdateAt,
		})
	}

	err = db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		bulk := make([]*ent.AuthCreate, len(auths))
		for i, val := range auths {
			bulk[i] = tx.Auth.
				Create().
				SetAppID(val.AppID).
				SetRoleID(val.RoleID).
				SetUserID(val.UserID).
				SetResource(val.Resource).
				SetMethod(val.Method).
				SetCreatedAt(val.CreatedAt).
				SetUpdatedAt(val.UpdatedAt)
		}
		_, err = tx.Auth.CreateBulk(bulk...).Save(_ctx)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

//nolint:gocyclo
func migrationKyc(ctx context.Context) (err error) {
	cli, err := db.Client()
	if err != nil {
		return err
	}

	newKycInfos, err := cli.Kyc.Query().Limit(1).All(ctx)
	if err != nil {
		return err
	}

	if len(newKycInfos) > 0 {
		return nil
	}

	logger.Sugar().Infow("Migrate", "Start", "...")
	defer func() {
		logger.Sugar().Infow("Migrate", "Done", "...", "error", err)
	}()

	kyc, err := open(kycconstant.ServiceName)
	if err != nil {
		return err
	}

	defer kyc.Close()
	kycCli := kycent.NewClient(kycent.Driver(entsql.OpenDB(dialect.MySQL, kyc)))

	review, err := open(reviewconstant.ServiceName)
	if err != nil {
		return err
	}

	defer review.Close()
	reviewCli := reviewent.NewClient(reviewent.Driver(entsql.OpenDB(dialect.MySQL, review)))

	offset := 0
	limit := 1000

	newKyc := []*ent.Kyc{}

	for {
		infos, err := kycCli.Kyc.Query().Offset(offset).Limit(limit).All(ctx)
		if err != nil {
			return err
		}

		if len(infos) == 0 {
			break
		}

		kycIDs := make([]interface{}, len(infos))
		for _, val := range infos {
			kycIDs = append(kycIDs, val.ID.String())
		}

		type Review struct {
			ID          string `json:"id"`
			ObjectID    string `json:"object_id"`
			SubObjectID string `json:"sub_object_id"`
			ReviewState string `json:"state"`
		}

		reviewInfos := []Review{}
		err = reviewCli.Review.
			Query().Select(
			reviewtb.FieldID,
			reviewtb.FieldState,
			reviewtb.FieldObjectID,
		).Modify(func(s *entsql.Selector) {
			subQuery := entsql.Select(
				entsql.As(reviewtb.FieldObjectID, "sub_object_id"),
				entsql.As(entsql.Max(reviewtb.FieldUpdateAt), reviewtb.FieldUpdateAt)).
				From(
					entsql.Table(reviewtb.Table),
				).
				Where(
					entsql.In(reviewtb.FieldObjectID, kycIDs...)).
				GroupBy(reviewtb.FieldObjectID)

			s.
				RightJoin(subQuery).
				On(
					s.C(reviewtb.FieldObjectID),
					subQuery.C("sub_object_id"),
				).On(
				s.C(reviewtb.FieldUpdateAt),
				subQuery.C(reviewtb.FieldUpdateAt),
			)
		}).Scan(ctx, &reviewInfos)
		if err != nil {
			return err
		}

		for _, kycInfo := range infos {
			reviewID := uuid.UUID{}.String()
			state := kycpb.KycState_DefaultState
			for _, info := range reviewInfos {
				if kycInfo.ID.String() != info.ObjectID {
					continue
				}
				reviewID = info.ID
				switch info.ReviewState {
				case "wait":
					state = kycpb.KycState_Reviewing
				case "approved":
					state = kycpb.KycState_Approved
				case "rejected":
					state = kycpb.KycState_Rejected
				}
			}

			newKycReviewID, err := uuid.Parse(reviewID)
			if err != nil {
				return err
			}

			cardType := kycpb.KycDocumentType_DefaultKycDocumentType
			switch kycInfo.CardType {
			case "passport":
				cardType = kycpb.KycDocumentType_Passport
			case "id-card":
				cardType = kycpb.KycDocumentType_IDCard
			}

			newKyc = append(newKyc, &ent.Kyc{
				ID:           kycInfo.ID,
				CreatedAt:    kycInfo.CreateAt,
				UpdatedAt:    kycInfo.UpdateAt,
				DeletedAt:    0,
				AppID:        kycInfo.AppID,
				UserID:       kycInfo.UserID,
				DocumentType: cardType.String(),
				IDNumber:     kycInfo.CardID,
				FrontImg:     kycInfo.FrontCardImg,
				BackImg:      kycInfo.BackCardImg,
				SelfieImg:    kycInfo.UserHandingCardImg,
				EntityType:   kycpb.KycEntityType_Individual.String(),
				ReviewID:     newKycReviewID,
				State:        state.String(),
			})
		}
		offset += limit
	}

	err = db.WithTx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		bulk := make([]*ent.KycCreate, len(newKyc))
		for i, info := range newKyc {
			bulk[i] = tx.Kyc.
				Create().
				SetID(info.ID).
				SetCreatedAt(info.CreatedAt).
				SetUpdatedAt(info.UpdatedAt).
				SetAppID(info.AppID).
				SetUserID(info.UserID).
				SetDocumentType(info.DocumentType).
				SetIDNumber(info.IDNumber).
				SetFrontImg(info.FrontImg).
				SetBackImg(info.BackImg).
				SetSelfieImg(info.SelfieImg).
				SetEntityType(info.EntityType).
				SetReviewID(info.ReviewID).
				SetState(info.State)
		}
		_, err = tx.Kyc.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}
