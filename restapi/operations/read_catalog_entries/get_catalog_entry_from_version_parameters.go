// Code generated by go-swagger; DO NOT EDIT.

package read_catalog_entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetCatalogEntryFromVersionParams creates a new GetCatalogEntryFromVersionParams object
// no default values defined in spec.
func NewGetCatalogEntryFromVersionParams() GetCatalogEntryFromVersionParams {

	return GetCatalogEntryFromVersionParams{}
}

// GetCatalogEntryFromVersionParams contains all the bound params for the get catalog entry from version operation
// typically these are obtained from a http.Request
//
// swagger:parameters getCatalogEntryFromVersion
type GetCatalogEntryFromVersionParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*CatalogEntryPackage providerId
	  Required: true
	  In: path
	*/
	ProviderID string
	/*CatalogEntryPackage version
	  Required: true
	  In: path
	*/
	Version string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetCatalogEntryFromVersionParams() beforehand.
func (o *GetCatalogEntryFromVersionParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rProviderID, rhkProviderID, _ := route.Params.GetOK("providerId")
	if err := o.bindProviderID(rProviderID, rhkProviderID, route.Formats); err != nil {
		res = append(res, err)
	}

	rVersion, rhkVersion, _ := route.Params.GetOK("version")
	if err := o.bindVersion(rVersion, rhkVersion, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindProviderID binds and validates parameter ProviderID from path.
func (o *GetCatalogEntryFromVersionParams) bindProviderID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ProviderID = raw

	return nil
}

// bindVersion binds and validates parameter Version from path.
func (o *GetCatalogEntryFromVersionParams) bindVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Version = raw

	return nil
}