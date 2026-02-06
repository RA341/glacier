package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var workingDir = ""

func main() {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	var err error
	workingDir, err = os.Getwd()
	must(err)

	frostOut := filepath.Join(workingDir, ".build", "frost")
	err = os.MkdirAll(frostOut, 0777)
	must(err)

	uiDir := filepath.Join(workingDir, "ui")

	electronOSBuildFlag := "-l"
	electronOut := "linux-unpacked"
	osOut := "linux"
	frostExe := "frost"
	if runtime.GOOS == "windows" {
		frostExe += ".exe"
		electronOSBuildFlag = "-w"
		electronOut = "win-unpacked"
		osOut = "windows"
	}

	frostOut = filepath.Join(frostOut, osOut)
	frostOutExe := filepath.Join(frostOut, frostExe)
	uiFrostOut := filepath.Join(frostOut, "ui")

	jobs := []*Job{
		{
			ID: 1,
			Stages: []Stage{
				{
					Dir: workingDir,
					Cmd: "task",
					Args: []string{
						"ui:b:frost",
					},
				},
				{
					Dir: workingDir,
					Cmd: "task",
					Args: []string{
						"go:b:frost", "PROD=1", "OUT_EXE=" + frostOutExe,
					},
				},
			},
		},
		{
			ID: 2,
			Stages: []Stage{
				{
					Dir:     uiDir,
					Cmd:     "npm",
					Sources: []string{"package*.json", "electron.cjs"},
					Args: []string{
						"run",
						"desk:build",
						"--",
						electronOSBuildFlag,
					},
					Debug: true,
				},
				{
					Dir: uiDir,
					RunFn: func() error {
						return CopyDir(
							filepath.Join(uiDir, "release/"+electronOut),
							uiFrostOut,
						)
					},
				},
			},
		},
		// todo
		//{
		//	ID: 1,
		//	Stages: []Stage{
		//		{
		//			Dir:  uiDir,
		//			Cmd:  "npm",
		//			Args: []string{"run", "buildfrost"},
		//		},
		//		{
		//			Dir:  uiDir,
		//			Cmd:  "cp",
		//			Args: []string{"-r", "build", frostCmd},
		//		},
		//		{
		//			Dir:       workingDir,
		//			Cmd:       gitCmd,
		//			IgnoreErr: true,
		//			Args: []string{
		//				"describe",
		//				"--tags",
		//				"--abbrev=0",
		//			},
		//			OnDone: func(output string) {
		//				if output != "" {
		//					version = output
		//				}
		//			},
		//		},
		//		{
		//			Dir:       workingDir,
		//			Cmd:       gitCmd,
		//			IgnoreErr: true,
		//			Args: []string{
		//				"rev-parse",
		//				"HEAD",
		//			},
		//			Debug: true,
		//			OnDone: func(output string) {
		//				if output != "" {
		//					commit = output
		//				}
		//			},
		//		},
		//		{
		//			Dir:       workingDir,
		//			Cmd:       gitCmd,
		//			IgnoreErr: true,
		//			Debug:     true,
		//			Args: []string{
		//				"rev-parse", "--abbrev-ref HEAD",
		//			},
		//			OnDone: func(output string) {
		//				if output != "" {
		//					branch = output
		//				}
		//			},
		//		},
		//		{
		//			Dir: core,
		//			Cmd: "go",
		//			Sources: []string{
		//				"**/*.go",
		//				"go.*",
		//			},
		//			Args: []string{
		//				"build",
		//				"-o", frostOut,
		//				"-ldflags",
		//				fmt.Sprintf("s -w %s %s %s %s",
		//					withPkgInf("Version", version),
		//					withPkgInf("CommitInfo", commit),
		//					withPkgInf("BuildDate", time.Now().UTC().Format(time.RFC3339)),
		//					withPkgInf("Branch", branch),
		//				),
		//				"./cmd/frost",
		//			},
		//		},
		//	},
		//},
	}

	startJobs(jobs)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
