package auth

import (
	"context"

	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/store"
)

func CanManageGroup(ctx context.Context, userID, groupID string, roleStore store.RoleStore) error {
	role, err := roleStore.GetRole(ctx, groupID, userID)
	if err != nil {
		return models.ErrForbidden
	}
	if role.Role != models.RoleOwner && role.Role != models.RoleAdmin {
		return models.ErrForbidden
	}
	return nil
}

func CanViewGroup(ctx context.Context, userID, groupID string, roleStore store.RoleStore) error {
	_, err := roleStore.GetRole(ctx, groupID, userID)
	if err != nil {
		return models.ErrForbidden
	}
	return nil
}

func IsOwner(ctx context.Context, userID, groupID string, roleStore store.RoleStore) error {
	role, err := roleStore.GetRole(ctx, groupID, userID)
	if err != nil {
		return models.ErrForbidden
	}
	if role.Role != models.RoleOwner {
		return models.ErrForbidden
	}
	return nil
}
