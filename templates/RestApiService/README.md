# ${CHART_NAME} Helm Chart (v${CHART_VERSION})

${CHART_DESCRIPTION}

## Chart Details

This chart will deploy rest api service.

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm add repo kathra-repo CHART_MUSEUM_URL --username=MY_LOGIN --password=MY_PASSWORD
$ helm install --name my-release ${CHART_NAME}
```

### Configuration

| Parameter                         | Description                          | Default                                   |
| --------------------------------- | ------------------------------------ | ----------------------------------------- |
| `image.repository`                    | Image name                    | `${REGISTRY_HOST}/${IMAGE_NAME}`                         |
| `image.tag`                      | Image tag                     | `${IMAGE_TAG}`                                     |
| `ingress.enabled`          | Enable ingress resource             | `false`                                  |
| `ingress.hosts[0].host`          | Ingress hostname  (if ingress enabled)            | `${CHART_NAME}-example.local`                                  |
| `resources.limits.cpu`             | CPU consumption limit              | `100m`                                         |
| `resources.limits.memory`      | Memory consumption limit  | `128Mi`                                      |
| `resources.requests.cpu`              | CPU requested for starting    | `100m`                                     |
| `resources.requests.memory`            | Memory requested for starting | `128Mi`                                   |
| `replicaCount`            | Number of replica | `1`                                   |


