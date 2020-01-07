// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"gitlab.com/kathra-project/kathra-catalogmanager-helm/restapi/operations"
	"gitlab.com/kathra-project/kathra-catalogmanager-helm/restapi/operations/read_catalog_entries"
)

//go:generate swagger generate server --target ../../kathra-catalogmanager-helm --name KathraCatalogmanagerHelm --spec ../../specifications/Services/catalogmanager/swagger.yaml

func configureFlags(api *operations.KathraCatalogmanagerHelmAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.KathraCatalogmanagerHelmAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.ReadCatalogEntriesGetAllCatalogServicesHandler == nil {
		api.ReadCatalogEntriesGetAllCatalogServicesHandler = read_catalog_entries.GetAllCatalogServicesHandlerFunc(func(params read_catalog_entries.GetAllCatalogServicesParams) middleware.Responder {
			return middleware.NotImplemented("operation read_catalog_entries.GetAllCatalogServices has not yet been implemented")
		})
	}
	if api.ReadCatalogEntriesGetCatalogEntryFromVersionHandler == nil {
		api.ReadCatalogEntriesGetCatalogEntryFromVersionHandler = read_catalog_entries.GetCatalogEntryFromVersionHandlerFunc(func(params read_catalog_entries.GetCatalogEntryFromVersionParams) middleware.Responder {
			return middleware.NotImplemented("operation read_catalog_entries.GetCatalogEntryFromVersion has not yet been implemented")
		})
	}
	if api.ReadCatalogEntriesGetCatalogEntryPackageHandler == nil {
		api.ReadCatalogEntriesGetCatalogEntryPackageHandler = read_catalog_entries.GetCatalogEntryPackageHandlerFunc(func(params read_catalog_entries.GetCatalogEntryPackageParams) middleware.Responder {
			return middleware.NotImplemented("operation read_catalog_entries.GetCatalogEntryPackage has not yet been implemented")
		})
	}
	if api.ReadCatalogEntriesGetCatalogEntryPackageVersionsHandler == nil {
		api.ReadCatalogEntriesGetCatalogEntryPackageVersionsHandler = read_catalog_entries.GetCatalogEntryPackageVersionsHandlerFunc(func(params read_catalog_entries.GetCatalogEntryPackageVersionsParams) middleware.Responder {
			return middleware.NotImplemented("operation read_catalog_entries.GetCatalogEntryPackageVersions has not yet been implemented")
		})
	}

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
