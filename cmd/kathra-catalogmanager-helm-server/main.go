// Code generated by go-swagger; DO NOT EDIT.

package main

import (
	"log"
	"os"

	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/kathra-project/kathra-catalogmanager-helm/controllers"
	"github.com/kathra-project/kathra-catalogmanager-helm/restapi"
	"github.com/kathra-project/kathra-catalogmanager-helm/restapi/operations"
	svc "github.com/kathra-project/kathra-catalogmanager-helm/services"
	"github.com/robfig/cron"
)

// This file was generated by the swagger tool.
// Make sure not to overwrite this file after you generated it because all your edits would be lost!

func main() {
	schedule()

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewKathraCatalogmanagerHelmAPI(swaggerSpec)
	registerHandlers(api)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Kathra Catalog Manager"
	parser.LongDescription = "KATHRA Catalog Management API permetting : \n * Generate source's packages from templates \n * Insert catalog entry from template  \n \n * Insert catalog entry from file  \n\n * Insert catalog entry from source repository  \n * Read catalog entries from catalog \n "

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}

func registerHandlers(api *operations.KathraCatalogmanagerHelmAPI) {
	api.ReadCatalogEntriesGetAllCatalogEntryPackagesHandler = controllers.GetAllCatalogEntries()
	api.ReadCatalogEntriesGetCatalogEntryPackageHandler = controllers.GetAllCatalogEntry()
	api.ReadCatalogEntriesGetCatalogEntryPackageVersionsHandler = controllers.GetAllCatalogEntryVersions()
	api.ReadCatalogEntriesGetCatalogEntryFromVersionHandler = controllers.GetAllCatalogEntryVersion()
}

func schedule() {

	var helmSvc = svc.GetHelmServiceInstance()
	var cronSettings = os.Getenv("HELM_UPDATE_INTERVAL")
	if cronSettings == "" {
		cronSettings = "1 * * * * *"
	}
	cronScheduler := cron.New()
	cronScheduler.AddFunc(cronSettings, func() {
		log.Printf("Update Helm Chart.. begin")
		helmSvc.UpdateFromResourceManager()
		helmSvc.HelmUpdate()
		helmSvc.HelmLoadAllInMemory()
		log.Printf("Update Helm Chart.. done")
	})
	cronScheduler.Start()
}
