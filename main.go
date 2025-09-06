package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

const (
	version   = "1.0.0"
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

func main() {
	// CLI flags
	urlList := flag.String("urls", "", "Comma-separated YouTube URLs or playlist URLs")
	audioOnly := flag.Bool("audio", false, "Download audio only")
	outputDir := flag.String("out", "downloads", "Output directory")
	quality := flag.String("quality", "best", "Quality: best, worst, 720p, 1080p, etc.")
	showVersion := flag.Bool("version", false, "Show Grabit version")
	flag.Parse()

	// Banner
	fmt.Println(banner)

	if *showVersion {
		fmt.Printf("Version: %s\nDeveloper: %s\nGitHub: %s\n", version, developer, github)
		return
	}

	if *urlList == "" {
		showHelp()
		return
	}

	// Ensure output directory exists
	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create output directory: %v", err)
	}

	urls := strings.Split(*urlList, ",")

	for _, url := range urls {
		url = strings.TrimSpace(url)
		fmt.Printf("📥 Downloading: %s\n", url)

		args := []string{"-o", fmt.Sprintf("%s/%%(title)s.%%(ext)s", *outputDir), "--newline"}

		if *audioOnly {
			args = append(args, "-f", "bestaudio")
			args = append(args, "--extract-audio", "--audio-format", "m4a")
		} else {
			if *quality == "best" || *quality == "worst" {
				args = append(args, "-f", *quality)
			} else {
				args = append(args, "-f", fmt.Sprintf("bestvideo[height<=%s]+bestaudio/best", *quality))
			}
		}

		args = append(args, url)

		cmd := exec.Command("yt-dlp", args...)
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("⚠️ Failed to get stdout pipe: %v", err)
			continue
		}
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			log.Printf("⚠️ Failed to start yt-dlp: %v", err)
			continue
		}

		scanner := bufio.NewScanner(stdoutPipe)
		var bar *pb.ProgressBar

		for scanner.Scan() {
			line := scanner.Text()
			line = sanitizeOutputLine(line)
			fmt.Println(line) // sanitized output

			// Progress bar logic
			if strings.HasPrefix(line, "[download]") {
				parts := strings.Fields(line)
				if len(parts) < 2 {
					continue
				}

				percentStr := parts[1]
				if strings.HasSuffix(percentStr, "%") {
					percentStr = strings.TrimSuffix(percentStr, "%")
					percent, err := strconv.Atoi(percentStr)
					if err != nil {
						continue
					}

					if bar == nil {
						bar = pb.New(100)
						bar.SetTemplateString(`{{counters . }} {{bar . }} {{percent . }} {{etime . }}`)
						bar.Start()
					}
					bar.SetCurrent(int64(percent))
				}
			}
		}

		if bar != nil {
			bar.Finish()
		}

		if err := scanner.Err(); err != nil {
			log.Printf("⚠️ Scanner error: %v", err)
		}

		if err := cmd.Wait(); err != nil {
			log.Printf("⚠️ yt-dlp failed: %v", err)
			continue
		}

		fmt.Printf("✅ Finished downloading: %s\n\n", url)
	}
}

// sanitizeOutputLine replaces illegal filename characters
func sanitizeOutputLine(line string) string {
	illegal := regexp.MustCompile(`[<>:"/\\|?*]`)
	return illegal.ReplaceAllString(line, "_")
}

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
