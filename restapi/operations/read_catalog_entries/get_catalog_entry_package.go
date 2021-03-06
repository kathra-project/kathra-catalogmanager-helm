// Code generated by go-swagger; DO NOT EDIT.

package read_catalog_entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetCatalogEntryPackageHandlerFunc turns a function with the right signature into a get catalog entry package handler
type GetCatalogEntryPackageHandlerFunc func(GetCatalogEntryPackageParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetCatalogEntryPackageHandlerFunc) Handle(params GetCatalogEntryPackageParams) middleware.Responder {
	return fn(params)
}

// GetCatalogEntryPackageHandler interface for that can handle valid get catalog entry package params
type GetCatalogEntryPackageHandler interface {
	Handle(GetCatalogEntryPackageParams) middleware.Responder
}

// NewGetCatalogEntryPackage creates a new http.Handler for the get catalog entry package operation
func NewGetCatalogEntryPackage(ctx *middleware.Context, handler GetCatalogEntryPackageHandler) *GetCatalogEntryPackage {
	return &GetCatalogEntryPackage{Context: ctx, Handler: handler}
}

/*GetCatalogEntryPackage swagger:route GET /catalogEntries/{providerId} Read catalog entries getCatalogEntryPackage

Get an entry in the catalog

*/
type GetCatalogEntryPackage struct {
	Context *middleware.Context
	Handler GetCatalogEntryPackageHandler
}

func (o *GetCatalogEntryPackage) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetCatalogEntryPackageParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
