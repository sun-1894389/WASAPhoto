/*
Questo modulo:
Fornisce un modo per creare e configurare un router HTTP per un'API web.
Utilizza httprouter per gestire il routing delle richieste HTTP.
Utilizza logrus per il logging.
Si aspetta di interagire con un database attraverso l'interfaccia database.AppDatabase.
Permette di creare un router HTTP configurato passando un'istanza di Config alla funzione New.
Package api exposes the main API engine. All HTTP APIs are handled here - so-called "business logic" should be here, or
in a dedicated package (if that logic is complex enough).

To use this package, you should create a new instance with New() passing a valid Config. The resulting Router will have
the Router.Handler() function that returns a handler that can be used in a http.Server (or in other middlewares).

Example:

	// Create the API router
	apirouter, err := api.New(api.Config{
		Logger:   logger,
		Database: appdb,
	})
	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("error creating the API server instance: %w", err)
	}
	router := apirouter.Handler()

	// ... other stuff here, like middleware chaining, etc.

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           router,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	apiserver.ListenAndServe()

See the `main.go` file inside the `cmd/webapi` for a full usage example.
*/
package api

import (
	"Wasa-Photo-1894389/service/database"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Photo media folder,indica un percorso dove vengono memorizzate delle foto.
var photoFolder = filepath.Join("/tmp", "media")

// è una struct utilizzata per passare le configurazioni e le dipendenze alla funzione New.
type Config struct {
	// Logger where log entries are sent
	Logger logrus.FieldLogger

	// Database is the instance of database.AppDatabase where data are saved
	Database database.AppDatabase
}

// Router is the package API interface representing an API handler builder
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package
	Handler() http.Handler

	// Close terminates any resource used in the package
	Close() error
}

// New returns a new Router instance
// New è una funzione che accetta un parametro cfg di tipo Config e restituisce un'istanza di Router e un error.
// Crea un nuovo router HTTP e lo configura in base ai parametri forniti.
// Controlla se le configurazioni fornite sono valide e, in caso contrario, restituisce un errore.
func New(cfg Config) (Router, error) {
	// Check if the configuration is correct
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}

	// Create a new router where we will register HTTP endpoints. The server will pass requests to this router to be
	// handled.
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router:     router,
		baseLogger: cfg.Logger,
		db:         cfg.Database,
	}, nil
}
// _router è una struct che implementa l'interfaccia Router e contiene un router HTTP, un logger, e un'istanza del database.
type _router struct {
	router *httprouter.Router

	// baseLogger is a logger for non-requests contexts, like goroutines or background tasks not started by a request.
	// Use context logger if available (e.g., in requests) instead of this logger.
	baseLogger logrus.FieldLogger

	db database.AppDatabase
}
