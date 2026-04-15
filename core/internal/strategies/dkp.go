package strategies

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/angelov-todor/lootforge/core/internal/models"
)

type DKPStrategy struct {
	WinCost int
}

func (s *DKPStrategy) Roll(participants []*models.Member) (string, error) {
	if len(participants) == 0 {
		return "", errors.New("no participants")
	}
	maxPoints := participants[0].Points
	for _, m := range participants[1:] {
		if m.Points > maxPoints {
			maxPoints = m.Points
		}
	}
	var candidates []*models.Member
	for _, m := range participants {
		if m.Points == maxPoints {
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

func (s *DKPStrategy) AdjustAfterRoll(winner *models.Member, losers []*models.Member) error {
	winner.Points -= s.WinCost
	if winner.Points < 0 {
		winner.Points = 0
	}
	return nil
}
