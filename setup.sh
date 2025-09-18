#!/data/data/com.termux/files/usr/bin/bash

INSTALL_DIR="$HOME/.local/bin"
CLI_NAME="screech"

# Ensure install dir exists
mkdir -p "$INSTALL_DIR"

# Copy CLI into install dir
cp "$CLI_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$CLI_NAME"

# Ensure profile.d exists
mkdir -p "$PREFIX/etc/profile.d"

# Write a PATH hook for all future sessions
cat > "$PREFIX/etc/profile.d/${CLI_NAME}.sh" <<EOF
# Added by $CLI_NAME installer
export PATH=\$PATH:$INSTALL_DIR
EOF

# Feedback
echo -e "\033[1;32m$CLI_NAME installed successfully!\033[0m"
echo -e "\033[1;33mRestart Termux and run '$CLI_NAME' from anywhere.\033[0m"
