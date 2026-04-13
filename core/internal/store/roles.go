package store

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreRoleStore struct {
	client *firestore.Client
}

func NewFirestoreRoleStore(client *firestore.Client) *FirestoreRoleStore {
	return &FirestoreRoleStore{client: client}
}

func (s *FirestoreRoleStore) SetRole(ctx context.Context, role *models.GroupUserRole) error {
	_, err := s.client.Collection("groups").Doc(role.GroupID).
		Collection("roles").Doc(role.UserID).Set(ctx, role)
	return err
}

func (s *FirestoreRoleStore) GetRole(ctx context.Context, groupID, userID string) (*models.GroupUserRole, error) {
	doc, err := s.client.Collection("groups").Doc(groupID).
		Collection("roles").Doc(userID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	var role models.GroupUserRole
	if err := doc.DataTo(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *FirestoreRoleStore) ListRolesForGroup(ctx context.Context, groupID string) ([]*models.GroupUserRole, error) {
	docs, err := s.client.Collection("groups").Doc(groupID).
		Collection("roles").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var roles []*models.GroupUserRole
	for _, doc := range docs {
		var role models.GroupUserRole
		if err := doc.DataTo(&role); err != nil {
			continue
		}
		roles = append(roles, &role)
	}
	return roles, nil
}
