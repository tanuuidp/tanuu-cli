# tanuu-cli
CLI tool for Tanuu Demo (and beta bootstrapping)

Download from [here](https://github.com/tanuuidp/tanuu-cli/releases)


note: compiled releases for linux and mac. Tested on macos. Because this is not signed with apple dev certificate, there is the 'this is not from a trusted developer' warning and workaround needed.


## Demo use
### PREREQUISITES 
1. You must have docker running, meaning /var/run/docker.sock must be accessible and running. The easiest way to accomplish this is with docker desktop. If docker desktop is not available, there are other methods to setup docker.sock, but they will not be addressed here.
2. github user account. You also need to create a token. This is done in github, settings, Developer settings, personal access tokens, classic. The scopes needed are repo, and read:user. This token never leaves your laptop.


### Run the demo
```
tanuu-cli demo --ghtoken <your token here>
```
Once running, the links to connect will be displayed in the terminal.



## Bootstrapping a Tanuu management cluster
Coming soon.