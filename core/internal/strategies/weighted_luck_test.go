package strategies_test

import (
	"math"
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/strategies"
)

func newWeightedLuck() *strategies.WeightedLuckStrategy {
	return &strategies.WeightedLuckStrategy{Increment: 1, Decrement: 1}
}

func TestWeightedLuck_HigherLuckWinsMore(t *testing.T) {
	s := newWeightedLuck()
	wins := map[string]int{}
	for i := 0; i < 10000; i++ {
		h := &models.Member{ID: "high", Luck: 10}
		l := &models.Member{ID: "low", Luck: 5}
		winner, err := s.Roll([]*models.Member{h, l})
		if err != nil {
			t.Fatal(err)
		}
		wins[winner]++
	}
	ratio := float64(wins["high"]) / float64(wins["low"])
	if math.Abs(ratio-2.0) > 0.3 {
		t.Errorf("expected ratio ~2.0, got %.2f (high=%d, low=%d)", ratio, wins["high"], wins["low"])
	}
}

func TestWeightedLuck_AdjustWinner(t *testing.T) {
	s := newWeightedLuck()
	winner := &models.Member{ID: "w", Luck: 5}
	losers := []*models.Member{{ID: "l1", Luck: 3}, {ID: "l2", Luck: 7}}
	err := s.AdjustAfterRoll(winner, losers)
	if err != nil {
		t.Fatal(err)
	}
	if winner.Luck != 4 {
		t.Errorf("winner luck: expected 4, got %d", winner.Luck)
	}
	if losers[0].Luck != 4 {
		t.Errorf("loser1 luck: expected 4, got %d", losers[0].Luck)
	}
	if losers[1].Luck != 8 {
		t.Errorf("loser2 luck: expected 8, got %d", losers[1].Luck)
	}
}

func TestWeightedLuck_CustomIncrementDecrement(t *testing.T) {
	s := &strategies.WeightedLuckStrategy{Increment: 2, Decrement: 3}
	winner := &models.Member{ID: "w", Luck: 10}
	losers := []*models.Member{{ID: "l", Luck: 5}}
	s.AdjustAfterRoll(winner, losers)
	if winner.Luck != 7 {
		t.Errorf("winner luck: expected 7, got %d", winner.Luck)
	}
	if losers[0].Luck != 7 {
		t.Errorf("loser luck: expected 7, got %d", losers[0].Luck)
	}
}

func TestWeightedLuck_ZeroLuckExcluded(t *testing.T) {
	s := newWeightedLuck()
	zero := &models.Member{ID: "zero", Luck: 0}
	active := &models.Member{ID: "active", Luck: 5}
	winner, err := s.Roll([]*models.Member{zero, active})
	if err != nil {
		t.Fatal(err)
	}
	if winner != "active" {
		t.Errorf("expected active to always win, got %s", winner)
	}
}

func TestWeightedLuck_AllZeroLuckError(t *testing.T) {
	s := newWeightedLuck()
	_, err := s.Roll([]*models.Member{{ID: "a", Luck: 0}, {ID: "b", Luck: -1}})
	if err == nil {
		t.Error("expected error when all participants have zero or negative luck")
	}
}
