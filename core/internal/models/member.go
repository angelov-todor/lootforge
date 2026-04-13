package models

import "time"

type Member struct {
	ID        string    `json:"id" firestore:"-"`
	GroupID   string    `json:"groupID" firestore:"groupID"`
	Name      string    `json:"name" firestore:"name"`
	Role      string    `json:"role" firestore:"role"`
	Luck      int       `json:"luck" firestore:"luck"`
	Priority  int       `json:"priority" firestore:"priority"`
	Points    int       `json:"points" firestore:"points"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
}
