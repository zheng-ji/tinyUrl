package main

import (
    "tinyurl/global"
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)

var (
	gitBranch, gitCommit string
	cpuprofile           = flag.Bool("cpuprofile", false, "CPU性能测试")
	showVersion          = flag.Bool("version", false, "显示当前版本号")
	configFile           = flag.String("c", "etc/server.yml", "配置文件路径，默认etc/server.yml")
)

func main() {

	flag.Parse()

	if *showVersion {
		fmt.Printf("Branch: %s\nCommit: %s\n", gitBranch, gitCommit)
		os.Exit(0)
	}

	if err := global.ParseConfigFile(*configFile); err != nil {
		panic(err)
	}

	if *cpuprofile {
		go func() {
			http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", global.GlobalConfig.MonitAddr), nil)
		}()
	}

	// Init Logger
	global.InitLogger()

	// GRPC Connection
	global.NewRedisPool(global.GlobalConfig.Redis)

	// DB Connection
	global.InitDBConnection()
	defer global.GDb.Close()

	router.InitRouter()
	svr := &http.Server{
		Addr:    global.GlobalConfig.Bind,
		Handler: router.Engine,
	}
	go func() {
		if err := svr.ListenAndServe(); err != nil {
			global.Runlogger.Infof("listen: %s\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	global.Runlogger.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Runlogger.Fatalf("Server Shutdown: %v", err)
	}
}
