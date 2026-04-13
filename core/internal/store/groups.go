package store

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreGroupStore struct {
	client *firestore.Client
}

func NewFirestoreGroupStore(client *firestore.Client) *FirestoreGroupStore {
	return &FirestoreGroupStore{client: client}
}

func (s *FirestoreGroupStore) CreateGroup(ctx context.Context, group *models.Group) (string, error) {
	group.CreatedAt = time.Now()
	ref, _, err := s.client.Collection("groups").Add(ctx, group)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

func (s *FirestoreGroupStore) GetGroup(ctx context.Context, id string) (*models.Group, error) {
	doc, err := s.client.Collection("groups").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	var group models.Group
	if err := doc.DataTo(&group); err != nil {
		return nil, err
	}
	group.ID = doc.Ref.ID
	return &group, nil
}

func (s *FirestoreGroupStore) ListGroupsForUser(ctx context.Context, userID string) ([]*models.Group, error) {
	roleDocs, err := s.client.CollectionGroup("roles").Where("userID", "==", userID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var groups []*models.Group
	for _, doc := range roleDocs {
		groupID := doc.Ref.Parent.Parent.ID
		group, err := s.GetGroup(ctx, groupID)
		if err != nil {
			continue
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (s *FirestoreGroupStore) UpdateGroup(ctx context.Context, group *models.Group) error {
	_, err := s.client.Collection("groups").Doc(group.ID).Set(ctx, group, firestore.MergeAll)
	return err
}

func (s *FirestoreGroupStore) DeleteGroup(ctx context.Context, id string) error {
	subcollections := []string{"members", "roles", "rolls"}
	for _, sub := range subcollections {
		docs, _ := s.client.Collection("groups").Doc(id).Collection(sub).Documents(ctx).GetAll()
		for _, doc := range docs {
			doc.Ref.Delete(ctx)
		}
	}
	_, err := s.client.Collection("groups").Doc(id).Delete(ctx)
	return err
}
