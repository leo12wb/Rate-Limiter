package entity

type Config struct {
	Limiter LimiterConfig `json:"limiter"`
}

type LimiterConfig struct {
	Database Database `json:"database"`
	Default  Default  `json:"default"`
	IPS      []IP     `json:"ips"`
	Tokens   []Token  `json:"tokens"`
}

type Database struct {
	InMemory bool `json:"inMemory"`
	Redis    bool `json:"redis"`
}

type Default struct {
	Requests int `json:"requests"`
	Every    int `json:"every"`
}

type IP struct {
	IP       string `json:"ip"`
	Requests int    `json:"requests"`
	Every    int    `json:"every"`
}

type Token struct {
	Token    string `json:"token"`
	Requests int    `json:"requests"`
	Every    int    `json:"every"`
}
