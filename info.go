package main

import "fmt"

const (
	version   = "1.1.0"
	developer = "Olivier M.K"
	github    = "https://github.com/oliviermugishak"
	banner    = `
 ______     ______     ______     ______     __     ______  
/\  ___\   /\  == \   /\  __ \   /\  == \   /\ \   /\__  _\ 
\ \ \__ \  \ \  __<   \ \  __ \  \ \  __<   \ \ \  \/_/\ \/ 
 \ \_____\  \ \_\ \_\  \ \_\ \_\  \ \_____\  \ \_\    \ \_\ 
  \/_____/   \/_/ /_/   \/_/\/_/   \/_____/   \/_/     \/_/ 
   Grabit - Your All-in-One YouTube Downloader`
)

// showHelp prints a stylish help message
func showHelp() {
	fmt.Println(`
Usage: grabit [options]

Options:
  -urls     Comma-separated YouTube URLs or playlist URLs (required)
  -audio    Download audio only (m4a)
  -quality  Video quality: best, worst, 720p, 1080p, etc. (default: best)
  -out      Output directory (default: downloads)
  -version  Show Grabit version and developer info

Examples:
  grabit -urls="https://www.youtube.com/watch?v=ID"
  grabit -urls="https://www.youtube.com/playlist?list=PLAYLIST_ID" -audio
  grabit -urls="https://youtu.be/ID1,https://youtu.be/ID2" -quality="720p"`)
}
