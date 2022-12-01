package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aretaja/godevmanapi/config"
	_ "github.com/aretaja/godevmanapi/docs"
	"github.com/aretaja/godevmanapi/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	Conf    *config.Configuration
	Router  *chi.Mux
	Handler *handlers.Handler
	Version string
}

func (a *App) Initialize() {
	// Config
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	a.Conf = c

	// Router instance
	a.Router = chi.NewRouter()
	a.initializeMiddleware()

	// Handler instance
	a.Handler = new(handlers.Handler)
	err = a.Handler.Initialize(a.Conf.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Handler instance
	a.initializeRoutes()
}

// Midleware activation
func (a *App) initializeMiddleware() {
	r := a.Router

	// httplog configuration
	logConf := httplog.Options{
		JSON:    true,
		Concise: false,
	}

	if os.Getenv("GODEVMANAPI_LOGPLAIN") != "" {
		logConf.JSON = false
		logConf.Concise = true
	}

	logger := httplog.NewLogger("godevmanapi", logConf)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

}

// Route definitions
func (a *App) initializeRoutes() {
	r := a.Router

	// Welcome
	r.Get("/", a.Handler.Hello)

	// Version
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		handlers.VersionSwagger() // Prevent function not used warning
		handlers.RespondJSON(w, r, http.StatusOK, handlers.StatusResponse{
			Code:    strconv.Itoa(http.StatusOK),
			Message: a.Version,
		})
	})

	// Routes for "/connections/providers" resource
	r.Route("/connections/providers", func(r chi.Router) {
		// Takes parameters: count(100), start(0). Uses default if not set.
		r.Get("/", a.Handler.GetConProviders)
		r.Get("/count", a.Handler.CountConProviders)
		r.Post("/", a.Handler.CreateConProvider)

		// Subroutes
		r.Route("/{con_prov_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetConProvider)
			r.Put("/", a.Handler.UpdateConProvider)
			r.Delete("/", a.Handler.DeleteConProvider)
			r.Get("/connections", a.Handler.GetConProviderConnections)
		})
	})

	// Routes for "/connections/capacities" resource
	r.Route("/connections/capacities", func(r chi.Router) {
		// Takes parameters: count(100), start(0). Uses default if not set.
		r.Get("/", a.Handler.GetConCapacities)
		r.Get("/count", a.Handler.CountConCapacities)
		r.Post("/", a.Handler.CreateConCapacity)

		// Subroutes
		r.Route("/{con_cap_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetConCapacity)
			r.Put("/", a.Handler.UpdateConCapacity)
			r.Delete("/", a.Handler.DeleteConCapacity)
			r.Get("/connections", a.Handler.GetConCapacityConnections)
		})
	})

	// Swagger
	r.Mount("/swagger", httpSwagger.WrapHandler)

	// Custom 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondError(w, r, http.StatusNotFound, "Route does not exist")
	})

	// Custom 405 handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondError(w, r, http.StatusMethodNotAllowed, "Method is not valid")
	})

}

func (a *App) Run() {
	fmt.Printf("Starting up on:%s\n", a.Conf.ApiListen)
	err := http.ListenAndServe(a.Conf.ApiListen, a.Router)
	if err != nil {
		log.Printf("Failed to launch api server:%+v\n", err)
	}
}
