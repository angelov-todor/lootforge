package strategies_test

import (
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/strategies"
)

func TestFactory_WeightedLuck(t *testing.T) {
	cfg := models.StrategyConfig{
		Type:         models.StrategyWeightedLuck,
		WeightedLuck: &models.WeightedLuckConfig{Increment: 2, Decrement: 1, DefaultLuck: 10},
	}
	s, err := strategies.NewStrategy(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.(*strategies.WeightedLuckStrategy); !ok {
		t.Error("expected WeightedLuckStrategy")
	}
}

func TestFactory_RoundRobin(t *testing.T) {
	s, err := strategies.NewStrategy(models.StrategyConfig{Type: models.StrategyRoundRobin})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.(*strategies.RoundRobinStrategy); !ok {
		t.Error("expected RoundRobinStrategy")
	}
}

func TestFactory_DKP(t *testing.T) {
	cfg := models.StrategyConfig{
		Type: models.StrategyDKP,
		DKP:  &models.DKPConfig{WinCost: 15},
	}
	s, err := strategies.NewStrategy(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.(*strategies.DKPStrategy); !ok {
		t.Error("expected DKPStrategy")
	}
}

func TestFactory_PureRandom(t *testing.T) {
	s, err := strategies.NewStrategy(models.StrategyConfig{Type: models.StrategyPureRandom})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.(*strategies.PureRandomStrategy); !ok {
		t.Error("expected PureRandomStrategy")
	}
}

func TestFactory_Unknown(t *testing.T) {
	_, err := strategies.NewStrategy(models.StrategyConfig{Type: "unknown"})
	if err == nil {
		t.Error("expected error for unknown strategy type")
	}
}
