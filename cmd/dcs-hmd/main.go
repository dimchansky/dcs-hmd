package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/silbinarywolf/preferdiscretegpu"

	dcshmd "github.com/dimchansky/dcs-hmd"
	"github.com/dimchansky/dcs-hmd/aircraft/ka-50/outputparser"
	"github.com/dimchansky/dcs-hmd/updlistener"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
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
