package strategies_test

import (
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
	// With base weight +1: high gets 11 entries, low gets 6 entries → ratio ~1.83
	ratio := float64(wins["high"]) / float64(wins["low"])
	if ratio < 1.4 || ratio > 2.3 {
		t.Errorf("expected ratio ~1.83, got %.2f (high=%d, low=%d)", ratio, wins["high"], wins["low"])
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

func TestWeightedLuck_ZeroLuckCanStillWin(t *testing.T) {
	s := newWeightedLuck()
	wins := map[string]int{}
	for i := 0; i < 1000; i++ {
		zero := &models.Member{ID: "zero", Luck: 0}
		active := &models.Member{ID: "active", Luck: 5}
		winner, err := s.Roll([]*models.Member{zero, active})
		if err != nil {
			t.Fatal(err)
		}
		wins[winner]++
	}
	// zero gets 1 entry, active gets 6 → zero should win ~14% of the time
	if wins["zero"] == 0 {
		t.Error("expected zero-luck member to win at least once in 1000 rolls")
	}
	if wins["active"] == 0 {
		t.Error("expected active member to win at least once in 1000 rolls")
	}
}

func TestWeightedLuck_AllZeroLuckStillWorks(t *testing.T) {
	s := newWeightedLuck()
	winner, err := s.Roll([]*models.Member{{ID: "a", Luck: 0}, {ID: "b", Luck: 0}})
	if err != nil {
		t.Fatalf("expected no error with zero-luck participants, got: %v", err)
	}
	if winner != "a" && winner != "b" {
		t.Errorf("unexpected winner: %s", winner)
	}
}

func TestWeightedLuck_LuckFloorAtZero(t *testing.T) {
	s := &strategies.WeightedLuckStrategy{Increment: 1, Decrement: 5}
	winner := &models.Member{ID: "w", Luck: 2}
	losers := []*models.Member{{ID: "l", Luck: 3}}
	s.AdjustAfterRoll(winner, losers)
	if winner.Luck != 0 {
		t.Errorf("winner luck: expected 0 (floored), got %d", winner.Luck)
	}
	if losers[0].Luck != 4 {
		t.Errorf("loser luck: expected 4, got %d", losers[0].Luck)
	}
}
