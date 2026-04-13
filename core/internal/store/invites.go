package store

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreInviteStore struct {
	client *firestore.Client
}

func NewFirestoreInviteStore(client *firestore.Client) *FirestoreInviteStore {
	return &FirestoreInviteStore{client: client}
}

func (s *FirestoreInviteStore) CreateInvite(ctx context.Context, invite *models.Invite) error {
	_, err := s.client.Collection("invites").Doc(invite.Token).Set(ctx, invite)
	return err
}

func (s *FirestoreInviteStore) GetInvite(ctx context.Context, token string) (*models.Invite, error) {
	doc, err := s.client.Collection("invites").Doc(token).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	var invite models.Invite
	if err := doc.DataTo(&invite); err != nil {
		return nil, err
	}
	invite.Token = doc.Ref.ID
	return &invite, nil
}

func (s *FirestoreInviteStore) DeleteInvite(ctx context.Context, token string) error {
	_, err := s.client.Collection("invites").Doc(token).Delete(ctx)
	return err
}
