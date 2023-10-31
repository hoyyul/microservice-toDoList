package main

import (
	"context"
	"go-micro-toDoList/app/gateway/routes"
	"go-micro-toDoList/app/gateway/rpc"
	"go-micro-toDoList/global"
	"go-micro-toDoList/setting"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// global setting
	setting.InitConfig()
	setting.InitLogger()

	// enable etcd discover services
	rpc.Init()

	// start service
	startListenAndServe()

}

func startListenAndServe() {
	r := routes.NewRouter()
	server := &http.Server{
		Addr:           global.Config.Server.Addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Println("gateway failed to run, err: ", err)
		}
	}()

	global.Logger.Println("gateway listen on: %s", global.Config.Server.Addr)
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-osSignals
	global.Logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//shut down http server after timeout
	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Println("closing http server gracefully failed: ", err)
	}

}
