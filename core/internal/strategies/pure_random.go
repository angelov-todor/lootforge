package strategies

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/angelov-todor/lootforge/core/internal/models"
)

type PureRandomStrategy struct{}

func (s *PureRandomStrategy) Roll(participants []*models.Member) (string, error) {
	if len(participants) == 0 {
		return "", errors.New("no participants")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(participants))))
	if err != nil {
		return "", err
	}
	return participants[n.Int64()].ID, nil
}

func (s *PureRandomStrategy) AdjustAfterRoll(winner *models.Member, losers []*models.Member) error {
	return nil
}
