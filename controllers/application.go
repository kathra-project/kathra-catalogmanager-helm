package controllers

import (
	"fmt"

	middleware "github.com/go-openapi/runtime/middleware"
	api "github.com/kathra-project/kathra-catalogmanager-helm/restapi/operations/read_catalog_entries"
	svc "github.com/kathra-project/kathra-catalogmanager-helm/services"
	"github.com/kathra-project/kathra-core-model-go/models"
)

func GetAllCatalogEntries() api.GetAllCatalogEntryPackagesHandlerFunc {
	return api.GetAllCatalogEntryPackagesHandlerFunc(func(params api.GetAllCatalogEntryPackagesParams) middleware.Responder {
		entries, err := svc.GetAllCatalogEntryPackage()
		if err != nil {
			fmt.Println(err)
			return api.NewGetAllCatalogEntryPackagesInternalServerError().WithPayload("Get CatalogEntries generates internal error")
		} else {
			return api.NewGetAllCatalogEntryPackagesOK().WithPayload(entries)
		}
	})
}

func GetAllCatalogEntry() api.GetCatalogEntryPackageHandler {
	return api.GetCatalogEntryPackageHandlerFunc(func(params api.GetCatalogEntryPackageParams) middleware.Responder {
		versions, err := svc.GetAllCatalogEntryPackageVersionVersions(params.ProviderID)
		if err != nil {
			fmt.Println(err)
			return api.NewGetCatalogEntryPackageInternalServerError().WithPayload("Get CatalogEntry generates internal error")
		} else {
			var versionsAsPointer []*models.CatalogEntryPackageVersion
			for i, _ := range versions {
				versionsAsPointer = append(versionsAsPointer, &versions[i])
			}
			return api.NewGetCatalogEntryPackageOK().WithPayload(models.CatalogEntryPackage{Versions: versionsAsPointer})
		}
	})
}

func GetAllCatalogEntryVersions() api.GetCatalogEntryPackageVersionsHandler {
	return api.GetCatalogEntryPackageVersionsHandlerFunc(func(params api.GetCatalogEntryPackageVersionsParams) middleware.Responder {
		versions, err := svc.GetAllCatalogEntryPackageVersionVersions(params.ProviderID)
		if err != nil {
			fmt.Println(err)
			return api.NewGetCatalogEntryPackageVersionsInternalServerError().WithPayload("Get CatalogEntry generates internal error")
		} else {
			return api.NewGetCatalogEntryPackageVersionsOK().WithPayload(versions)
		}
	})
}

func GetAllCatalogEntryVersion() api.GetCatalogEntryFromVersionHandler {
	return api.GetCatalogEntryFromVersionHandlerFunc(func(params api.GetCatalogEntryFromVersionParams) middleware.Responder {
		version, err := svc.GetCatalogEntryPackageVersionFromProviderId(params.ProviderID, params.Version)
		if err != nil {
			fmt.Println(err)
			return api.NewGetCatalogEntryFromVersionInternalServerError().WithPayload("Get CatalogEntry generates internal error")
		} else {
			return api.NewGetCatalogEntryFromVersionOK().WithPayload(*version)
		}
	})
}
