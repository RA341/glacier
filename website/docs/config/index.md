# Configuration

Almost all configuration for both frost and glacier can be set by either via (in the order of importance)

* Env var
* Web UI

If an env var is set you cannot change it on the webui,

This is due to the fact that if a env var is set any changes you make will be overwritten by the env var when the app
restarts.

To prevent this mismatch, the webui will prevent you from editing the config in the UI.

If you want to edit the config in the webui, you should unset the env var and restart the app.

# Other configs

Users, Indexers, Downloaders, Metadata providers cannot be set via the above the method, and can be modified only via
the WebUI
