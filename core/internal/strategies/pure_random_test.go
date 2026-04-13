package strategies_test

import (
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/strategies"
)

func makeMember(id string) *models.Member {
	return &models.Member{ID: id, Name: id, Luck: 10, Priority: 0, Points: 0}
}

func TestPureRandom_EmptyParticipants(t *testing.T) {
	s := &strategies.PureRandomStrategy{}
	_, err := s.Roll(nil)
	if err == nil {
		t.Error("expected error for empty participants")
	}
}

func TestPureRandom_SingleParticipant(t *testing.T) {
	s := &strategies.PureRandomStrategy{}
	winner, err := s.Roll([]*models.Member{makeMember("alice")})
	if err != nil {
		t.Fatal(err)
	}
	if winner != "alice" {
		t.Errorf("expected alice, got %s", winner)
	}
}

func TestPureRandom_AllCanWin(t *testing.T) {
	s := &strategies.PureRandomStrategy{}
	members := []*models.Member{makeMember("a"), makeMember("b"), makeMember("c")}
	wins := make(map[string]int)
	for i := 0; i < 10000; i++ {
		winner, err := s.Roll(members)
		if err != nil {
			t.Fatal(err)
		}
		wins[winner]++
	}
	for _, m := range members {
		if wins[m.ID] == 0 {
			t.Errorf("member %s never won in 10000 rolls", m.ID)
		}
	}
}

func TestPureRandom_AdjustIsNoop(t *testing.T) {
	s := &strategies.PureRandomStrategy{}
	winner := makeMember("w")
	winner.Luck = 5
	losers := []*models.Member{makeMember("l")}
	losers[0].Luck = 3
	err := s.AdjustAfterRoll(winner, losers)
	if err != nil {
		t.Fatal(err)
	}
	if winner.Luck != 5 {
		t.Errorf("winner luck changed: %d", winner.Luck)
	}
	if losers[0].Luck != 3 {
		t.Errorf("loser luck changed: %d", losers[0].Luck)
	}
}
