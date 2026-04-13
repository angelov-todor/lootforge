package models

import "time"

type User struct {
	ID          string    `json:"id" firestore:"-"`
	Email       string    `json:"email" firestore:"email"`
	DisplayName string    `json:"displayName" firestore:"displayName"`
	PhotoURL    string    `json:"photoURL" firestore:"photoURL"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
	LastLoginAt time.Time `json:"lastLoginAt" firestore:"lastLoginAt"`
}
