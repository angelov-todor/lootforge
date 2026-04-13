package strategies

import "github.com/angelov-todor/lootforge/core/internal/models"

type Strategy interface {
	Roll(participants []*models.Member) (winnerID string, err error)
	AdjustAfterRoll(winner *models.Member, losers []*models.Member) error
}
