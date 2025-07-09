package config

import (
	"fmt"
	"time"

	"github.com/knadh/koanf"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(config *koanf.Koanf, log *zerolog.Logger) *gorm.DB {
	username := config.String("database.username")
	password := config.String("database.password")
	host := config.String("database.host")
	port := config.String("database.port")
	database := config.String("database.name")
	idleConnection := config.Int("database.pool.idle")
	maxConnection := config.Int("database.pool.max")
	maxLifeTimeConnection := config.Int("database.pool.lifetime")

	if idleConnection == 0 {
		idleConnection = 20
	}
	if maxConnection == 0 {
		maxConnection = 100
	}
	if maxLifeTimeConnection == 0 {
		maxLifeTimeConnection = 30
	}
	fmt.Println(host, port, username, password, database)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&zerologWritter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Duration(maxLifeTimeConnection) * time.Minute)

	return db
}

type zerologWritter struct {
	Logger *zerolog.Logger
}

func (l *zerologWritter) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Logger.Trace().Msg(msg)
}
