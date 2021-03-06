// Code generated by go-swagger; DO NOT EDIT.

package read_catalog_entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetAllCatalogEntryPackagesHandlerFunc turns a function with the right signature into a get all catalog entry packages handler
type GetAllCatalogEntryPackagesHandlerFunc func(GetAllCatalogEntryPackagesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAllCatalogEntryPackagesHandlerFunc) Handle(params GetAllCatalogEntryPackagesParams) middleware.Responder {
	return fn(params)
}

// GetAllCatalogEntryPackagesHandler interface for that can handle valid get all catalog entry packages params
type GetAllCatalogEntryPackagesHandler interface {
	Handle(GetAllCatalogEntryPackagesParams) middleware.Responder
}

// NewGetAllCatalogEntryPackages creates a new http.Handler for the get all catalog entry packages operation
func NewGetAllCatalogEntryPackages(ctx *middleware.Context, handler GetAllCatalogEntryPackagesHandler) *GetAllCatalogEntryPackages {
	return &GetAllCatalogEntryPackages{Context: ctx, Handler: handler}
}

/*GetAllCatalogEntryPackages swagger:route GET /catalogEntries Read catalog entries getAllCatalogEntryPackages

Get all entries in the catalog

*/
type GetAllCatalogEntryPackages struct {
	Context *middleware.Context
	Handler GetAllCatalogEntryPackagesHandler
}

func (o *GetAllCatalogEntryPackages) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetAllCatalogEntryPackagesParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
