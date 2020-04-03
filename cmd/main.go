package main

import (
	"flag"
	"fmt"
	"hrserver/app"
	"hrserver/util/config"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	configFile = flag.String("c", "", "config file path")
)

type IServer interface {
	OnStart(cfg *config.Config) error
	Shutdown()
}

func interceptSignal(s IServer) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.Shutdown()
		os.Exit(0)
	}()
}

func main() {

	fmt.Println("begin start server")

	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := config.LoadConfigFile(*configFile)

	profPort := cfg.GetString("prof")

	go func() {
		fmt.Println(http.ListenAndServe(":"+profPort, nil))
	}()

	dls := &app.App{}

	interceptSignal(dls)

	if err := dls.OnStart(cfg); err != nil {
		fmt.Println(err)
	}
	return
}
