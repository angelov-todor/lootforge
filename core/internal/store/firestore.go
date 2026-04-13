package store

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func NewFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	log.Printf("Firestore client initialized for project: %s", projectID)
	return client, nil
}
