package main

import (
	"flag"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/silbinarywolf/preferdiscretegpu"

	dcshmd "github.com/dimchansky/dcs-hmd"
	"github.com/dimchansky/dcs-hmd/aircraft/ka-50/outputparser"
	"github.com/dimchansky/dcs-hmd/cmd"
	"github.com/dimchansky/dcs-hmd/updlistener"
)

func main() {
	var showVersion bool
	var targetDir string
	flag.BoolVar(&showVersion, "v", false, "show version information")
	flag.StringVar(&targetDir, "i", "", "target DCS scripts directory (usually %USERPROFILE%/Saved Games/DCS/ScriptsFS)")
	flag.Parse()

	if showVersion {
		fmt.Printf("Version: %s\nBuild Time: %s\nGit Hash: %s\n",
			cmd.Version, cmd.BuildTime, cmd.GitHash)
		return
	}

	if flag.Lookup("i") != nil {
		// check if -i flag is provided without a target directory name
		if targetDir == "" {
			fmt.Println("error: target DCS scripts directory is not specified")
			return
		}

		if err := dcshmd.InstallScripts(targetDir, true); err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Printf("all scripts are successfully installed to folder: %s", targetDir)
		fmt.Println()
		return
	}

	if err := run(); err != nil {
		fmt.Println("error:", err)
	}
}

const udpPortToListen = 19089

func run() error {
	hud, err := dcshmd.NewHUD()
	if err != nil {
		return fmt.Errorf("failed to create HUD: %w", err)
	}

	defer func() {
		_ = hud.Close()
	}()

	l, err := updlistener.New(udpPortToListen, outputparser.New(hud))
	if err != nil {
		return fmt.Errorf("failed to create UDP listener: %w", err)
	}

	defer func() {
		_ = l.Close()
	}()

	if err := ebiten.RunGame(hud); err != nil {
		return fmt.Errorf("failed to run HUD: %w", err)
	}

	return nil
}
