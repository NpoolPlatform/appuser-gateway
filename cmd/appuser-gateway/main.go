package main

import (
	"fmt"
	"os"

	ossconst "github.com/NpoolPlatform/go-service-framework/pkg/oss/const"
	kycconstant "github.com/NpoolPlatform/kyc-management/pkg/message/const"
	reviewconstant "github.com/NpoolPlatform/review-service/pkg/message/const"

	mysqlconst "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"

	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"

	"github.com/NpoolPlatform/go-service-framework/pkg/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	rabbitmqconst "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/const"
	redisconst "github.com/NpoolPlatform/go-service-framework/pkg/redis/const"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/message/const"
	authconst "github.com/NpoolPlatform/authing-gateway/pkg/message/const"

	cli "github.com/urfave/cli/v2"
)

func main() {
	commands := cli.Commands{
		runCmd,
	}

	description := fmt.Sprintf("my %v service cli\nFor help on any individual command run <%v COMMAND -h>\n",
		servicename.ServiceName, servicename.ServiceName)
	err := app.Init(
		servicename.ServiceName,
		description,
		"",
		"",
		"./",
		nil,
		commands,
		mysqlconst.MysqlServiceName,
		rabbitmqconst.RabbitMQServiceName,
		redisconst.RedisServiceName,
		appusermgrconst.ServiceName,
		authconst.ServiceName,
		authconst.ServiceName,
		ossconst.S3NameSpace,
		reviewconstant.ServiceName,
		kycconstant.ServiceName,
	)
	if err != nil {
		logger.Sugar().Errorf("fail to create %v: %v", servicename.ServiceName, err)
		return
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Sugar().Errorf("fail to run %v: %v", servicename.ServiceName, err)
	}
}
