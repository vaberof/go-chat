package postgres

type PostgresDatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}
