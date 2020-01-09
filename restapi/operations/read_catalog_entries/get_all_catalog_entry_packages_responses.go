// Code generated by go-swagger; DO NOT EDIT.

package read_catalog_entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	CatalogEntryPackage "github.com/kathra-project/kathra-core-model-go/models"
)

// GetAllCatalogEntryPackagesOKCode is the HTTP code returned for type GetAllCatalogEntryPackagesOK
const GetAllCatalogEntryPackagesOKCode int = 200

/*GetAllCatalogEntryPackagesOK CatalogEntryPackage with providerId

swagger:response getAllCatalogEntryPackagesOK
*/
type GetAllCatalogEntryPackagesOK struct {

	/*
	  In: Body
	*/
	Payload []CatalogEntryPackage.CatalogEntryPackage `json:"body,omitempty"`
}

// NewGetAllCatalogEntryPackagesOK creates GetAllCatalogEntryPackagesOK with default headers values
func NewGetAllCatalogEntryPackagesOK() *GetAllCatalogEntryPackagesOK {

	return &GetAllCatalogEntryPackagesOK{}
}

// WithPayload adds the payload to the get all catalog entry packages o k response
func (o *GetAllCatalogEntryPackagesOK) WithPayload(payload []CatalogEntryPackage.CatalogEntryPackage) *GetAllCatalogEntryPackagesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all catalog entry packages o k response
func (o *GetAllCatalogEntryPackagesOK) SetPayload(payload []CatalogEntryPackage.CatalogEntryPackage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllCatalogEntryPackagesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]CatalogEntryPackage.CatalogEntryPackage, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetAllCatalogEntryPackagesInternalServerErrorCode is the HTTP code returned for type GetAllCatalogEntryPackagesInternalServerError
const GetAllCatalogEntryPackagesInternalServerErrorCode int = 500

/*GetAllCatalogEntryPackagesInternalServerError Internal error

swagger:response getAllCatalogEntryPackagesInternalServerError
*/
type GetAllCatalogEntryPackagesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetAllCatalogEntryPackagesInternalServerError creates GetAllCatalogEntryPackagesInternalServerError with default headers values
func NewGetAllCatalogEntryPackagesInternalServerError() *GetAllCatalogEntryPackagesInternalServerError {

	return &GetAllCatalogEntryPackagesInternalServerError{}
}

// WithPayload adds the payload to the get all catalog entry packages internal server error response
func (o *GetAllCatalogEntryPackagesInternalServerError) WithPayload(payload string) *GetAllCatalogEntryPackagesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all catalog entry packages internal server error response
func (o *GetAllCatalogEntryPackagesInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllCatalogEntryPackagesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}