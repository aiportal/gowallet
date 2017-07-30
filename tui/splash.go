// +build en_us !zh_cn

package tui

import (
	ui "github.com/aiportal/termui"
)

const (
	//-------------------------------------------------------------------------
	// Splash start view
	//-------------------------------------------------------------------------
	splashTitleFmt  = "                \n    GoWallet    \n                "
	splashTitleRows = 4

	splashTipFmt = `
  GoWallet is a safe brain wallet for bitcoin.
  It uses a secret phrase and a salt phrase to generate wallets.

  It is recommended that use a more complex secret and put it on paper.
  It's also recommended that keep your salt in mind.`
	splashTipRows = 8

	splashAddrFmt  = "  Project address: [https://github.com/aiportal/gowallet](fg-bold,fg-underline)"
	splashAddrRows = 2

	splashEnterFmt  = "  Press <ENTER> to continue."
	splashEnterRows = 2

	//-------------------------------------------------------------------------
	// Splash create view
	//-------------------------------------------------------------------------
	splashCreateTipFmt = `
  Secret phrase at least 16 characters.
  Salt phrase at least 6 characters.

  Secret phrases should contain uppercase letters, lowercase letters, numbers, and special characters.

  Both secret phrases and salt phrases can use hexadecimal notation such as \xff or \xFF to represent a character.`
	splashCreateTipRows = 10

	splashCreateNextFmt  = "  Press <ENTER> to create a wallet."
	splashCreateNextRows = 2
)

func SplashStart() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	blank := ui.NewBlock()
	blank.Border = false

	title := ui.NewPar(splashTitleFmt)
	title.Height = splashTitleRows
	title.Border = false
	title.TextFgColor = ui.ColorYellow | ui.AttrBold
	title.TextBgColor = ui.ColorBlue

	tip := ui.NewPar(splashTipFmt)
	tip.Height = splashTipRows
	tip.Border = false
	tip.TextFgColor = ui.ColorGreen | ui.AttrBold

	addr := ui.NewPar(splashAddrFmt)
	addr.Height = splashAddrRows
	addr.Border = false
	addr.TextFgColor = ui.ColorWhite

	enter := ui.NewPar(splashEnterFmt)
	enter.Height = splashEnterRows
	enter.Border = false
	enter.TextFgColor = ui.ColorRed | ui.AttrBold

	// build layout
	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(3, 5, blank)),
		ui.NewRow(ui.NewCol(3, 5, title)),
		ui.NewRow(ui.NewCol(12, 0, tip)),
		ui.NewRow(ui.NewCol(12, 0, addr)),
		ui.NewRow(ui.NewCol(12, 0, enter)))

	ui.Body.Align()
	ui.Clear()
	ui.Render(ui.Body)

	ui.Handle("sys/kbd/<enter>", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		//ui.Body.Width = ui.TermWidth()
		//ui.Body.Align()
		//ui.Clear()
		//ui.Render(ui.Body)
	})

	ui.Loop()
}

func SplashCreate() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	blank := ui.NewBlock()
	blank.Border = false

	title := ui.NewPar(splashTitleFmt)
	title.Height = splashTitleRows
	title.Border = false
	title.TextFgColor = ui.ColorYellow | ui.AttrBold
	title.TextBgColor = ui.ColorBlue

	tip := ui.NewPar(splashCreateTipFmt)
	tip.Height = splashCreateTipRows
	tip.Border = false
	tip.TextFgColor = ui.ColorGreen | ui.AttrBold

	next := ui.NewPar(splashCreateNextFmt)
	next.Height = splashCreateNextRows
	next.Border = false
	next.TextFgColor = ui.ColorRed | ui.AttrBold

	// build layout
	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(3, 5, blank)),
		ui.NewRow(ui.NewCol(3, 5, title)),
		ui.NewRow(ui.NewCol(12, 0, tip)),
		ui.NewRow(ui.NewCol(12, 0, next)))

	ui.Body.Align()
	ui.Clear()
	ui.Render(ui.Body)

	ui.Handle("sys/kbd/<enter>", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		//ui.Body.Width = ui.TermWidth()
		//ui.Body.Align()
		//ui.Clear()
		//ui.Render(ui.Body)
	})

	ui.Loop()
}
