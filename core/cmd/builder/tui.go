package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

func renderLoop(ctx context.Context, jobs []*Job, doneCh <-chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-doneCh:
			render(jobs)
			fmt.Println("\n\033[1;32mALL DONE\033[0m")
			fmt.Println("Summary:")

			for _, j := range jobs {
				j.mu.RLock()

				// Calculate total job time
				var totalDuration time.Duration
				for _, s := range j.Stages {
					totalDuration += s.Duration
				}

				fmt.Printf("\n\033[1mJob %d Total Time: %s\033[0m\n", j.ID, totalDuration.Round(time.Millisecond))

				for i, s := range j.Stages {
					dur := s.Duration.Truncate(time.Millisecond).String()

					// Handle failed/skipped status
					if s.Duration == 0 {
						if j.status == "failed" && i >= j.CurrentStage-1 {
							dur = "\033[31mfailed/skipped\033[0m"
						} else {
							dur = "0s"
						}
					}

					fmt.Printf("  [Stage %d-%d] [%s %s] - %s\n",
						j.ID, i+1, s.Cmd, strings.Join(s.Args, " "), dur)

					if s.Debug {
						fmt.Printf("    \033[90m> Debug Output: %s\033[0m\n", strings.ReplaceAll(s.Output, "\n", "\n    "))
					}

				}
				j.mu.RUnlock()
			}
			return
		case <-ticker.C:
			render(jobs)
		case <-ctx.Done():
			return
		}
	}
}

func render(jobs []*Job) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24
	}

	numJobs := len(jobs)
	if numJobs == 0 {
		return
	}

	availLines := height - numJobs
	if availLines < 0 {
		availLines = 0
	}

	linesPerJob := availLines / numJobs
	remainder := availLines % numJobs

	var screenLines []string

	for i, j := range jobs {
		j.mu.RLock()

		paneHeight := linesPerJob
		if i == 0 {
			paneHeight += remainder
		}

		idx := j.CurrentStage - 1
		if idx < 0 {
			idx = 0
		}
		if idx >= len(j.Stages) {
			idx = len(j.Stages) - 1
		}

		currentStageSpec := j.Stages[idx]
		cmdDisplay := fmt.Sprintf("%s %s", currentStageSpec.Cmd, strings.Join(currentStageSpec.Args, " "))
		cmdLim := 30
		if len(cmdDisplay) > cmdLim {
			cmdDisplay = cmdDisplay[:cmdLim] + "..."
		}

		duration := time.Since(j.StartTime).Seconds()
		displayStageNum := j.CurrentStage
		if displayStageNum == 0 {
			displayStageNum = 1
		}

		header := fmt.Sprintf("[%6.1fs] [stage %d-%d] %s (running...)",
			duration, j.ID, j.CurrentStage, cmdDisplay)

		headerFormatted := fmt.Sprintf("\033[1m%s\033[0m", header)
		screenLines = append(screenLines, headerFormatted)

		count := len(j.Lines)
		start := 0
		if count > paneHeight {
			start = count - paneHeight
		}
		visibleLines := j.Lines[start:]

		for _, line := range visibleLines {
			if len(line) > width {
				line = line[:width]
			}
			screenLines = append(screenLines, fmt.Sprintf("=> %s", line))
		}

		padding := paneHeight - len(visibleLines)
		for k := 0; k < padding; k++ {
			screenLines = append(screenLines, "")
		}
		j.mu.RUnlock()
	}

	output := strings.Join(screenLines, "\n")
	fmt.Print("\033[H\033[2J" + output)
}
