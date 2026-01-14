//go:build go1.8

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shunshun/internal/api-gateway/router"
	"shunshun/internal/pkg/initialization"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	initialization.GatewayInit()
	
	// 注册API网关到Consul
	consul := initialization.NewConsul("14.103.173.254:8500")
	kv := initialization.ConsulKV{
		Name:    "api-gateway",
		Tags:    []string{"api-gateway"},
		Address: "127.0.0.1",
		Port:    8080,
	}
	serviceID, err := consul.RegisterServer(kv)
	if err != nil {
		log.Printf("failed to register service: %v", err)
	}
	defer consul.DeregisterServer(serviceID)
	
	r := router.LoadRouter()
	r.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	//端口设置
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
