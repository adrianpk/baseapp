package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/app"
	"gitlab.com/kabestan/repo/baseapp/internal/mig"
	repo "gitlab.com/kabestan/repo/baseapp/internal/repo/pg"
	"gitlab.com/kabestan/repo/baseapp/internal/seed"
	"gitlab.com/kabestan/repo/baseapp/internal/svc"
)

type contextKey string

const (
	// Replace by prefered
	appName = "kbs"
)

var (
	a *app.App
)

func main() {
	// Replace by custom envar prefix
	cfg := kbs.LoadConfig("kbs")
	log := kbs.NewLogger(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	checkStopEvents(ctx, cancel)

	// App
	a, err := app.NewApp(cfg, log, appName)
	if err != nil {
		exit(log, err)
	}

	// Database connection
	db, err := kbs.NewPostgresConn(cfg, log)
	if err != nil {
		exit(log, err)
	}

	// Migrator
	mg, err := mig.NewMigrator(cfg, log, "migrator", db)
	if err != nil {
		log.Error(err)
	}

	// Seeder
	sd, err := seed.NewSeeder(cfg, log, "seeder", db)
	if err != nil {
		log.Error(err)
	}

	// Mailer
	ml, err := kbs.NewSESMailer(cfg, log, "ses-mailer")
	if err != nil {
		exit(log, err)
	}

	// Repos
	userRepo := repo.NewUserRepo(cfg, log, "user-repo", db)
	accountRepo := repo.NewAccountRepo(cfg, log, "account-repo", db)
	authRepo := repo.NewAuthRepo(cfg, log, "auth-repo", db)

	// Volatile Repo
	//authRepo := vrepo.NewAuthRepo(cfg, log, "auth-repo")

	// Core service
	svc := svc.NewService(cfg, log, "core-service", db)

	// Service dependencies
	svc.Mailer = ml
	svc.UserRepo = userRepo
	svc.AccountRepo = accountRepo
	svc.AuthRepo = authRepo

	// App dependencies
	a.Migrator = mg
	a.Seeder = sd
	a.WebEP.SetService(svc)

	// Init service
	a.Init()

	// Start service
	a.Start()

	log.Error(err, fmt.Sprintf("%s service stoped", appName))
}

func exit(log kbs.Logger, err error) {
	log.Error(err)
	os.Exit(1)
}

func checkStopEvents(ctx context.Context, cancel context.CancelFunc) {
	go checkSigterm(cancel)
	go checkCancel(ctx)
}

func checkSigterm(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cancel()
}

func checkCancel(ctx context.Context) {
	<-ctx.Done()
	a.Stop()
	os.Exit(1)
}
