package strategies

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/angelov-todor/lootforge/core/internal/models"
)

type RoundRobinStrategy struct{}

func (s *RoundRobinStrategy) Roll(participants []*models.Member) (string, error) {
	if len(participants) == 0 {
		return "", errors.New("no participants")
	}
	maxPriority := participants[0].Priority
	for _, m := range participants[1:] {
		if m.Priority > maxPriority {
			maxPriority = m.Priority
		}
	}
	var candidates []*models.Member
	for _, m := range participants {
		if m.Priority == maxPriority {
			candidates = append(candidates, m)
		}
	}
	if len(candidates) == 1 {
		return candidates[0].ID, nil
	}
	idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates))))
	if err != nil {
		return "", err
	}
	return candidates[idx.Int64()].ID, nil
}

func (s *RoundRobinStrategy) AdjustAfterRoll(winner *models.Member, losers []*models.Member) error {
	winner.Priority = 0
	for _, l := range losers {
		l.Priority++
	}
	return nil
}
