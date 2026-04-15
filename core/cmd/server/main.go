package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/router"
	"github.com/angelov-todor/lootforge/core/internal/store"
)

func main() {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = "lootforge-dev"
	}

	fsClient, err := store.NewFirestoreClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer fsClient.Close()

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID})
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firebase Auth client: %v", err)
	}

	rollStore := store.NewFirestoreRollStore(fsClient)
	deps := router.Deps{
		Verifier:    auth.NewFirebaseTokenVerifier(authClient),
		UserStore:   store.NewFirestoreUserStore(fsClient),
		GroupStore:  store.NewFirestoreGroupStore(fsClient),
		RoleStore:   store.NewFirestoreRoleStore(fsClient),
		MemberStore: store.NewFirestoreMemberStore(fsClient),
		RollStore:   rollStore,
		InviteStore: store.NewFirestoreInviteStore(fsClient),
		RollTxStore: rollStore,
	}

	handler := router.New(ctx, deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("LootForge API starting on :%s", port)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited cleanly")
}
