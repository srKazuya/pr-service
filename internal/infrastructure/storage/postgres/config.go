package postgres

type Config struct {
	DSN            string
	Seed           bool
	MigrationsPath string
}
