package config

type config struct {
	FeaturezHost     string `split_words:"true" default:"localhost"`
	FeaturezPort     string `split_words:"true" default:"1000"`
	RedisNoDB        bool   `split_words:"true" default:"false"`
	RedisHost        string `split_words:"true" default:"localhost"`
	RedisPort        string `split_words:"true" default:"6379"`
	PostgresHost     string `split_words:"true" default:"localhost"`
	PostgresPort     string `split_words:"true" default:"5432"`
	PostgresUser     string `split_words:"true" default:"gorm"`
	PostgresPassword string `split_words:"true" default:"gorm"`
	PostgresDBName   string `split_words:"true" default:"gorm"`
	APIRoute         string `split_words:"true" default:"/api"`
	DebugMode        bool   `split_words:"true" default:"false"`
}

var (
	Cfg config
)
