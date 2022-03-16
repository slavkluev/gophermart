package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/slavkluev/gophermart/internal/app/handler"
	"github.com/slavkluev/gophermart/internal/app/middleware"
	"github.com/slavkluev/gophermart/internal/app/repository"
	"github.com/slavkluev/gophermart/internal/app/service"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey            string
	MigrationDir         string
}

func main() {
	cfg := parseVariables()

	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", cfg.MigrationDir), "gophermart", driver)
	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	userRepository := repository.CreateUserRepository(db)
	orderRepository := repository.CreateOrderRepository(db)
	withdrawalRepository := repository.CreateWithdrawalRepository(db)
	cookieAuthenticator := service.NewCookieAuthenticator([]byte(cfg.SecretKey))
	pointAccrualService := service.NewPointAccrualService(cfg.AccrualSystemAddress, orderRepository, userRepository)
	pointAccrualService.Start()
	authenticator := middleware.NewAuthenticator(cookieAuthenticator)

	mws := []handler.Middleware{
		middleware.GzipEncoder{},
		middleware.GzipDecoder{},
	}

	handler := handler.NewHandler(
		cfg.AccrualSystemAddress,
		userRepository,
		orderRepository,
		withdrawalRepository,
		cookieAuthenticator,
		pointAccrualService,
		authenticator,
		mws,
	)
	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: handler,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-c
		server.Close()
		pointAccrualService.Stop()
	}()

	log.Fatal(server.ListenAndServe())
}

func parseVariables() Config {
	var cfg = Config{
		RunAddress:           "localhost:8080",
		DatabaseURI:          "postgresql://localhost/gophermart?user=user&password=password",
		AccrualSystemAddress: "",
		SecretKey:            "secret key",
		MigrationDir:         "./migrations",
	}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "Run address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Database URI")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "Accrual system address")
	flag.Parse()

	return cfg
}
