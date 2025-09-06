package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	urlList := flag.String("urls", "", "Comma-separated YouTube URLs or playlist URLs")
	audioOnly := flag.Bool("audio", false, "Download audio only")
	outputDir := flag.String("out", "downloads", "Output directory")
	quality := flag.String("quality", "best", "Quality: best, worst, 720p, 1080p, etc.")
	flag.Parse()

	if *urlList == "" {
		log.Fatal("❌ Please provide at least one YouTube URL using -urls")
	}

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
			fmt.Println(line) // print full yt-dlp output as well

			// Progress lines start with [download]
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
