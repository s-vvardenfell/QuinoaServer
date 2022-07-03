package config

type Config struct {
	ServerHost     string `mapstructure:"server_addr"`
	ServerPort     string `mapstructure:"server_port"`
	ParserPort     string `mapstructure:"parser_port"`
	RedisServPort  string `mapstructure:"redis_serv_port"`
	WithReflection bool   `mapstructure:"with_reflection"`
	ExpTime        int    `mapstructure:"exp_time"`
	Logrus         Logrus `mapstructure:"logrus"`
}

type Logrus struct {
	LogLvl int    `mapstructure:"log_level"`
	ToFile bool   `mapstructure:"to_file"`
	ToJson bool   `mapstructure:"to_json"`
	LogDir string `mapstructure:"log_dir"`
}
