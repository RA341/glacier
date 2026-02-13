---
title: Docker
---

Glacier uses ghcr.io as its docker distribution

Official hub image
```
ghcr.io/ra341/glacier
```

## Docker run

You can use docker run to quickly test the project 

:::important
For a more permanent setup, compose is recommended
:::

```
docker run ghcr.io/ra341/glacier:canary  
```

## Compose

```yaml
services:  
  glacier:
    container_name: glacier
    image: ghcr.io/ra341/glacier:canary
    environment:
      GLACIER_LOG_VERBOSE: true
      GLACIER_GAME_DIR:  /media/gamestop
      GLACIER_INCOMPLETE_DIR: /media/downloads/games
      GLACIER_CONFIG_YML_PATH: /config/glacier.yml
    volumes:
      - ./config:/config
      -  ${media}:/media
    ports:
      - "6699:6699"
    restart: unless-stopped
```
