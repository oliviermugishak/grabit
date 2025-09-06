# Grabit - Your All-in-One YouTube Downloader 🎥🎵

[![Version](https://img.shields.io/badge/version-1.1.0-blue)](https://github.com/OlivierMugishaK/grabit)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Grabit is a **CLI YouTube downloader** written in **Go**, designed to download **single videos, playlists, or multiple videos** in **audio or video format**, with a **progress bar**, **sanitized filenames**, **quality selection**, and **parallel downloads** for faster processing.

---

## Features

* 🎬 Download **single YouTube videos**
* 📃 Download **playlists or multiple videos**
* 🎵 **Audio-only mode** (m4a)
* 📺 **Quality selection**: best, worst, 720p, 1080p, etc.
* ⏱ **Parallel downloads** with configurable concurrency (`-concurrent`)
* 📊 **CLI progress bar** with percentage and ETA per video
* 🛡 **Sanitized filenames** for safe downloads
* 🏷 **Developer info** and version display with `--version`
* 💻 Easy to install and run globally as `grabit`

---

## Installation

### Requirements

* [Go](https://golang.org/dl/) (for building Grabit)
* curl (for downloading yt-dlp if missing)

### Install

1. Clone the repo:

```bash
git clone https://github.com/OlivierMugishaK/grabit.git
cd grabit
```

2. Ensure Go is installed:

```bash
go version
```

3. Run the installer:

```bash
chmod +x install.sh
./install.sh
```

4. Test the installation:

```bash
grabit --help
grabit --version
```

---

## Usage Examples

### Download a single video

```bash
grabit -urls="https://www.youtube.com/watch?v=V1Pl8CzNzCw"
```

### Download audio only

```bash
grabit -urls="https://www.youtube.com/watch?v=V1Pl8CzNzCw" -audio
```

### Download multiple videos

```bash
grabit -urls="https://youtu.be/ID1,https://youtu.be/ID2"
```

### Download a playlist

```bash
grabit -urls="https://www.youtube.com/playlist?list=PLAYLIST_ID"
```

### Specify video quality (e.g., 720p)

```bash
grabit -urls="https://www.youtube.com/watch?v=V1Pl8CzNzCw" -quality="720p"
```

### Download multiple videos in parallel

```bash
grabit -urls="https://youtu.be/ID1,https://youtu.be/ID2,https://youtu.be/ID3" -concurrent=3
```

> Use `-concurrent=N` to control the number of simultaneous downloads. Default is 3.

---

## Flags

| Flag          | Description                                       |
| ------------- | ------------------------------------------------- |
| `-urls`       | Comma-separated video or playlist URLs (required) |
| `-audio`      | Download audio only (m4a)                         |
| `-quality`    | Video quality (`best`, `worst`, `720p`, etc.)     |
| `-out`        | Output directory (default: `downloads`)           |
| `-version`    | Show Grabit version and developer info            |
| `-concurrent` | Number of parallel downloads (default: 3)         |

---

## Developer

Olivier M. Kwizera – [GitHub](https://github.com/OlivierMugishaK)

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

