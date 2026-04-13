package strategies_test

import (
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/strategies"
)

func TestRoundRobin_HighestPriorityWins(t *testing.T) {
	s := &strategies.RoundRobinStrategy{}
	members := []*models.Member{
		{ID: "low", Priority: 1},
		{ID: "high", Priority: 5},
		{ID: "mid", Priority: 3},
	}
	winner, err := s.Roll(members)
	if err != nil {
		t.Fatal(err)
	}
	if winner != "high" {
		t.Errorf("expected high, got %s", winner)
	}
}

func TestRoundRobin_TieBreaking(t *testing.T) {
	s := &strategies.RoundRobinStrategy{}
	a := &models.Member{ID: "a", Priority: 5}
	b := &models.Member{ID: "b", Priority: 5}
	wins := map[string]int{}
	for i := 0; i < 1000; i++ {
		winner, _ := s.Roll([]*models.Member{a, b})
		wins[winner]++
	}
	if wins["a"] == 0 || wins["b"] == 0 {
		t.Errorf("tie breaking not random: a=%d, b=%d", wins["a"], wins["b"])
	}
}

func TestRoundRobin_AdjustResets(t *testing.T) {
	s := &strategies.RoundRobinStrategy{}
	winner := &models.Member{ID: "w", Priority: 5}
	losers := []*models.Member{{ID: "l1", Priority: 2}, {ID: "l2", Priority: 3}}
	s.AdjustAfterRoll(winner, losers)
	if winner.Priority != 0 {
		t.Errorf("winner priority: expected 0, got %d", winner.Priority)
	}
	if losers[0].Priority != 3 {
		t.Errorf("loser1 priority: expected 3, got %d", losers[0].Priority)
	}
	if losers[1].Priority != 4 {
		t.Errorf("loser2 priority: expected 4, got %d", losers[1].Priority)
	}
}
