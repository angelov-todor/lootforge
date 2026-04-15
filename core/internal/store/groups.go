package store

import (
	"context"
	"fmt"
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
	if len(roleDocs) == 0 {
		return nil, nil
	}

	refs := make([]*firestore.DocumentRef, len(roleDocs))
	for i, doc := range roleDocs {
		refs[i] = s.client.Collection("groups").Doc(doc.Ref.Parent.Parent.ID)
	}
	snapshots, err := s.client.GetAll(ctx, refs)
	if err != nil {
		return nil, err
	}

	var groups []*models.Group
	for _, doc := range snapshots {
		if !doc.Exists() {
			continue
		}
		var group models.Group
		if err := doc.DataTo(&group); err != nil {
			continue
		}
		group.ID = doc.Ref.ID
		groups = append(groups, &group)
	}
	return groups, nil
}

func (s *FirestoreGroupStore) UpdateGroup(ctx context.Context, group *models.Group) error {
	_, err := s.client.Collection("groups").Doc(group.ID).Set(ctx, group)
	return err
}

func (s *FirestoreGroupStore) DeleteGroup(ctx context.Context, id string) error {
	subcollections := []string{"members", "roles", "rolls", "meta"}
	for _, sub := range subcollections {
		docs, err := s.client.Collection("groups").Doc(id).Collection(sub).Documents(ctx).GetAll()
		if err != nil {
			return fmt.Errorf("failed to list %s: %w", sub, err)
		}
		for i := 0; i < len(docs); i += 500 {
			batch := s.client.Batch()
			end := i + 500
			if end > len(docs) {
				end = len(docs)
			}
			for _, doc := range docs[i:end] {
				batch.Delete(doc.Ref)
			}
			if _, err := batch.Commit(ctx); err != nil {
				return fmt.Errorf("failed to delete %s batch: %w", sub, err)
			}
		}
	}

	// Delete invites for this group
	inviteDocs, err := s.client.Collection("invites").Where("groupID", "==", id).Documents(ctx).GetAll()
	if err == nil {
		batch := s.client.Batch()
		for _, doc := range inviteDocs {
			batch.Delete(doc.Ref)
		}
		if len(inviteDocs) > 0 {
			batch.Commit(ctx)
		}
	}

	_, err = s.client.Collection("groups").Doc(id).Delete(ctx)
	return err
}
