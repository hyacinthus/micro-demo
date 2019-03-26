package main

// Config is the global settings of this project
type Config struct {
	APP struct {
		Debug     bool   `default:"false"`
		Host      string `default:"0.0.0.0"`
		Port      string `default:"80"`
		PageSize  int    `default:"10"`
		JWTSecret string `default:"secret"`
		BaseURL   string `default:"https://api.example.com/"`
	}

	DB struct {
		Host     string `default:"mysql"`
		Port     string `default:"3306"`
		User     string `default:"root"`
		Password string `default:"root"`
		Name     string `default:"demo"`
		Lifetime int    `default:"3000"`
	}

	Redis struct {
		Host     string `default:"redis"`
		Port     string `default:"6379"`
		Password string
		DB       int `default:"0"`
	}

	NSQ struct {
		NsqdAddr       string `default:"nsqd:4150"`
		NsqLookupdAddr string `default:"nsqlookupd:4161"`
	}

	QCloud struct {
		SecretID  string
		SecretKey string
		Region    string `default:"ap-shanghai"`
		AppID     string `default:"1234567"`
	}
}
