//go:build integration

package store_test

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/store"
)

func newTestClient(t *testing.T) *firestore.Client {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST not set, skipping integration test")
	}
	client, err := firestore.NewClient(context.Background(), "lootforge-test")
	if err != nil {
		t.Fatalf("failed to create Firestore client: %v", err)
	}
	t.Cleanup(func() { client.Close() })
	return client
}

// --- UserStore ---

func TestFirestoreUserStore_UpsertAndGet(t *testing.T) {
	client := newTestClient(t)
	s := store.NewFirestoreUserStore(client)
	ctx := context.Background()

	user := &models.User{
		ID:          "test-user-" + time.Now().Format("150405"),
		Email:       "test@example.com",
		DisplayName: "Test User",
	}
	if err := s.UpsertUser(ctx, user); err != nil {
		t.Fatalf("UpsertUser: %v", err)
	}

	got, err := s.GetUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %q", got.Email)
	}
	if got.DisplayName != "Test User" {
		t.Errorf("expected displayName 'Test User', got %q", got.DisplayName)
	}
}

// --- GroupStore ---

func TestFirestoreGroupStore_CRUD(t *testing.T) {
	client := newTestClient(t)
	gs := store.NewFirestoreGroupStore(client)
	ctx := context.Background()

	group := &models.Group{
		Name:    "Integration Test Group",
		OwnerID: "owner-1",
		Strategy: models.StrategyConfig{
			Type: models.StrategyPureRandom,
		},
	}

	id, err := gs.CreateGroup(ctx, group)
	if err != nil {
		t.Fatalf("CreateGroup: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty ID")
	}

	got, err := gs.GetGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetGroup: %v", err)
	}
	if got.Name != "Integration Test Group" {
		t.Errorf("expected name 'Integration Test Group', got %q", got.Name)
	}

	got.Name = "Updated Group"
	if err := gs.UpdateGroup(ctx, got); err != nil {
		t.Fatalf("UpdateGroup: %v", err)
	}

	updated, err := gs.GetGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetGroup after update: %v", err)
	}
	if updated.Name != "Updated Group" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}

	if err := gs.DeleteGroup(ctx, id); err != nil {
		t.Fatalf("DeleteGroup: %v", err)
	}

	_, err = gs.GetGroup(ctx, id)
	if err != models.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

// --- MemberStore ---

