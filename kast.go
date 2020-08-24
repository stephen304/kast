package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stephen304/kast/internal"
	"github.com/stephen304/kast/internal/modules/backdrop"
	"github.com/stephen304/kast/internal/modules/media"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	r := gin.Default()
	display := internal.NewDisplayMutex()

	backdrop := backdrop.New(r.Group("/backdrop"), display)
	media := media.New(r.Group("/media"), display)

	r.POST("/stop", func(c *gin.Context) {
		display.Assign(backdrop)
	})

	r.POST("/reboot", func(c *gin.Context) {
		exec.Command("reboot").Run()
	})

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("API closed")
	log.Println("Shutting down modules...")
	backdrop.Stop()
	media.Stop()
}
