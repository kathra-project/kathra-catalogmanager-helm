package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

type KeycloakAuthService struct {
	Host         string
	Realm        string
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	Token        *KeycloakAuthServiceToken
}

type KeycloakAuthServiceToken struct {
	Token    string
	ExpireAt int64
}

var keycloakAuthServiceInstance *KeycloakAuthService
var keycloakAuthServiceOnce sync.Once

func GetKeycloakAuthServiceInstance() *KeycloakAuthService {
	keycloakAuthServiceOnce.Do(func() {
		keycloakAuthServiceInstance = NewKeycloakAuthServiceFromEnvVar()
	})
	return keycloakAuthServiceInstance
}

func NewKeycloakAuthServiceFromEnvVar() *KeycloakAuthService {
	var host = os.Getenv("KEYCLOAK_AUTH_URL")
	var realm = os.Getenv("KEYCLOAK_REALM")
	var username = os.Getenv("USERNAME")
	var password = os.Getenv("PASSWORD")
	var clientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	var clientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	return &KeycloakAuthService{Host: host, Realm: realm, Username: username, Password: password, ClientID: clientID, ClientSecret: clientSecret}
}

func (svc *KeycloakAuthService) GetToken() (string, error) {

	if svc.Token != nil && time.Now().Unix() < svc.Token.ExpireAt {
		return svc.Token.Token, nil
	}

	formData := url.Values{
		"username":      {svc.Username},
		"password":      {svc.Password},
		"grant_type":    {"password"},
		"client_id":     {svc.ClientID},
		"client_secret": {svc.ClientSecret},
	}
	url := svc.Host + "/realms/" + svc.Realm + "/protocol/openid-connect/token"
	fmt.Printf("url: %+v\n", url)
	fmt.Printf("formData: %+v\n", formData)
	resp, err := http.PostForm(url, formData)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["access_token"] == nil {
		fmt.Printf("%+v\n", result)
		return "", errors.New("No token")
	}
	var expireIn, errorParseExpireIn = strconv.Atoi(fmt.Sprintf("%v", result["expires_in"]))
	if errorParseExpireIn != nil {
		fmt.Printf("%+v\n", result)
		return "", errors.New("Unable to parse 'expires_in'")
	}
	svc.Token = &KeycloakAuthServiceToken{Token: fmt.Sprintf("%v", result["access_token"]), ExpireAt: time.Now().Unix() + int64(expireIn)}
	return svc.Token.Token, nil
}
