package store

import (
	"context"

	"github.com/angelov-todor/lootforge/core/internal/models"
)

type UserStore interface {
	UpsertUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id string) (*models.User, error)
}

type GroupStore interface {
	CreateGroup(ctx context.Context, group *models.Group) (string, error)
	GetGroup(ctx context.Context, id string) (*models.Group, error)
	ListGroupsForUser(ctx context.Context, userID string) ([]*models.Group, error)
	UpdateGroup(ctx context.Context, group *models.Group) error
	DeleteGroup(ctx context.Context, id string) error
}

type RoleStore interface {
	SetRole(ctx context.Context, role *models.GroupUserRole) error
	GetRole(ctx context.Context, groupID, userID string) (*models.GroupUserRole, error)
	ListRolesForGroup(ctx context.Context, groupID string) ([]*models.GroupUserRole, error)
}

type MemberStore interface {
	ListMembers(ctx context.Context, groupID string) ([]*models.Member, error)
	AddMember(ctx context.Context, groupID string, member *models.Member) (string, error)
	GetMember(ctx context.Context, groupID, memberID string) (*models.Member, error)
	UpdateMember(ctx context.Context, groupID string, member *models.Member) error
	DeleteMember(ctx context.Context, groupID, memberID string) error
	BatchUpdateMembers(ctx context.Context, groupID string, members []*models.Member) error
}

type RollStore interface {
	CreateRoll(ctx context.Context, groupID string, roll *models.RollSession) (string, error)
	ListRolls(ctx context.Context, groupID string, opts ListRollsOpts) ([]*models.RollSession, string, error)
	GetRollStats(ctx context.Context, groupID string) (map[string]int, error)
	IncrementWinCount(ctx context.Context, groupID, memberID string) error
}

type ListRollsOpts struct {
	Cursor   string
	Limit    int
	MemberID string
	Item     string
}

type InviteStore interface {
	CreateInvite(ctx context.Context, invite *models.Invite) error
	GetInvite(ctx context.Context, token string) (*models.Invite, error)
	DeleteInvite(ctx context.Context, token string) error
}

// RollTxStore provides transactional roll persistence.
// Implementations atomically create the roll, update member stats,
// and increment the win counter in a single transaction.
type RollTxStore interface {
	ExecuteRollTx(ctx context.Context, groupID string, roll *models.RollSession, members []*models.Member) (string, error)
}
