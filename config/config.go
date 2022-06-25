package config

type Config struct {
	Field1 string `mapstructure:"field1"`
	Field2 string `mapstructure:"field2"`
	Logrus Logrus `mapstructure:"logrus"`
}

type Logrus struct {
	LogLvl int    `mapstructure:"log_level"`
	ToFile bool   `mapstructure:"to_file"`
	ToJson bool   `mapstructure:"to_json"`
	LogDir string `mapstructure:"log_dir"`
}
