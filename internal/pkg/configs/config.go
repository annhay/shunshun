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
	Zap struct {
		LogDir      string `json:"logDir"`
		MaxAge      int    `json:"maxAge"`
		Compress    bool   `json:"compress"`
		Level       string `json:"level"`
		Development bool   `json:"development"`
	}
	Huyi struct {
		APIID  string `json:"APIID,omitempty"`
		APIKEY string `json:"APIKEY,omitempty"`
	}
	AliYun struct {
		AccessKeyID     string `json:"accessKeyID"`
		AccessKeySecret string `json:"accessKeySecret"`
	}
}
