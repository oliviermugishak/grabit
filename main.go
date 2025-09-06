package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

type PlaylistJSON struct {
	Entries []struct {
		ID string `json:"id"`
	} `json:"entries"`
}

func main() {
	// CLI flags
	urlList := flag.String("urls", "", "Comma-separated YouTube URLs or playlist URLs")
	audioOnly := flag.Bool("audio", false, "Download audio only")
	outputDir := flag.String("out", "downloads", "Output directory")
	quality := flag.String("quality", "best", "Quality: best, worst, 720p, 1080p, etc.")
	showVersion := flag.Bool("version", false, "Show Grabit version")
	concurrency := flag.Int("c", 3, "Number of concurrent downloads")
	flag.Parse()

	// Banner and version
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

	// Prepare all video URLs (split playlists)
	allVideos := []string{}
	for _, url := range strings.Split(*urlList, ",") {
		url = strings.TrimSpace(url)
		videos, err := extractVideos(url)
		if err != nil {
			log.Printf("⚠️ Failed to extract videos from %s: %v", url, err)
			continue
		}
		allVideos = append(allVideos, videos...)
	}

	// Worker pool
	var wg sync.WaitGroup
	videoChan := make(chan string)

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for videoURL := range videoChan {
				downloadVideo(workerID, videoURL, *outputDir, *audioOnly, *quality)
			}
		}(i + 1)
	}

	for _, v := range allVideos {
		videoChan <- v
	}
	close(videoChan)
	wg.Wait()
}

// extractVideos returns a list of individual video URLs
func extractVideos(url string) ([]string, error) {
	cmd := exec.Command("yt-dlp", "--flat-playlist", "-J", url)
	out, err := cmd.Output()
	if err != nil {
		// If not a playlist, return the URL itself
		return []string{url}, nil
	}

	var playlist PlaylistJSON
	if err := json.Unmarshal(out, &playlist); err != nil {
		return nil, err
	}

	videoURLs := []string{}
	for _, entry := range playlist.Entries {
		videoURLs = append(videoURLs, "https://www.youtube.com/watch?v="+entry.ID)
	}
	return videoURLs, nil
}

// downloadVideo downloads a single video with progress bar
func downloadVideo(workerID int, url, outputDir string, audioOnly bool, quality string) {
	fmt.Printf("Worker %d downloading: %s\n", workerID, url)

	args := []string{"-o", fmt.Sprintf("%s/%%(title)s.%%(ext)s", outputDir), "--newline"}

	if audioOnly {
		args = append(args, "-f", "bestaudio", "--extract-audio", "--audio-format", "m4a")
	} else {
		if quality == "best" || quality == "worst" {
			args = append(args, "-f", quality)
		} else {
			args = append(args, "-f", fmt.Sprintf("bestvideo[height<=%s]+bestaudio/best", quality))
		}
	}

	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("⚠️ Worker %d failed to get stdout pipe: %v", workerID, err)
		return
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Printf("⚠️ Worker %d failed to start yt-dlp: %v", workerID, err)
		return
	}

	scanner := bufio.NewScanner(stdoutPipe)
	var bar *pb.ProgressBar

	for scanner.Scan() {
		line := sanitizeOutputLine(scanner.Text())
		fmt.Printf("Worker %d: %s\n", workerID, line)

		if strings.HasPrefix(line, "[download]") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}

			percentStr := strings.TrimSuffix(parts[1], "%")
			if percent, err := strconv.Atoi(percentStr); err == nil {
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
		log.Printf("⚠️ Worker %d scanner error: %v", workerID, err)
	}
	if err := cmd.Wait(); err != nil {
		log.Printf("⚠️ Worker %d yt-dlp failed: %v", workerID, err)
	}

	fmt.Printf("✅ Worker %d finished downloading: %s\n", workerID, url)
}

// sanitizeOutputLine replaces illegal filename characters
func sanitizeOutputLine(line string) string {
	illegal := regexp.MustCompile(`[<>:"/\\|?*]`)
	return illegal.ReplaceAllString(line, "_")
}
