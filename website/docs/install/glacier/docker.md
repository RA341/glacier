# Docker

Glacier Docker images are distributed only through the GitHub Container Registry (GHCR).

## Official Image

```
ghcr.io/ra341/glacier
```

You can pull the latest canary build with:

```
docker pull ghcr.io/ra341/glacier:canary
```

## Quick Start (docker run)

To quickly test Glacier:

```bash
docker run -p 6699:6699 ghcr.io/ra341/glacier:canary
```

This will start Glacier with default settings and expose the Web UI on port **6699**.

Default Credentials will be

* User: `admin`
* Pass: `admin6699`

:::info
It is recommended to use Docker compose for a more permanent setup
:::

## Docker Compose

Using Docker Compose allows you to configure environment variables, volumes, and restart policies more cleanly.

```yaml
services:
  glacier:
    container_name: glacier
    image: ghcr.io/ra341/glacier:canary
    environment:
      GLACIER_LOG_VERBOSE: true
      GLACIER_GAME_DIR: /media/gamestop
      GLACIER_INCOMPLETE_DIR: /media/downloads/games
      GLACIER_CONFIG_YML_PATH: /config/glacier.yml
    volumes:
      - ./config:/config
      - ${media}:/media
    ports:
      - "6699:6699"
    restart: unless-stopped
```


After the container starts, open your browser and navigate to:

```
http://localhost:6699
```

Then proceed to the [Setup Guide](../../config/glacier) to complete configuration.


### Configuration Breakdown

Below is a non-exhaustive list of commonly used environment variables, especially those relevant during initial setup.

For a complete list of available configuration options:

* Check the server logs at startup (a configuration table is printed automatically), or

* Visit the Config page in the Web UI.

#### Common Environment Variables

* **GLACIER_LOG_VERBOSE** – Enables verbose logging
* **GLACIER_GAME_DIR** – Directory where completed games are stored
* **GLACIER_INCOMPLETE_DIR** – Directory for in-progress downloads
* **GLACIER_CONFIG_YML_PATH** – Path to the Glacier configuration file

### Volumes

* `./config:/config` — Persists Glacier configuration
* `${media}:/media` — Mounts your media directory into the container
