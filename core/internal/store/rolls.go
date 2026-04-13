package store

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/angelov-todor/lootforge/core/internal/models"
)

type FirestoreRollStore struct {
	client *firestore.Client
}

func NewFirestoreRollStore(client *firestore.Client) *FirestoreRollStore {
	return &FirestoreRollStore{client: client}
}

func (s *FirestoreRollStore) rollsCol(groupID string) *firestore.CollectionRef {
	return s.client.Collection("groups").Doc(groupID).Collection("rolls")
}

func (s *FirestoreRollStore) CreateRoll(ctx context.Context, groupID string, roll *models.RollSession) (string, error) {
	roll.GroupID = groupID
	roll.CreatedAt = time.Now()
	ref, _, err := s.rollsCol(groupID).Add(ctx, roll)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

func (s *FirestoreRollStore) ListRolls(ctx context.Context, groupID string, opts ListRollsOpts) ([]*models.RollSession, string, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}

	q := s.rollsCol(groupID).OrderBy("createdAt", firestore.Desc).Limit(limit + 1)

	if opts.Cursor != "" {
		cursorDoc, err := s.rollsCol(groupID).Doc(opts.Cursor).Get(ctx)
		if err == nil {
			q = q.StartAfter(cursorDoc.Data()["createdAt"])
		}
	}

	if opts.MemberID != "" {
		q = q.Where("winnerID", "==", opts.MemberID)
	}

	if opts.Item != "" {
		q = q.Where("item", "==", opts.Item)
	}

	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(docs) > limit {
		nextCursor = docs[limit-1].Ref.ID
		docs = docs[:limit]
	}

	rolls := make([]*models.RollSession, 0, len(docs))
	for _, doc := range docs {
		var r models.RollSession
		if err := doc.DataTo(&r); err != nil {
			continue
		}
		r.ID = doc.Ref.ID
		r.GroupID = groupID
		rolls = append(rolls, &r)
	}
	return rolls, nextCursor, nil
}

func (s *FirestoreRollStore) GetRollStats(ctx context.Context, groupID string) (map[string]int, error) {
	docs, err := s.rollsCol(groupID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	stats := make(map[string]int)
	for _, doc := range docs {
		var r models.RollSession
		if err := doc.DataTo(&r); err != nil {
			continue
		}
		stats[r.WinnerID]++
	}
	return stats, nil
}
