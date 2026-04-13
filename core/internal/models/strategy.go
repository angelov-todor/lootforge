package models

type StrategyType string

const (
	StrategyWeightedLuck StrategyType = "weighted_luck"
	StrategyRoundRobin   StrategyType = "round_robin"
	StrategyDKP          StrategyType = "dkp"
	StrategyPureRandom   StrategyType = "pure_random"
)

type WeightedLuckConfig struct {
	Increment   int `json:"increment" firestore:"increment"`
	Decrement   int `json:"decrement" firestore:"decrement"`
	DefaultLuck int `json:"defaultLuck" firestore:"defaultLuck"`
}

type DKPConfig struct {
	WinCost int `json:"winCost" firestore:"winCost"`
}

type StrategyConfig struct {
	Type         StrategyType        `json:"type" firestore:"type"`
	WeightedLuck *WeightedLuckConfig `json:"weightedLuck,omitempty" firestore:"weightedLuck,omitempty"`
	DKP          *DKPConfig          `json:"dkp,omitempty" firestore:"dkp,omitempty"`
}
