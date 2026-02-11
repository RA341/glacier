#!/bin/bash

REPO="RA341/glacier"
ASSET_NAME="frost-linux.zip"
DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/$ASSET_NAME"

TARGET_DIR=$(pwd)

echo "Current download directory is: $TARGET_DIR"
read -p "Do you want to specify a different download path? (y/n): " choice

if [[ "$choice" =~ ^[Yy]$ ]]; then
    # -e allows the user to use Tab completion to find their folder
    read -e -p "Enter the full path to the directory: " input_dir
    # Expand tilde (~) to home directory if used
    TARGET_DIR="${input_dir/#\~/$HOME}"

    # Create directory if it doesn't exist
    mkdir -p "$TARGET_DIR"
fi

ZIP_PATH="$TARGET_DIR/$ASSET_NAME"
EXTRACT_PATH="$TARGET_DIR/frost"

echo "Downloading $ASSET_NAME to $TARGET_DIR..."
if command -v curl >/dev/null 2>&1; then
    curl -L "$DOWNLOAD_URL" -o "$ZIP_PATH"
elif command -v wget >/dev/null 2>&1; then
    wget -O "$ZIP_PATH" "$DOWNLOAD_URL"
else
    echo "Error: Neither curl nor wget found. Please install one."
    exit 1
fi

echo "Unzipping to folder: $EXTRACT_PATH..."
mkdir -p "$EXTRACT_PATH"

if command -v unzip >/dev/null 2>&1; then
    unzip -o "$ZIP_PATH" -d "$EXTRACT_PATH"
else
    echo "Error: 'unzip' utility not found. Please install it (e.g., sudo apt install unzip). or unzip manually $ZIP_PATH"
    exit 1
fi

read -p "Would you like to delete the downloaded zip file? (y/n): " delete_choice
if [[ "$delete_choice" =~ ^[Yy]$ ]]; then
    rm "$ZIP_PATH"
    echo "Zip file deleted."
else
    echo "Zip file kept at: $ZIP_PATH"
fi

echo "Done!"