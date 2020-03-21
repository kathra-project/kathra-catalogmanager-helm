package services

import (
	"fmt"
	"log"
	"net/url"
	"os"

	openApiRt "github.com/go-openapi/runtime"
	openApiRtClient "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"
	apiClient "github.com/kathra-project/kathra-resourcemanager-client-go/client"
	apiClientBR "github.com/kathra-project/kathra-resourcemanager-client-go/client/binary_repositories"
	apiClientGroup "github.com/kathra-project/kathra-resourcemanager-client-go/client/groups"
	apiClientUser "github.com/kathra-project/kathra-resourcemanager-client-go/client/users"
	models "github.com/kathra-project/kathra-resourcemanager-client-go/models"
)

type ResourceManagerService struct {
	AuthInfo             *openApiRtClient.Runtime
	ClientAuthInfoWriter openApiRt.ClientAuthInfoWriter
	Formats              strfmt.Registry
	Config               *apiClient.TransportConfig
	KeycloakService      *KeycloakAuthService
}

func NewResourceManagerService() ResourceManagerService {
	var token, errToken = GetKeycloakAuthServiceInstance().GetToken()
	if errToken != nil {
		panic(errToken)
	}
	var formats = strfmt.NewFormats()
	var config = apiClient.DefaultTransportConfig()

	resourceManagerURL, _ := url.Parse(os.Getenv("RESOURCE_MANAGER_URL"))
	schemes := []string{resourceManagerURL.Scheme}
	config.Host = resourceManagerURL.Host
	config.BasePath = resourceManagerURL.Path
	config.Schemes = schemes
	var authInfo = openApiRtClient.New(config.Host, config.BasePath, schemes)
	var ClientAuthInfoWriter = openApiRtClient.BearerToken(token)
	return ResourceManagerService{AuthInfo: authInfo, Formats: formats, Config: config, ClientAuthInfoWriter: ClientAuthInfoWriter, KeycloakService: GetKeycloakAuthServiceInstance()}
}

func (svc ResourceManagerService) initToken() {
	var token, errToken = GetKeycloakAuthServiceInstance().GetToken()
	if errToken != nil {
		panic(errToken)
	}
	svc.ClientAuthInfoWriter = openApiRtClient.BearerToken(token)
}

func (svc ResourceManagerService) getBinaryRepositories() []*models.BinaryRepository {
	svc.initToken()
	var cli = apiClient.NewHTTPClientWithConfig(svc.Formats, svc.Config)

	var params = apiClientBR.NewGetBinaryRepositoriesParams()
	var response, err = cli.BinaryRepositories.GetBinaryRepositories(params, svc.ClientAuthInfoWriter)
	if err != nil {
		fmt.Printf("%+v\n", response)
		log.Println(err)
		return nil
	}
	return response.Payload
}

func (svc ResourceManagerService) getGroupById(id string) *models.Group {
	svc.initToken()
	var cli = apiClient.NewHTTPClientWithConfig(svc.Formats, svc.Config)

	var params = apiClientGroup.NewGetGroupParams()
	params.ResourceID = id
	var response, err = cli.Groups.GetGroup(params, svc.ClientAuthInfoWriter)
	if err != nil {
		fmt.Printf("%+v\n", response)
		log.Println(err)
		return nil
	}
	return response.Payload
}

func (svc ResourceManagerService) getUserById(id string) *models.User {
	svc.initToken()
	var cli = apiClient.NewHTTPClientWithConfig(svc.Formats, svc.Config)

	var params = apiClientUser.NewGetUserParams()
	params.ResourceID = id
	var response, err = cli.Users.GetUser(params, svc.ClientAuthInfoWriter)
	if err != nil {
		fmt.Printf("%+v\n", response)
		log.Println(err)
		return nil
	}
	return response.Payload
}
