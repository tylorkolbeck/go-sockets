package config

type Config struct {
	Server ServerConfig `json:"server"`
	Game   GameConfig   `json:"game"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type GameConfig struct {
	TickRate    int     `json:"tick_rate"`
	WorldWidth  float64 `json:"world_width"`
	WorldHeight float64 `json:"world_height"`
	MaxPlayers  int     `json:"max_players"`
}
