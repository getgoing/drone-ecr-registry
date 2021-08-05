A registry credential extension to Drone registry extension to pull images from private AWS ECR repo. _Please note this project requires Drone server version 1.4 or higher._

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env=ECR_REGISTRY_IDS="awsAccountId1,awsAccountId2,..." \
  --restart=always \
  --name=credentials reptiloid666/drone-ecr-registry
```

Update your runner configuration to include the plugin address and the shared secret.

```text
DRONE_REGISTRY_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_REGISTRY_PLUGIN_TOKEN=bea26a2221fd8090ea38720fc445eca6
```
