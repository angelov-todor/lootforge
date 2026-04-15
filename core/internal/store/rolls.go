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
	doc, err := s.client.Collection("groups").Doc(groupID).Collection("meta").Doc("roll_stats").Get(ctx)
	if err != nil {
		// Fall back to scanning all rolls if stats doc doesn't exist yet
		return s.getRollStatsLegacy(ctx, groupID)
	}

	data := doc.Data()
	stats := make(map[string]int, len(data))
	for k, v := range data {
		if count, ok := v.(int64); ok {
			stats[k] = int(count)
		}
	}
	return stats, nil
}

func (s *FirestoreRollStore) getRollStatsLegacy(ctx context.Context, groupID string) (map[string]int, error) {
	docs, err := s.rollsCol(groupID).OrderBy("createdAt", firestore.Desc).Limit(500).Documents(ctx).GetAll()
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

func (s *FirestoreRollStore) IncrementWinCount(ctx context.Context, groupID, memberID string) error {
	ref := s.client.Collection("groups").Doc(groupID).Collection("meta").Doc("roll_stats")
	_, err := ref.Set(ctx, map[string]interface{}{
		memberID: firestore.Increment(1),
	}, firestore.MergeAll)
	return err
}

// ExecuteRollTx atomically creates a roll session, batch-updates member stats,
// and increments the winner's win counter in a single Firestore transaction.
func (s *FirestoreRollStore) ExecuteRollTx(ctx context.Context, groupID string, roll *models.RollSession, members []*models.Member) (string, error) {
	var rollID string
	err := s.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// 1. Create roll document
		rollRef := s.rollsCol(groupID).NewDoc()
		rollID = rollRef.ID
		roll.ID = rollID
		if err := tx.Set(rollRef, roll); err != nil {
			return err
		}

		// 2. Update all member stats
		for _, m := range members {
			memberRef := s.client.Collection("groups").Doc(groupID).Collection("members").Doc(m.ID)
			if err := tx.Set(memberRef, m); err != nil {
				return err
			}
		}

		// 3. Increment win counter
		statsRef := s.client.Collection("groups").Doc(groupID).Collection("meta").Doc("roll_stats")
		if err := tx.Set(statsRef, map[string]interface{}{
			roll.WinnerID: firestore.Increment(1),
		}, firestore.MergeAll); err != nil {
			return err
		}

		return nil
	})
	return rollID, err
}
