# About

**Glacier** is the core server component of the project.

It is responsible for:

* Managing system configuration
* Connecting to metadata providers
* Coordinating download clients
* Indexing and organizing content
* Handling user accounts and permissions

Glacier acts as the central hub. All clients and supporting services communicate with Glacier to
retrieve metadata, manage downloads, and maintain your library.

It is designed to be self-hosted, giving you full control over your data, configuration, and infrastructure.

## Supported Platforms

| Platform | Architecture   | Installation         |
|----------|----------------|----------------------|
| Docker   | x86_64 / Arm64 | [Install](docker.md) |

Additional platforms and variants may be supported based on feedback in the future.

# After Installation

Once Glacier is installed and running:

1. Open the Web UI in your browser.
2. Complete the initial configuration.
3. Configure metadata providers, indexers, and download clients.

Continue to the
[Setup Guide](../../config/glacier)
to complete the full configuration process.
