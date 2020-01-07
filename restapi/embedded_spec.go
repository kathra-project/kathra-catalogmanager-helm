// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "KATHRA Catalog Management API permetting : \n * Generate source's packages from templates \n * Insert catalog entry from template  \n \n * Insert catalog entry from file  \n\n * Insert catalog entry from source repository  \n * Read catalog entries from catalog \n ",
    "title": "Kathra Catalog Manager",
    "version": "1.1.0-SNAPSHOT",
    "x-artifactName": "catalogManager",
    "x-groupId": "org.kathra"
  },
  "basePath": "/api/v1",
  "paths": {
    "/catalogEntries": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "Read catalog entries"
        ],
        "summary": "Get all entries in the catalog",
        "operationId": "getAllCatalogServices",
        "responses": {
          "200": {
            "description": "CatalogEntryPackage with providerId",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/CatalogEntryPackage"
              }
            }
          }
        }
      }
    },
    "/catalogEntries/{providerId}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "Read catalog entries"
        ],
        "summary": "Get an entry in the catalog",
        "operationId": "getCatalogEntryPackage",
        "parameters": [
          {
            "type": "string",
            "x-exportParamName": "ProviderId",
            "description": "CatalogEntryPackage providerId",
            "name": "providerId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "CatalogEntryPackage with details",
            "schema": {
              "$ref": "#/definitions/CatalogEntryPackage"
            }
          }
        }
      }
    },
    "/catalogEntries/{providerId}/versions": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "Read catalog entries"
        ],
        "summary": "Get all version for an entry in the catalog",
        "operationId": "getCatalogEntryPackageVersions",
        "parameters": [
          {
            "type": "string",
            "x-exportParamName": "ProviderId",
            "description": "CatalogEntryPackage providerId",
            "name": "providerId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "All versions for CatalogEntryPackage",
            "schema": {
              "$ref": "#/definitions/CatalogEntryPackageVersion"
            }
          }
        }
      }
    },
    "/catalogEntries/{providerId}/versions/{version}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "Read catalog entries"
        ],
        "summary": "Get an entry in the catalog for specific version",
        "operationId": "getCatalogEntryFromVersion",
        "parameters": [
          {
            "type": "string",
            "x-exportParamName": "ProviderId",
            "description": "CatalogEntryPackage providerId",
            "name": "providerId",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "x-exportParamName": "Version",
            "description": "CatalogEntryPackage version",
            "name": "version",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "CatalogEntryVersion with details",
            "schema": {
              "$ref": "#/definitions/CatalogEntryPackageVersion"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "CatalogEntryPackage": {
      "type": "object",
      "x-artifactId": "kathra-core-model",
      "x-go-type": {
        "import": {
          "alias": "CatalogEntryPackageVersion",
          "package": "github.com/kathra-project/kathra-core-model-go"
        },
        "type": "CatalogEntryPackageVersion"
      }
    },
    "CatalogEntryPackageVersion": {
      "type": "object",
      "x-artifactId": "kathra-core-model",
      "x-go-type": {
        "import": {
          "alias": "CatalogEntryPackageVersion",
          "package": "github.com/kathra-project/kathra-core-model-go"
        },
        "type": "CatalogEntryPackageVersion"
      }
    }
  },
  "x-dependencies": [
    {
      "artifactId": "kathra-core-model",
      "artifactVersion": "1.1.0-SNAPSHOT",
      "groupId": "org.kathra",
      "modelPackage": "core.model"
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "KATHRA Catalog Management API permetting : \n * Generate source's packages from templates \n * Insert catalog entry from template  \n \n * Insert catalog entry from file  \n\n * Insert catalog entry from source repository  \n * Read catalog entries from catalog \n ",
    "title": "Kathra Catalog Manager",
    "version": "1.1.0-SNAPSHOT",
    "x-artifactName": "catalogManager",
    "x-groupId": "org.kathra"
  },
  "basePath": "/api/v1",
  "paths": {
    "/catalogEntries": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "Read catalog entries"
        ],
        "summary": "Get all entries in the catalog",
        "operationId": "getAllCatalogServices",
        "responses": {
          "200": {
            "description": "CatalogEntryPackage with providerId",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/CatalogEntryPackage"
              }
            }
          }
        }
      }
    },
    "/catalogEntries/{providerId}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "Read catalog entries"
        ],
        "summary": "Get an entry in the catalog",
        "operationId": "getCatalogEntryPackage",
        "parameters": [
          {
            "type": "string",
            "x-exportParamName": "ProviderId",
            "description": "CatalogEntryPackage providerId",
            "name": "providerId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "CatalogEntryPackage with details",
            "schema": {
              "$ref": "#/definitions/CatalogEntryPackage"
            }
          }
        }
      }
    },
    "/catalogEntries/{providerId}/versions": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "Read catalog entries"
        ],
        "summary": "Get all version for an entry in the catalog",
        "operationId": "getCatalogEntryPackageVersions",
        "parameters": [
          {
            "type": "string",
            "x-exportParamName": "ProviderId",
            "description": "CatalogEntryPackage providerId",
            "name": "providerId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "All versions for CatalogEntryPackage",
            "schema": {
              "$ref": "#/definitions/CatalogEntryPackageVersion"
            }
          }
        }
      }
    },
    "/catalogEntries/{providerId}/versions/{version}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "Read catalog entries"
        ],
        "summary": "Get an entry in the catalog for specific version",
        "operationId": "getCatalogEntryFromVersion",
        "parameters": [
          {
            "type": "string",
            "x-exportParamName": "ProviderId",
            "description": "CatalogEntryPackage providerId",
            "name": "providerId",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "x-exportParamName": "Version",
            "description": "CatalogEntryPackage version",
            "name": "version",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "CatalogEntryVersion with details",
            "schema": {
              "$ref": "#/definitions/CatalogEntryPackageVersion"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "CatalogEntryPackage": {
      "type": "object",
      "x-artifactId": "kathra-core-model",
      "x-go-type": {
        "import": {
          "alias": "CatalogEntryPackageVersion",
          "package": "github.com/kathra-project/kathra-core-model-go"
        },
        "type": "CatalogEntryPackageVersion"
      }
    },
    "CatalogEntryPackageVersion": {
      "type": "object",
      "x-artifactId": "kathra-core-model",
      "x-go-type": {
        "import": {
          "alias": "CatalogEntryPackageVersion",
          "package": "github.com/kathra-project/kathra-core-model-go"
        },
        "type": "CatalogEntryPackageVersion"
      }
    }
  },
  "x-dependencies": [
    {
      "artifactId": "kathra-core-model",
      "artifactVersion": "1.1.0-SNAPSHOT",
      "groupId": "org.kathra",
      "modelPackage": "core.model"
    }
  ]
}`))
}
