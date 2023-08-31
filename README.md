# tanuu-cli
CLI tool for Tanuu Demo (and beta bootstrapping)

Download from [here](https://github.com/tanuuidp/tanuu-cli/releases)


note: compiled releases for linux and mac. Tested on macos. Because this is not signed with apple dev certificate, there is the 'this is not from a trusted developer' warning and workaround needed.


## Demo use
### PREREQUISITES 
1. You must have docker or podman (with socket) running, meaning /var/run/docker.sock must be accessible and running.
   > The easiest way to accomplish this is with docker desktop. If docker desktop is not available, there are other methods to setup docker.sock, but they will not be addressed here.
2. Github user account. You also need to create a token. This is done in github, settings, Developer settings, personal access tokens, classic. The scopes needed are repo, and read:user. This token never leaves your laptop.

### Run the demo

```sh
tanuu-cli demo --ghtoken <your token here>
```

<details>
<summary>Using podman or a different path to <code>docker.sock</code></summary>
For example, when using rootless podman, the socket location must be reconfigured or the following error is likely encountered:

```text
failed to list containers: Got permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock
```

To set the socket location, run tanuu-cli with the following environment variables set:

```sh
DOCKER_SOCK="/run/user/1000/podman/podman.sock" DOCKER_HOST="unix://$DOCKER_SOCK" tanuu-cli demo --ghtoken <your token here>
```

The rootless socket will start automatically once any application connects to it. Check `systemctl --user status podman.socket` for status.

</details>

Once running, the links to connect to the services will be displayed in the terminal, e.g. as follows:

```text
READY

ArgoCD URL: http://127.0.0.1:8081
AdminPW:  VEfTyuWC4N1mmiNj

ArgoWorkflow URL: http://127.0.0.1:8082

Backstage URL: http://127.0.0.1:7007

Demo app (once deployed): http://127.0.0.1:8084
```

### Setting log level

To set the log level to `debug` for example, run tanuu-cli with the following environment variable set:

```shell
LOG_LEVEL=debug tanuu-cli demo --ghtoken <your token here>
```

## Bootstrapping a Tanuu management cluster
Coming soon.
