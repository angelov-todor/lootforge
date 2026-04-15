package strategies

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/angelov-todor/lootforge/core/internal/models"
)

type WeightedLuckStrategy struct {
	Increment int
	Decrement int
}

func (s *WeightedLuckStrategy) Roll(participants []*models.Member) (string, error) {
	var pool []string
	for _, m := range participants {
		weight := m.Luck + 1 // base weight of 1 so luck-0 members can still win
		if weight < 1 {
			weight = 1
		}
		for i := 0; i < weight; i++ {
			pool = append(pool, m.ID)
		}
	}
	if len(pool) == 0 {
		return "", errors.New("no eligible participants")
	}

	// Fisher-Yates shuffle with crypto/rand
	for i := len(pool) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		pool[i], pool[j.Int64()] = pool[j.Int64()], pool[i]
	}

	idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
	if err != nil {
		return "", err
	}
	return pool[idx.Int64()], nil
}

func (s *WeightedLuckStrategy) AdjustAfterRoll(winner *models.Member, losers []*models.Member) error {
	winner.Luck -= s.Decrement
	if winner.Luck < 0 {
		winner.Luck = 0
	}
	for _, l := range losers {
		l.Luck += s.Increment
	}
	return nil
}
