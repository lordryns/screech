#!/data/data/com.termux/files/usr/bin/bash

INSTALL_DIR="$HOME/.local/bin"
CLI_NAME="screech"
PROFILE_HOOK="$PREFIX/etc/profile.d/${CLI_NAME}.sh"

# Ensure install dir exists
mkdir -p "$INSTALL_DIR"

# Install or update CLI
if [ -f "$INSTALL_DIR/$CLI_NAME" ]; then
    echo -e "\033[1;34mUpdating existing $CLI_NAME...\033[0m"
else
    echo -e "\033[1;34mInstalling $CLI_NAME...\033[0m"
fi
cp "$CLI_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$CLI_NAME"

# Ensure profile.d exists
mkdir -p "$PREFIX/etc/profile.d"

# Add PATH hook only if not present
if [ ! -f "$PROFILE_HOOK" ]; then
    cat > "$PROFILE_HOOK" <<EOF
# Added by $CLI_NAME installer
export PATH=\$PATH:$INSTALL_DIR
EOF
    echo -e "\033[1;34mPATH hook created in $PROFILE_HOOK\033[0m"
fi

# Install termux-api package if missing
if ! command -v termux-notification >/dev/null 2>&1; then
    echo -e "\033[1;34mInstalling termux-api package...\033[0m"
    pkg install -y termux-api
fi

# Check Termux:API app availability
if ! command -v termux-notification >/dev/null 2>&1; then
    echo -e "\033[1;31mTermux:API app is required but not installed.\033[0m"
    echo -e "\033[1;33mDownload it from:\033[0m https://f-droid.org/en/packages/com.termux.api/"
    exit 1
fi

# Run storage setup if not already done
if [ ! -d "$HOME/storage" ]; then
    echo -e "\033[1;34mRequesting storage access...\033[0m"
    termux-setup-storage
fi

# Final feedback
echo -e "\033[1;32m$CLI_NAME installed successfully!\033[0m"
echo -e "\033[1;33mRestart Termux and run '$CLI_NAME' from anywhere.\033[0m"
