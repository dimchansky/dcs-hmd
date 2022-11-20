package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	dcshmd "github.com/dimchansky/dcs-hmd"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

/*
|================================|
How to make a program Un-Clickable or Click-through-able?

#c::
    WinSet, ExStyle, ^0x20, A
    WinSet, Transparent, 255, A
    return


^!x::       ; Make mouse-click transparent
  WinGet, currentWindow, ID, A
  WinSet, ExStyle, +0x80020, ahk_id %currentWindow%
Return

^!z::       ; Undo mouse-click transparent
  WinSet, ExStyle, -0x20, ahk_id %currentWindow%
Return

https://jacks-autohotkey-blog.com/2016/07/22/the-winset-exstyle-command-for-mouse-click-transparent-windows-intermediate-autohotkey-tip/
|========================================|
*/

func run() error {
	hud, err := dcshmd.NewHUD()
	if err != nil {
		return err
	}
	defer func() {
		_ = hud.Close()
	}()

	if err := ebiten.RunGame(hud); err != nil {
		return err
	}

	return nil
}
