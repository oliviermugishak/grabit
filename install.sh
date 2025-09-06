#!/bin/bash

echo "🔧 Building Grabit..."

# Build the whole module
go build -o grabit ./...

if [ $? -ne 0 ]; then
    echo "❌ Build failed. Make sure Go is installed and your module is correct."
    exit 1
fi

echo "📦 Installing Grabit to /usr/local/bin..."
sudo mv grabit /usr/local/bin/
sudo chmod +x /usr/local/bin/grabit

# Check if yt-dlp exists
if ! command -v yt-dlp &> /dev/null
then
    echo "⚠️ yt-dlp not found. Installing yt-dlp..."
    sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
    sudo chmod +x /usr/local/bin/yt-dlp
    echo "✅ yt-dlp installed successfully!"
fi

echo "✅ Grabit installed successfully!"
echo "Try running: grabit --help"
