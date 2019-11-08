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
| `KATHRA_REPO_NAME`            | Chart repository name for Helm           | `kathra-local`                 |
| `KATHRA_REPO_URL`             | Chart repository URL                     |                                |
| `KATHRA_REPO_CREDENTIAL_ID`   | Chart repository Credential ID           |                                |
| `KATHRA_REPO_SECRET`          | Chart repository Credential secret       |                                |
| `HELM_UPDATE_INTERVAL`            | Cron settings for helm update            | `* * * * *`                    |



## How to run

```
go run main.go
```