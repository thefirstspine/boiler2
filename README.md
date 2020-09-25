```
██████╗      ██████╗     ██╗    ██╗         ███████╗    ██████╗
██╔══██╗    ██╔═══██╗    ██║    ██║         ██╔════╝    ██╔══██╗
██████╔╝    ██║   ██║    ██║    ██║         █████╗      ██████╔╝
██╔══██╗    ██║   ██║    ██║    ██║         ██╔══╝      ██╔══██╗
██████╔╝    ╚██████╔╝    ██║    ███████╗    ███████╗    ██║  ██║
╚═════╝      ╚═════╝     ╚═╝    ╚══════╝    ╚══════╝    ╚═╝  ╚═╝
```

## About

Boiler is a Docker orchestration software built on top of nginx & docker.

## Prerequisites

You need a server with git, docker, nginx and certbot installed. To install these dependencies, you can run the ./setup.sh script.

Also, your server needs to access your repositories in order to boil some apps.

## Usage

Deploy a project

```
./boiler deploy [project] [options]
```

Here's the options available:
- `tag_or_branch`: The tag or the branch to deploy. Default to master.
- `skip_build`: Skip the docker build. 1 or 0. Defaults to 0.
- `skip_sign`: Skip the certbot call. 1 or 0. Defaults to 0.
- `skip_clean`: Skip the clean at the end of the deployment. 1 or 0. Defaults to 0.

Example with Arena

`./boiler deploy arena --tag_or_branch=1.0.0`

## Configuring a project

The configurations are stored inside in JSON located at `boiler.json`. Here’s an explaination of the content of the file:

```
{
  "projects": {
    "{project name}": {
      "forward": "{domain}::{container's port}",
      "repository": "{repository - HTTPS or SSH}",
      "env": [
        "{environment key}={value}"
      ]
    }
  }
}
```