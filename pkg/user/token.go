package user

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"

	"github.com/google/uuid"
)

type Metadata struct {
	AppID       uuid.UUID
	Account     string
	AccountType string
	UserID      uuid.UUID
	ClientIP    net.IP
	UserAgent   string
	User        *usermwpb.User
}

func MetadataFromContext(ctx context.Context) (*Metadata, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("fail get metadata")
	}

	clientIP := ""
	if forwards, ok := meta["x-forwarded-for"]; ok {
		if len(forwards) > 0 {
			clientIP = strings.Split(forwards[0], ",")[0]
		}
	}

	userAgent := ""
	if agents, ok := meta["grpcgateway-user-agent"]; ok {
		if len(agents) > 0 {
			userAgent = agents[0]
		}
	}

	return &Metadata{
		ClientIP:  net.ParseIP(clientIP),
		UserAgent: userAgent,
	}, nil
}

func (meta *Metadata) ToJWTClaims() jwt.MapClaims {
	claims := jwt.MapClaims{}

	claims["app_id"] = meta.AppID
	claims["user_id"] = meta.UserID
	claims["account"] = meta.Account
	claims["account_type"] = meta.AccountType
	claims["client_ip"] = meta.ClientIP
	claims["user_agent"] = meta.UserAgent

	return claims
}

func (meta *Metadata) ValidateJWTClaims(claims jwt.MapClaims) error {
	appID, ok := claims["app_id"]
	if !ok || appID.(string) != meta.AppID.String() {
		return fmt.Errorf("invalid app id")
	}
	userID, ok := claims["user_id"]
	if !ok || userID.(string) != meta.UserID.String() {
		return fmt.Errorf("invalid user id")
	}
	account, ok := claims["account"]
	if !ok || account.(string) != meta.Account {
		return fmt.Errorf("invalid account")
	}
	loginType, ok := claims["account_type"]
	if !ok || loginType.(string) != meta.AccountType {
		return fmt.Errorf("invalid account type")
	}
	clientIP, ok := claims["client_ip"]
	if !ok || clientIP.(string) != meta.ClientIP.String() {
		return fmt.Errorf("invalid client ip")
	}
	userAgent, ok := claims["user_agent"]
	if !ok || userAgent.(string) != meta.UserAgent {
		return fmt.Errorf("invalid user agent")
	}
	return nil
}

func createToken(meta *Metadata) (string, error) {
	tokenAccessSecret := os.Getenv("LOGIN_TOKEN_ACCESS_SECRET")
	if tokenAccessSecret == "" {
		return "", fmt.Errorf("invalid login token access secret")
	}

	claims := meta.ToJWTClaims()
	candidate := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := candidate.SignedString([]byte(tokenAccessSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func verifyToken(meta *Metadata, token string) error {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		tokenAccessSecret := os.Getenv("LOGIN_TOKEN_ACCESS_SECRET")
		if tokenAccessSecret == "" {
			return "", fmt.Errorf("invalid login token access secret")
		}
		return []byte(tokenAccessSecret), nil
	})
	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return fmt.Errorf("invalid token")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("type assertion fail of jwt token")
	}

	err = meta.ValidateJWTClaims(claims)
	if err != nil {
		return err
	}

	return nil
}
