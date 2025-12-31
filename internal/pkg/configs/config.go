package configs

type AppConfig struct {
	Mysql struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	}
	Redis struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	}
	HuYi struct {
		APIID  string `json:"APIID,omitempty"`
		APIKEY string `json:"APIKEY,omitempty"`
	}
}