func TestFirestoreMemberStore_CRUD(t *testing.T) {
	client := newTestClient(t)
	gs := store.NewFirestoreGroupStore(client)
	ms := store.NewFirestoreMemberStore(client)
	ctx := context.Background()

	gid, err := gs.CreateGroup(ctx, &models.Group{
		Name:     "Member Test Group",
		OwnerID:  "owner-1",
		Strategy: models.StrategyConfig{Type: models.StrategyPureRandom},
	})
	if err != nil {
		t.Fatalf("CreateGroup: %v", err)
	}
	t.Cleanup(func() { gs.DeleteGroup(ctx, gid) })

	member := &models.Member{Name: "Alice", Role: "member"}
	mid, err := ms.AddMember(ctx, gid, member)
	if err != nil {
		t.Fatalf("AddMember: %v", err)
	}

	got, err := ms.GetMember(ctx, gid, mid)
	if err != nil {
		t.Fatalf("GetMember: %v", err)
	}
	if got.Name != "Alice" {
		t.Errorf("expected Alice, got %q", got.Name)
	}

	got.Name = "Alice Updated"
	if err := ms.UpdateMember(ctx, gid, got); err != nil {
		t.Fatalf("UpdateMember: %v", err)
	}

	list, err := ms.ListMembers(ctx, gid)
	if err != nil {
		t.Fatalf("ListMembers: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 member, got %d", len(list))
	}

	if err := ms.DeleteMember(ctx, gid, mid); err != nil {
		t.Fatalf("DeleteMember: %v", err)
	}

	_, err = ms.GetMember(ctx, gid, mid)
	if err != models.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

// --- RoleStore ---

func TestFirestoreRoleStore_SetAndGet(t *testing.T) {
	client := newTestClient(t)
	gs := store.NewFirestoreGroupStore(client)
	rs := store.NewFirestoreRoleStore(client)
	ctx := context.Background()

	gid, _ := gs.CreateGroup(ctx, &models.Group{
		Name:     "Role Test Group",
		OwnerID:  "owner-1",
		Strategy: models.StrategyConfig{Type: models.StrategyPureRandom},
	})
	t.Cleanup(func() { gs.DeleteGroup(ctx, gid) })

	role := &models.GroupUserRole{UserID: "user-1", GroupID: gid, Role: models.RoleOwner}
	if err := rs.SetRole(ctx, role); err != nil {
		t.Fatalf("SetRole: %v", err)
	}

	got, err := rs.GetRole(ctx, gid, "user-1")
	if err != nil {
		t.Fatalf("GetRole: %v", err)
	}
	if got.Role != models.RoleOwner {
		t.Errorf("expected owner, got %q", got.Role)
	}

	roles, err := rs.ListRolesForGroup(ctx, gid)
	if err != nil {
		t.Fatalf("ListRolesForGroup: %v", err)
	}
	if len(roles) != 1 {
		t.Errorf("expected 1 role, got %d", len(roles))
	}
}

// --- RollStore ---

func TestFirestoreRollStore_CreateAndList(t *testing.T) {
	client := newTestClient(t)
	gs := store.NewFirestoreGroupStore(client)
	rollS := store.NewFirestoreRollStore(client)
	ctx := context.Background()

	gid, _ := gs.CreateGroup(ctx, &models.Group{
		Name:     "Roll Test Group",
		OwnerID:  "owner-1",
		Strategy: models.StrategyConfig{Type: models.StrategyPureRandom},
	})
	t.Cleanup(func() { gs.DeleteGroup(ctx, gid) })

	session := &models.RollSession{
		Strategy:       models.StrategyConfig{Type: models.StrategyPureRandom},
		ParticipantIDs: []string{"m1", "m2"},
		WinnerID:       "m1",
		Item:           "Sword",
	}

	rid, err := rollS.CreateRoll(ctx, gid, session)
	if err != nil {
		t.Fatalf("CreateRoll: %v", err)
	}
	if rid == "" {
		t.Fatal("expected non-empty roll ID")
	}

	if err := rollS.IncrementWinCount(ctx, gid, "m1"); err != nil {
		t.Fatalf("IncrementWinCount: %v", err)
	}

	rolls, _, err := rollS.ListRolls(ctx, gid, store.ListRollsOpts{Limit: 10})
	if err != nil {
		t.Fatalf("ListRolls: %v", err)
	}
	if len(rolls) != 1 {
		t.Fatalf("expected 1 roll, got %d", len(rolls))
	}
	if rolls[0].WinnerID != "m1" {
		t.Errorf("expected winner m1, got %q", rolls[0].WinnerID)
	}

	stats, err := rollS.GetRollStats(ctx, gid)
	if err != nil {
		t.Fatalf("GetRollStats: %v", err)
	}
	if stats["m1"] != 1 {
		t.Errorf("expected 1 win for m1, got %d", stats["m1"])
	}
}

// --- InviteStore ---

func TestFirestoreInviteStore_CRUD(t *testing.T) {
	client := newTestClient(t)
	is := store.NewFirestoreInviteStore(client)
	ctx := context.Background()

	invite := &models.Invite{
		Token:     "test-token-" + time.Now().Format("150405"),
		GroupID:   "group-1",
		CreatedBy: "user-1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := is.CreateInvite(ctx, invite); err != nil {
		t.Fatalf("CreateInvite: %v", err)
	}

	got, err := is.GetInvite(ctx, invite.Token)
	if err != nil {
		t.Fatalf("GetInvite: %v", err)
	}
	if got.GroupID != "group-1" {
		t.Errorf("expected groupID group-1, got %q", got.GroupID)
	}

	if err := is.DeleteInvite(ctx, invite.Token); err != nil {
		t.Fatalf("DeleteInvite: %v", err)
	}

	_, err = is.GetInvite(ctx, invite.Token)
	if err != models.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}
