package main

import (
	"time"

	"github.com/devaartana/e01-oprec-rpl/internal/auth"
	"github.com/devaartana/e01-oprec-rpl/internal/db"
	"github.com/devaartana/e01-oprec-rpl/internal/env"
	"github.com/devaartana/e01-oprec-rpl/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	err := godotenv.Load(".env.local")
	if err != nil {
		logger.Info("No .env file found or error loading it: %v", err)
	}

	cfg := config{
		addr: env.GetString("GO_ADDR", "localhost:8000"),
		db: dbConfig{
			addr:              env.GetString("DB_ADDR", "mongodb://opet:pasipuri@localhost:27017"),
			maxOpenConnection: env.GetInt("MONGO_MAX_OPEN_CONNECTION", 30),
			maxIdleConnection: env.GetInt("MONGO_MAX_IDLE_CONNECTION", 30),
			maxIdleTime:       env.GetString("MONGO_MAX_IDLE_TIME", "15m"),
		},
		auth: authConfig{
			user:   env.GetString("AUTH_USER", "admin"),
			secret: env.GetString("AUTH_SECRET", "admin"),
			exp:    time.Hour * time.Duration(env.GetInt("AUTH_EXP", 72)),
			iss:    env.GetString("AUTH_ISS", "opet"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConnection,
		cfg.db.maxIdleConnection,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("database connection pool established")

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.secret,
		cfg.auth.iss,
		cfg.auth.iss,
	)

	store := store.NewStorage(db)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
