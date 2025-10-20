#!/bin/bash

set -e

echo "Building Axion..."
go build -buildvcs=false -o axion

# Detect shell
if [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="$HOME/.zshrc"
    SHELL_NAME="zsh"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_RC="$HOME/.bashrc"
    SHELL_NAME="bash"
else
    # Fallback - check default shell
    case "$SHELL" in
        */zsh)
            SHELL_RC="$HOME/.zshrc"
            SHELL_NAME="zsh"
            ;;
        */bash)
            SHELL_RC="$HOME/.bashrc"
            SHELL_NAME="bash"
            ;;
        *)
            SHELL_RC="$HOME/.profile"
            SHELL_NAME="unknown"
            ;;
    esac
fi

echo "Detected shell: $SHELL_NAME"

# Create local bin if it doesn't exist
mkdir -p ~/.local/bin

INSTALL_PATH="$HOME/.local/bin/axion"
if [ -f "$INSTALL_PATH" ]; then
    echo "Removing existing axion..."
    rm -f "$INSTALL_PATH"
fi

cp "$(pwd)/axion" "$INSTALL_PATH"
chmod +x "$INSTALL_PATH"
echo "✓ Installed to: $INSTALL_PATH"

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo ""
    echo "Adding ~/.local/bin to PATH in $SHELL_RC..."
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$SHELL_RC"
    echo "✓ Updated $SHELL_RC"
    echo ""
    echo "⚠ Run: source $SHELL_RC"
    echo "   or restart your terminal"
else
    echo "✓ ~/.local/bin already in PATH"
fi

echo ""
echo "✓ Installation complete!"
echo ""
echo "Run 'axion --help' to get started"