package api

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

type Api struct {
	router   *chi.Mux
	handlers *handlers
}

func New() *Api {
	a := &Api{}
	a.handlers = loadDependencies()
	a.setupRouter()

	return a
}

func (a *Api) Run(ctx context.Context, addr string) error {
	httpServer := &http.Server{
		Addr:        addr,
		Handler:     a.router,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()
		fmt.Printf("shutting down HTTP server\n")
		return httpServer.Shutdown(ctx)
	})

	g.Go(func() error {
		fmt.Printf("starting HTTP server on %s\n", addr)
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	return g.Wait()
}

func (a *Api) setupRouter() {
	a.router = chi.NewRouter()
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)

	a.router.Get("/ping", a.handlers.PingHandler.PingHandler())
	a.router.Get("/swagger/*", a.handlers.SwaggerHandler.SwaggerHandler())

	a.router.Get("/api/payments/{id}", a.handlers.PaymentsHandler.GetPaymentHandler())
	a.router.Post("/api/payments", a.handlers.PaymentsHandler.ProcessPaymentHandler())

}
