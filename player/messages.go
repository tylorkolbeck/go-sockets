package player

type JoinMsg struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Player Player `json:"player"`
}

type PlayerLeaveMsg struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type PlayerWsMsg struct {
	ID  string       `json:"id"`
	Msg InputMapping `json:"msg"`
}
