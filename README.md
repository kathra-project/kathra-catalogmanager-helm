# Kathra Catalog Manager Helm

This service implements Kathra Catalog Manager 1.1.0.
Using helm-client, it can do :

* Get all charts available from : stable, appscode and local kathra chart museum
* Extract parameters from README.md files (or questions.yaml for Rancher Charts)
* Generate chart from template (eg : Rest Api)
* Push a new chart into kathra chart museum


## Configuration

| Env var                         | Description                          | Default                                   |
| --------------------------------- | ------------------------------------ | ----------------------------------------- |
| `REPOSITORIES_CONFIG`            | File repositories settings          | `repositories.yaml`                 |
| `RESOURCE_MANAGER_URL`            | Resource Manager URL          | |
| `KEYCLOAK_AUTH_URL`            | Keycloak Auth Url          | |
| `KEYCLOAK_REALM`            | Keycloak Realm          | |
| `USERNAME`            | Keycloak's username for technical user          | |
| `PASSWORD`            | Keycloak's password for technical user          | |
| `KEYCLOAK_CLIENT_ID`            | Keycloak client ID          | |
| `KEYCLOAK_CLIENT_SECRET`            | Keycloak client Secret          | |
| `HELM_UPDATE_INTERVAL`            | Cron settings for helm update            | `* * * * *`                    |



## How to run

```
go run cmd/kathra-catalogmanager-helm-server/main.go
```