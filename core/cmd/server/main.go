package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	deps := router.Deps{
		Verifier:    auth.NewFirebaseTokenVerifier(authClient),
		UserStore:   store.NewFirestoreUserStore(fsClient),
		GroupStore:  store.NewFirestoreGroupStore(fsClient),
		RoleStore:   store.NewFirestoreRoleStore(fsClient),
		MemberStore: store.NewFirestoreMemberStore(fsClient),
		RollStore:   store.NewFirestoreRollStore(fsClient),
		InviteStore: store.NewFirestoreInviteStore(fsClient),
	}

	handler := router.New(deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("LootForge API starting on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
