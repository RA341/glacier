package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Stage represents a single command step within a job pipeline
type Stage struct {
	Cmd       string
	Args      []string
	Dir       string
	Duration  time.Duration
	OnDone    func(output string)
	IgnoreErr bool
	Sources   []string
	Debug     bool
	Output    string // Add this to store the result
	RunFn     func() error
}

// Job represents a lane in the TUI that runs multiple stages sequentially
type Job struct {
	ID           int
	Stages       []Stage
	CurrentStage int
	StartTime    time.Time
	Lines        []string

	mu     sync.RWMutex
	done   bool
	status string
}

func startJobs(jobs []*Job) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		j.StartTime = time.Now()
		go func(job *Job) {
			defer wg.Done()
			runJob(ctx, job)
		}(j)
	}

	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	renderLoop(ctx, jobs, doneCh)
}

func runJob(ctx context.Context, job *Job) {
	allPreviousCached := true

	for i := range job.Stages {
		job.mu.Lock()
		job.CurrentStage = i + 1
		job.status = "running"
		job.mu.Unlock()

		start := time.Now()
		var err error
		stage := job.Stages[i]
		stageKey := fmt.Sprintf("%d-%d", job.ID, i)

		// isCurrentCached is true if file exists AND (no sources OR hashes match)
		isCurrentCached, newHashContent := checkCache(stageKey, &stage)

		var resultSuffix string
		var output string
		if isCurrentCached && allPreviousCached {
			job.status = "cached"
			resultSuffix = "\033[33m(CACHED)\033[0m"
		} else {
			allPreviousCached = false

			output, err = runStage(ctx, job, &stage)
			if err == nil {
				job.status = "done"
				resultSuffix = "\033[32m(DONE)\033[0m"
				// Mark as completed in cache, even if newHashContent is empty
				commitCache(stageKey, newHashContent)
			} else {
				job.status = "failed"
				resultSuffix = "\033[31m(FAILED)\033[0m"
			}
		}

		// TUI Logging
		elapsed := time.Since(start)
		durationTotal := time.Since(job.StartTime).Seconds()
		cmdStr := fmt.Sprintf("%s %s", stage.Cmd, strings.Join(stage.Args, " "))
		lim := 80
		if len(cmdStr) > lim {
			cmdStr = cmdStr[:lim] + "..."
		}

		completionLine := fmt.Sprintf("[%6.1fs] [stage %d-%d] %s %s",
			durationTotal, job.ID, i+1, cmdStr, resultSuffix)

		job.mu.Lock()
		job.Stages[i].Duration = elapsed
		job.Stages[i].Output = output
		job.Lines = append(job.Lines, completionLine)
		job.mu.Unlock()

		if err != nil {
			return
		}
	}

	job.mu.Lock()
	job.done = true
	job.mu.Unlock()
}

func runStage(ctx context.Context, job *Job, stage *Stage) (string, error) {
	if stage.RunFn != nil {
		return "", stage.RunFn()
	}

	cmd := exec.CommandContext(ctx, stage.Cmd, stage.Args...)
	if stage.Dir != "" {
		cmd.Dir = stage.Dir
	}

	cmd.Env = append(os.Environ(),
		"FORCE_COLOR=1",
		"CLICOLOR_FORCE=1",
		"TERM=xterm-256color",
	)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	multiReader := io.MultiReader(stdout, stderr)

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var capturedOutput strings.Builder
	scanner := bufio.NewScanner(multiReader)
	for scanner.Scan() {
		txt := scanner.Text()
		job.appendLine(txt)
		// Always capture if Debug is on or OnDone is present
		if stage.OnDone != nil || stage.Debug {
			capturedOutput.WriteString(txt + "\n")
		}
	}

	err := cmd.Wait()
	output := strings.TrimSpace(capturedOutput.String())

	if stage.OnDone != nil {
		stage.OnDone(output)
	}

	if stage.IgnoreErr {
		return output, nil
	}

	return output, err
}

const cacheDir = ".bob.cache"

func checkCache(i string, stage *Stage) (bool, string) {
	base := filepath.Join(workingDir, cacheDir)
	cachePath := filepath.Join(base, i)

	oldContent, err := os.ReadFile(cachePath)
	if err != nil {
		return false, ""
	}

	var currentLines []string
	currentLines = append(currentLines, fmt.Sprintf("CMD: %s %s", stage.Cmd, strings.Join(stage.Args, " ")))

	files := expandSources(stage.Dir, stage.Sources)
	for _, fpath := range files {
		stat, err := os.Stat(fpath)
		if err != nil {
			continue
		}
		relPath, _ := filepath.Rel(stage.Dir, fpath)
		currentLines = append(currentLines, fmt.Sprintf("%s:%s", relPath, stat.ModTime().String()))
	}

	newContent := strings.Join(currentLines, "\n")
	return string(oldContent) == newContent, newContent
}

func commitCache(i string, content string) {
	base := filepath.Join(workingDir, cacheDir)
	err := os.MkdirAll(base, 0777)
	must(err)
	// Write hashes or an empty string to signify completion
	err = os.WriteFile(filepath.Join(base, i), []byte(content), 0644)
	must(err)
}

func expandSources(wd string, patterns []string) []string {
	var allFiles []string
	for _, pattern := range patterns {
		// Handle recursive pattern **
		if strings.Contains(pattern, "**") {
			baseDir := strings.Split(pattern, "**")[0]
			searchDir := filepath.Join(wd, baseDir)

			err := filepath.WalkDir(searchDir, func(path string, d os.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					// Match the specific extension or suffix if provided after **
					suffix := filepath.Ext(pattern)
					if suffix == "" || strings.HasSuffix(path, suffix) {
						allFiles = append(allFiles, path)
					}
				}
				return nil
			})
			must(err)

		} else {
			// Standard non-recursive glob
			matches, _ := filepath.Glob(filepath.Join(wd, pattern))
			for _, m := range matches {
				if stat, err := os.Stat(m); err == nil && !stat.IsDir() {
					allFiles = append(allFiles, m)
				}
			}
		}
	}
	return allFiles
}

func (j *Job) appendLine(line string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Lines = append(j.Lines, line)
	if len(j.Lines) > 2000 {
		j.Lines = j.Lines[len(j.Lines)-2000:]
	}
}
