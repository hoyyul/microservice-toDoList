package main

import (
	"context"
	"micro-toDoList/app/gateway/routes"
	"micro-toDoList/app/gateway/rpc"
	"micro-toDoList/global"
	"micro-toDoList/setting"

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
		Addr:           ":" + global.Config.Server.Address,
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

	global.Logger.Printf("gateway listen on: %s", global.Config.Server.Address)
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
