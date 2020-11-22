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

## Deployment

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

## Serving for Github Webhooks

```
./boiler serve
```

This will serve on the port 3000

## Configuring a project

The configurations are stored inside in JSON located at `boiler.json`. Here’s an explaination of the content of the file:

```
{
  "config": {
    "githubKey": ""
  },
  "common": {
    "env": [
      "{environment key}={value}"
    ]
  },
  "projects": [
    {
      "name": "{project name}",
      "domain": "{domain}",
      "repository": "{repository - HTTPS or SSH}",
      "env": [
        "{environment key}={value}"
      ]
    }
  ]
}

```

## Complete example

```
# Write config
echo '{
  "config": {
    "githubKey": ""
  },
  "common": {
    "env": [
      "PORT=8080",
    ]
  },
  "projects": [
    {
      "name": "rest",
      "domain": "rest.sandbox.thefirstspine.fr",
      "repository": "https://github.com/thefirstspine/rest.git",
      "env": [
        "ARENA_URL=https://arena.thefirstspine.fr",
        "WEBSITE_URL=https://www.thefirstspine.fr"
      ]
    }
  ]
}' > boiler.json
# Download from github
wget https://github.com/thefirstspine/boiler2/releases/download/v-0.1.0/boiler2_v-0.1.0_linux_amd64.tar.gz boiler2_v-0.1.0_linux_amd64.tar.gz
# Untar release
tar vxf boiler2_v-0.1.0_linux_amd64.tar.gz
# Make it executable
chmod +x boiler2
# Execute file
./boiler2 deploy rest
```