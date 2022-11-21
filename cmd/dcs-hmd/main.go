package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	dcshmd "github.com/dimchansky/dcs-hmd"
	"github.com/dimchansky/dcs-hmd/updlistener"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	hud, err := dcshmd.NewHUD()
	if err != nil {
		return err
	}
	defer func() {
		_ = hud.Close()
	}()

	l, err := updlistener.New(19089)
	if err != nil {
		return err
	}
	defer func() {
		_ = l.Close()
	}()

	if err := ebiten.RunGame(hud); err != nil {
		return err
	}

	return nil
}
