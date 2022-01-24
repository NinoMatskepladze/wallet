package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NinoMatskepladze/wallet/configs"
	"github.com/NinoMatskepladze/wallet/db"
	"github.com/NinoMatskepladze/wallet/handle"
	"github.com/NinoMatskepladze/wallet/responder"
	"github.com/NinoMatskepladze/wallet/service"
	"github.com/go-chi/chi"
	"github.com/go-kit/log"
	"github.com/kelseyhightower/envconfig"
)

var (
	fs       = flag.NewFlagSet("wallet", flag.ExitOnError)
	httpAddr = fs.String("http-address", ":8080", "HTTP address to listen")
)

func main() {
	fs.Parse(os.Args[1:])

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)

	var cfg configs.DBConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Log(err)
	}
	dsn := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	datastore := db.NewDataStore(dsn)
	svc := service.NewService(datastore)
	handler := handle.NewServiceRoutes(svc, responder.NewResponder(logger))

	r := chi.NewRouter()
	r.Route("/wallets", func(w chi.Router) {
		w.Post("/", handler.CreateWallet)
		w.Put("/{wallet_id}", handler.UpdateBalance)
		w.Get("/{wallet_id}", handler.GetWallet)
	})

	server := http.Server{Addr: *httpAddr, Handler: r}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log(
				"transport", "HTTP",
				"error", err,
			)
		}
	}()

	logger.Log("Starting http server on", httpAddr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
