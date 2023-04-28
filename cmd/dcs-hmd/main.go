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
	showVersion := flag.Bool("v", false, "show version information")
	installDir := flag.String("i", "", `install scripts to the target DCS scripts directory (usually "%USERPROFILE%\Saved Games\DCS.openbeta\Scripts")`)
	unInstallDir := flag.String("u", "", `uninstall scripts from the target DCS scripts directory (usually "%USERPROFILE%\Saved Games\DCS.openbeta\Scripts")`)
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\nBuild Time: %s\nGit Hash: %s\n",
			cmd.Version, cmd.BuildTime, cmd.GitHash)
		return
	}

	if *installDir != "" {
		if err := dcshmd.InstallScripts(*installDir, true); err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Printf("all scripts are successfully installed to the folder: %s", *installDir)
		fmt.Println()
		return
	}

	if *unInstallDir != "" {
		if err := dcshmd.UninstallScripts(*unInstallDir, true); err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Printf("all scripts were successfully uninstalled from the folder: %s", *unInstallDir)
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
