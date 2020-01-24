// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	_ "github.com/go-sql-driver/mysql" //não sei pq precisa
	"github.com/joelfelipe94/ModuloAutenticacao/models"
	"github.com/joelfelipe94/ModuloAutenticacao/restapi/operations"
)

//go:generate swagger generate server --target ../../ModuloAutenticacao --name MóduloDeLogin --spec ../swagger.json

func configureFlags(api *operations.MóduloDeLoginAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MóduloDeLoginAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()
	api.LoginHandler = operations.LoginHandlerFunc(func(params operations.LoginParams) middleware.Responder {
		var nome string
		var senha string
		db, err := sql.Open("mysql", "joel:12345@/SistemaAutenticacao")
		if err != nil {
			fmt.Println(err.Error())
			mensagem := models.Erro{Codigo: 500, Mensagem: err.Error()}
			return operations.NewLoginDefault(500).WithPayload(&mensagem)
		}

		db.QueryRow("SELECT nome, senha FROM usuarios where nome = ?", *params.Credenciais.Usuario).Scan(&nome, &senha)
		if *params.Credenciais.Usuario == nome && *params.Credenciais.Senha == senha {
			return operations.NewLoginOK()
		}
		return operations.NewLoginUnauthorized()
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
