#!/bin/bash

set -e

REPO_URL="https://github.com/razzat008/letsgodb.git"
PROJECT_DIR="letsgodb"

echo "---- letsgodb Setup Script (Arch-based Systems) ----"

# 1. Install Go if not present
if ! command -v go &> /dev/null; then
    echo "Go not found. Installing Go with pacman..."
    sudo pacman -S --noconfirm go
else
    echo "Go is already installed."
fi

# 2. Clone the repository if not present
if [ ! -d "$PROJECT_DIR" ]; then
    echo "Cloning letsgodb repository..."
    git clone "$REPO_URL"
else
    echo "Repository already cloned."
fi

cd "$PROJECT_DIR"

# 3. Initialize Go modules
echo "Initializing Go modules..."
go mod tidy

# 4. Create data directory if not present
if [ ! -d "data" ]; then
    mkdir data
    echo "Created data directory."
fi

echo "---- Setup Complete ----"
echo "To run the project:"
echo "  cd $PROJECT_DIR"
echo "  go run ."
echo
echo "For help, type: `help;` in the letsgodb prompt."
echo "For full help, type: `helpall;` in the letsgodb prompt."
