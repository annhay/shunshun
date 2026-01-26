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
	Amap struct {
		APIKey    string `json:"apiKey" yaml:"apiKey"`
		SecretKey string `json:"secretKey" yaml:"secretKey"`
		BaseURL   string `json:"baseURL" yaml:"baseURL"`
	}
	Tongyi struct {
		APIKey  string `json:"apiKey" yaml:"apiKey"`
		BaseURL string `json:"baseURL" yaml:"baseURL"`
		Model   string `json:"model" yaml:"model"`
	}
	AES struct {
		SecretKey string `json:"secretKey" yaml:"secretKey"`
	}
	RabbitMQ struct {
		Host       string `json:"host" yaml:"host"`
		Port       int    `json:"port" yaml:"port"`
		User       string `json:"user" yaml:"user"`
		Password   string `json:"password" yaml:"password"`
		VHost      string `json:"vHost" yaml:"vHost"`
		Exchange   string `json:"exchange" yaml:"exchange"`
		Queue      string `json:"queue" yaml:"queue"`
		RoutingKey string `json:"routingKey" yaml:"routingKey"`
	}
}
