#!/bin/sh
set -e

echo "Glacier entrypoint"

# Ensure we are in the /app folder
cd /app

PUID=${PUID:-1000}
PGID=${PGID:-1000}

# If Root (PUID 0)
if [ "$PUID" -eq 0 ]; then
    echo "Running as root..."
    # no perm changed needed
    # Run the app directly
    exec "$@"
fi

# Create the Group
if ! getent group appgroup >/dev/null; then
    addgroup -g "$PGID" appgroup
fi

# Create the User
if ! getent passwd appuser >/dev/null; then
    adduser -u "$PUID" -G appgroup -h /app -D appuser
fi

# todo load any appdirs
APP_DIRS="GLACIER_CONFIG_DIR GLACIER_GAME_DIR"

# Initialize Directories
# loop through the env variables
for DIR_VAR in $APP_DIRS; do
    # value of the variable name
    DIR_PATH=$(eval echo "\$$DIR_VAR")

    if [ -n "$DIR_PATH" ]; then
        echo "Setting dir permissions: $DIR_PATH"
        mkdir -p "$DIR_PATH"
        chown -R "$PUID:$PGID" "$DIR_PATH"
    else
        echo "Warning: $DIR_VAR is not set, skipping..."
    fi
done

# Drop permissions and run the application
echo "Launching Glacier as $PUID:$PGID"
exec su-exec appuser "$@"