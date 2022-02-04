package main

import (
	"MyKubernetes/pkg/deploy"
	"MyKubernetes/pkg/events"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var mydeployment deploy.MyDeployment

// 初始化日志信息
var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer func(sugarLogger *zap.SugaredLogger) {
		err := sugarLogger.Sync()
		if err != nil {
			return
		}
	}(sugarLogger)
	//pods.GetPodsFromNamespace("kube-system", true)
	//deploy.GetDeployFromNamespace()
	//deploy.DeleteDeployment("web", "default")
	//mydeployment.Namespace = "default"
	//mydeployment.Name = "web"
	//mydeployment.PortName = "tcp-80"
	//mydeployment.ContainerName = "web"
	//mydeployment.ContainerPort = 80
	//mydeployment.Replicas = 2
	//mydeployment.Image = "nginx:1.19.3"
	//mydeployment.MatchLabel = "nginx"
	//mydeployment.ContainerName = "web"
	//deploy.CreateDeployment(mydeployment)
	//deploy.UpdateDeployment("web", "default", "nginx:1.19.3", 3)
	events.WatchEvents("default")

}

// InitLogger 初始化sugarLogger
func InitLogger() {
	//logger, _ := zap.NewProduction()
	//sugarLogger = logger.Sugar()
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger = logger.Sugar()
}

//
func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

//
func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./log.log")
	return zapcore.AddSync(file)
}
