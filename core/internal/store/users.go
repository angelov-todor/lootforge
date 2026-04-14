package store

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/angelov-todor/lootforge/core/internal/models"
)

type FirestoreUserStore struct {
	client *firestore.Client
}

func NewFirestoreUserStore(client *firestore.Client) *FirestoreUserStore {
	return &FirestoreUserStore{client: client}
}

func (s *FirestoreUserStore) UpsertUser(ctx context.Context, user *models.User) error {
	_, err := s.client.Collection("users").Doc(user.ID).Set(ctx, map[string]interface{}{
		"email":       user.Email,
		"displayName": user.DisplayName,
		"photoURL":    user.PhotoURL,
		"lastLoginAt": user.LastLoginAt,
	}, firestore.MergeAll)
	return err
}

func (s *FirestoreUserStore) GetUser(ctx context.Context, id string) (*models.User, error) {
	doc, err := s.client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	user.ID = doc.Ref.ID
	return &user, nil
}
