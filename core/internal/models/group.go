package models

import "time"

type Group struct {
	ID        string         `json:"id" firestore:"-"`
	Name      string         `json:"name" firestore:"name"`
	OwnerID   string         `json:"ownerID" firestore:"ownerID"`
	Strategy  StrategyConfig `json:"strategy" firestore:"strategy"`
	CreatedAt time.Time      `json:"createdAt" firestore:"createdAt"`
}

type GroupUserRole struct {
	UserID  string `json:"userID" firestore:"userID"`
	GroupID string `json:"groupID" firestore:"groupID"`
	Role    string `json:"role" firestore:"role"`
}

const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)
