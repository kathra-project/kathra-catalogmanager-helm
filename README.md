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
| `HELM_UPDATE_INTERVAL`            | Cron settings for helm update            | `* * * * *`                    |



## How to run

```
go run main.go
```