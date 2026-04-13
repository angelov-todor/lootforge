package models

import "time"

type RollSession struct {
	ID             string         `json:"id" firestore:"-"`
	GroupID        string         `json:"groupID" firestore:"groupID"`
	Strategy       StrategyConfig `json:"strategy" firestore:"strategy"`
	ParticipantIDs []string       `json:"participantIDs" firestore:"participantIDs"`
	WinnerID       string         `json:"winnerID" firestore:"winnerID"`
	Item           string         `json:"item" firestore:"item"`
	CreatedAt      time.Time      `json:"createdAt" firestore:"createdAt"`
}
