package models

import "time"

type Invite struct {
	Token     string    `json:"token" firestore:"-"`
	GroupID   string    `json:"groupID" firestore:"groupID"`
	CreatedBy string    `json:"createdBy" firestore:"createdBy"`
	ExpiresAt time.Time `json:"expiresAt" firestore:"expiresAt"`
}
