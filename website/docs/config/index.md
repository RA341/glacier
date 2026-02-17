Here’s a clearer, more professional version with better structure and flow:

---

# Configuration

Almost all configuration for both **Glacier** and **Frost** can be set in one of the following ways (in order of priority):

1. **Environment variables**
2. **Web UI**

Environment variables take precedence over the Web UI.

If a configuration value is set via an environment variable, it **cannot** be modified in the Web UI. This is because any changes made in the UI would be overwritten the next time the application restarts.

To prevent this mismatch, the Web UI automatically disables editing for any settings controlled by environment variables.

If you want to manage a setting through the Web UI instead, unset the corresponding environment variable and restart the application.

---

# Other Configuration

The following cannot be configured via environment variables:

* Users
* Indexers
* Downloaders
* Metadata providers

These can only be created and managed through the Web UI.
