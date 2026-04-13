package strategies_test

import (
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/strategies"
)

func TestDKP_HighestPointsWins(t *testing.T) {
	s := &strategies.DKPStrategy{WinCost: 10}
	members := []*models.Member{
		{ID: "poor", Points: 5},
		{ID: "rich", Points: 50},
		{ID: "mid", Points: 20},
	}
	winner, err := s.Roll(members)
	if err != nil {
		t.Fatal(err)
	}
	if winner != "rich" {
		t.Errorf("expected rich, got %s", winner)
	}
}

func TestDKP_WinCostDeducted(t *testing.T) {
	s := &strategies.DKPStrategy{WinCost: 10}
	winner := &models.Member{ID: "w", Points: 50}
	losers := []*models.Member{{ID: "l", Points: 20}}
	s.AdjustAfterRoll(winner, losers)
	if winner.Points != 40 {
		t.Errorf("expected 40, got %d", winner.Points)
	}
	if losers[0].Points != 20 {
		t.Errorf("loser points should be unchanged, got %d", losers[0].Points)
	}
}

func TestDKP_TieBreaking(t *testing.T) {
	s := &strategies.DKPStrategy{WinCost: 5}
	a := &models.Member{ID: "a", Points: 30}
	b := &models.Member{ID: "b", Points: 30}
	wins := map[string]int{}
	for i := 0; i < 1000; i++ {
		winner, _ := s.Roll([]*models.Member{a, b})
		wins[winner]++
	}
	if wins["a"] == 0 || wins["b"] == 0 {
		t.Errorf("tie not broken randomly: a=%d, b=%d", wins["a"], wins["b"])
	}
}
