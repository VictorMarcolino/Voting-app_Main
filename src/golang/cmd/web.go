package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"votingMicroservicesApp/pkg/handlers"
)

func web() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	route := gin.Default()
	route.POST("/create", handlers.CreateVoting)
	route.GET("/list", handlers.ListVoting)
	route.GET("/:sessionUuid", handlers.GetVotingSessionInfo)
	route.DELETE("/:sessionUuid", handlers.DeleteVotingSession)
	route.POST("/:sessionUuid/:candidateUuid", handlers.VoteOnCandidate)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: route,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
