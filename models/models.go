package models

type PlayerLeaveMsg struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Color struct {
	R int64 `json:"r"`
	G int64 `json:"g"`
	B int64 `json:"b"`
}
