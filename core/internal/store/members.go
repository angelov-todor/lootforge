package store

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreMemberStore struct {
	client *firestore.Client
}

func NewFirestoreMemberStore(client *firestore.Client) *FirestoreMemberStore {
	return &FirestoreMemberStore{client: client}
}

func (s *FirestoreMemberStore) membersCol(groupID string) *firestore.CollectionRef {
	return s.client.Collection("groups").Doc(groupID).Collection("members")
}

func (s *FirestoreMemberStore) ListMembers(ctx context.Context, groupID string) ([]*models.Member, error) {
	docs, err := s.membersCol(groupID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	members := make([]*models.Member, 0, len(docs))
	for _, doc := range docs {
		var m models.Member
		if err := doc.DataTo(&m); err != nil {
			continue
		}
		m.ID = doc.Ref.ID
		m.GroupID = groupID
		members = append(members, &m)
	}
	return members, nil
}

func (s *FirestoreMemberStore) AddMember(ctx context.Context, groupID string, member *models.Member) (string, error) {
	member.GroupID = groupID
	member.CreatedAt = time.Now()
	ref, _, err := s.membersCol(groupID).Add(ctx, member)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

func (s *FirestoreMemberStore) GetMember(ctx context.Context, groupID, memberID string) (*models.Member, error) {
	doc, err := s.membersCol(groupID).Doc(memberID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	var m models.Member
	if err := doc.DataTo(&m); err != nil {
		return nil, err
	}
	m.ID = doc.Ref.ID
	m.GroupID = groupID
	return &m, nil
}

func (s *FirestoreMemberStore) UpdateMember(ctx context.Context, groupID string, member *models.Member) error {
	_, err := s.membersCol(groupID).Doc(member.ID).Set(ctx, member)
	return err
}

func (s *FirestoreMemberStore) DeleteMember(ctx context.Context, groupID, memberID string) error {
	_, err := s.membersCol(groupID).Doc(memberID).Delete(ctx)
	return err
}

func (s *FirestoreMemberStore) BatchUpdateMembers(ctx context.Context, groupID string, members []*models.Member) error {
	batch := s.client.Batch()
	for _, m := range members {
		ref := s.membersCol(groupID).Doc(m.ID)
		batch.Set(ref, m)
	}
	_, err := batch.Commit(ctx)
	return err
}
