package main

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	Exe()
}

const (
	cleanArg = "--clean"

	BuildFile = "build.go"
	NogFile   = "nog.go"

	NogExec = "nog"
)

var (
	exe        = ""
	NogExecOld = NogExec + ".old" + exe
)

func Exe() {
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}
}

func MustStat(path string) os.FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		Fatal(
			"Failed to stat",
			"err", err,
			"path", path,
		)
	}
	return info
}

func GoRebuildUrself(opts ...Opt) {
	parseOpts(opts...)

	if len(os.Args) > 1 && os.Args[1] == cleanArg {
		cleanupAfterUrSelf()
		os.Exit(0)
	}

	executablePath := GetExecutable()
	rebuild := shouldRebuild(executablePath)
	if rebuild {
		slog.Info("rebuilding myself")
		BuildExec(executablePath)
	} else {
		slog.Info("up to date")
	}
}

type RebuildOpts struct {
	// set working dir to look for build.go and nob.go
	workingDir string
}

var rebuildOptions RebuildOpts

func parseOpts(opts ...Opt) {
	finalOpts := RebuildOpts{}
	for _, opt := range opts {
		opt(&finalOpts)
	}
	if finalOpts.workingDir == "" {
		finalOpts.workingDir = "./"
	}

	rebuildOptions = finalOpts
}

func WithWD(path string) Opt {
	return func(opts *RebuildOpts) {
		opts.workingDir = path
	}
}

type Opt func(opts *RebuildOpts)

func shouldRebuild(executablePath string) bool {
	execStat, err := os.Stat(executablePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		Fatal("Failed to stat executablePath: %v", err)
	}

	nogStat := MustStat(NogFile)
	buildStat := MustStat(BuildFile)

	isBuildMod := buildStat.ModTime().After(execStat.ModTime())
	isNogMod := nogStat.ModTime().After(execStat.ModTime())
	return isBuildMod || isNogMod
}

func cleanupAfterUrSelf() {
	path := getOldExecName()
	slog.Info("Cleaning up after MySelf", "path", path)
	err := os.RemoveAll(path)
	if err != nil {
		slog.Error("Failed to remove old executable: %v", err)
		return
	}
}

func GetExecutable() string {
	return rebuildOptions.workingDir + NogExec + "" + exe
}

func BuildExec(executablePath string) {
	rename := getOldExecName()

	slog.Info("Renaming old build exec",
		"from", executablePath,
		"to", rename,
	)
	err := os.Rename(executablePath, rename)
	if err != nil {
		Fatal("Failed to rename path",
			"err", err,
			"from", executablePath,
			"to", rename,
		)
	}

	cmd := []string{
		"go",
		"build",
		"-o", executablePath,
		NogFile,
		BuildFile,
	}
	RunCmd(cmd[0], cmd[1:]...)

	slog.Info("Calling new executable", "path", executablePath)
	RunCmd(executablePath, os.Args[1:]...)

	slog.Debug("Removing old exec")

	slog.Info("Goodbye...")
	// call the new binary to remove the old binary
	// calls ./nog --clean
	StartCmdSilent(executablePath, cleanArg)

	os.Exit(0)
}

func getOldExecName() string {
	return rebuildOptions.workingDir + NogExecOld + exe
}

// StartCmdSilent runs exec.Run() with no attached streams
func StartCmdSilent(executable string, args ...string) {
	execCmd := exec.Command(executable, args...)

	err := execCmd.Start()
	if err != nil {
		Fatal("Unable to run cmd: %v, %v", args, err)
	}
}

func RunCmd(executable string, args ...string) {
	execCmd := exec.Command(executable, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	err := execCmd.Run()
	if err != nil {
		Fatal("Unable to run cmd: %v, %v", args, err)
	}
}

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
