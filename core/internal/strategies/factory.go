package strategies

import (
	"fmt"

	"github.com/angelov-todor/lootforge/core/internal/models"
)

func NewStrategy(cfg models.StrategyConfig) (Strategy, error) {
	switch cfg.Type {
	case models.StrategyWeightedLuck:
		inc, dec := 1, 1
		if cfg.WeightedLuck != nil {
			if cfg.WeightedLuck.Increment > 0 {
				inc = cfg.WeightedLuck.Increment
			}
			if cfg.WeightedLuck.Decrement > 0 {
				dec = cfg.WeightedLuck.Decrement
			}
		}
		return &WeightedLuckStrategy{Increment: inc, Decrement: dec}, nil
	case models.StrategyRoundRobin:
		return &RoundRobinStrategy{}, nil
	case models.StrategyDKP:
		cost := 10
		if cfg.DKP != nil && cfg.DKP.WinCost > 0 {
			cost = cfg.DKP.WinCost
		}
		return &DKPStrategy{WinCost: cost}, nil
	case models.StrategyPureRandom:
		return &PureRandomStrategy{}, nil
	default:
		return nil, fmt.Errorf("unknown strategy type: %s", cfg.Type)
	}
}
