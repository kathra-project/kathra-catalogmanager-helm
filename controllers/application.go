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
			var entriesAsValues []models.CatalogEntryPackage
			for _, item := range entries {
				entriesAsValues = append(entriesAsValues, *item)
			}
			return api.NewGetAllCatalogEntryPackagesOK().WithPayload(entriesAsValues)
		}
	})
}

func GetAllCatalogEntry() api.GetCatalogEntryPackageHandler {
	return api.GetCatalogEntryPackageHandlerFunc(func(params api.GetCatalogEntryPackageParams) middleware.Responder {
		catalogEntryPackage, err := svc.GetCatalogEntryPackage(params.ProviderID)
		if err != nil {
			fmt.Println(err)
			return api.NewGetCatalogEntryPackageInternalServerError().WithPayload("Get CatalogEntry generates internal error")
		}
		if catalogEntryPackage == nil {
			return api.NewGetCatalogEntryPackageNotFound().WithPayload("Not found")
		}

		versions, err := svc.GetAllCatalogEntryPackageVersionVersions(params.ProviderID)
		if err != nil {
			fmt.Println(err)
			return api.NewGetCatalogEntryPackageInternalServerError().WithPayload("Get CatalogEntry generates internal error")
		}

		var versionsAsPointer []*models.CatalogEntryPackageVersion
		for i, _ := range versions {
			versionsAsPointer = append(versionsAsPointer, versions[i])
		}
		catalogEntryPackage.Versions = versionsAsPointer
		return api.NewGetCatalogEntryPackageOK().WithPayload(*catalogEntryPackage)
	})
}

func GetAllCatalogEntryVersions() api.GetCatalogEntryPackageVersionsHandler {
	return api.GetCatalogEntryPackageVersionsHandlerFunc(func(params api.GetCatalogEntryPackageVersionsParams) middleware.Responder {
		versions, err := svc.GetAllCatalogEntryPackageVersionVersions(params.ProviderID)
		if err != nil {
			fmt.Println(err)
			return api.NewGetCatalogEntryPackageVersionsInternalServerError().WithPayload("Get CatalogEntry generates internal error")
		} else {
			var entriesAsValues []models.CatalogEntryPackageVersion
			for _, item := range versions {
				entriesAsValues = append(entriesAsValues, *item)
			}
			return api.NewGetCatalogEntryPackageVersionsOK().WithPayload(entriesAsValues)
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
