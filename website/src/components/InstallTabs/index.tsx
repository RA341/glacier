import React, { useState, useEffect } from "react";
import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";

// Detects Windows vs Linux/Mac
function detectOS() {
  if (typeof navigator === "undefined") return "linux";
  const ua = navigator.userAgent.toLowerCase();
  if (ua.includes("win")) return "windows";
  return "linux";
}

export default function InstallTabs() {
  const [frostOS, setFrostOS] = useState("linux");

  useEffect(() => {
    setFrostOS(detectOS());
  }, []);

  return (
    <Tabs groupId="component">
      {/* ── FROST ── */}
      <TabItem value="frost" label="❄️ Frost">
        <Tabs groupId="frost-os" defaultValue={frostOS}>
          <TabItem value="windows" label="🪟 Windows">

## Windows Installation

Download and run the Frost installer for Windows.

```powershell
# Download the latest Frost installer
winget install frost
```

Or grab the `.exe` from the [releases page](https://github.com/your-org/frost/releases) and run it.

After installation, launch **Frost** from the Start Menu and sign in to connect to your Glacier server.

          </TabItem>

          <TabItem value="linux" label="🐧 Linux">

## Linux Installation

```bash
# Debian / Ubuntu
sudo apt install frost

# Arch
yay -S frost

# Or via the install script
curl -fsSL https://get.frost.sh | bash
```

After installation, start the client:

```bash
frost
```

Sign in to connect to your Glacier server.

          </TabItem>
        </Tabs>
      </TabItem>

      {/* ── GLACIER ── */}
      <TabItem value="glacier" label="🏔️ Glacier">
        <Tabs groupId="glacier-env">
          <TabItem value="docker" label="🐳 Docker">

## Docker Installation

Make sure you have [Docker](https://docs.docker.com/get-docker/) installed, then run:

```bash
docker pull your-org/glacier:latest
```

Start the Glacier server:

```bash
docker run -d \
  --name glacier \
  -p 8080:8080 \
  -v glacier-data:/data \
  your-org/glacier:latest
```

Or using **Docker Compose** — create a `docker-compose.yml`:

```yaml
version: "3.9"
services:
  glacier:
    image: your-org/glacier:latest
    container_name: glacier
    ports:
      - "8080:8080"
    volumes:
      - glacier-data:/data
    restart: unless-stopped

volumes:
  glacier-data:
```

Then start it with:

```bash
docker compose up -d
```

Glacier will be available at `http://localhost:8080`. Point your Frost client at this address to connect.

          </TabItem>
        </Tabs>
      </TabItem>
    </Tabs>
  );
}