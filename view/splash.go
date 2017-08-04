package view

import (
	ui "github.com/aiportal/termui"
)

func rowBlank(height int) *ui.Row {
	blank := ui.NewBlock()
	blank.Height = height
	blank.Border = false
	return ui.NewRow(ui.NewCol(12, 0, blank))
}

func rowTitle() *ui.Row {
	const splashTitleFmt = "                \n    GoWallet    \n                "

	title := ui.NewPar(splashTitleFmt)
	title.Height = 4
	title.Border = false
	title.TextFgColor = ui.ColorYellow | ui.AttrBold
	title.TextBgColor = ui.ColorBlue
	return ui.NewRow(ui.NewCol(3, 5, title))
}

func rowEnter() *ui.Row {
	const splashEnterFmt  = "  Press <ENTER> to continue."

	enter := ui.NewPar(splashEnterFmt)
	enter.Height = 2
	enter.Border = false
	enter.TextFgColor = ui.ColorRed | ui.AttrBold
	return ui.NewRow(ui.NewCol(12, 0, enter))
}

func rowStartTip() *ui.Row {
	const splashTipFmt = `
  GoWallet is a safe brain wallet for bitcoin.
  It uses a secret phrase and a salt phrase to generate wallets.

  It is recommended that use a more complex secret and put it on paper.
  It's also recommended that keep your salt in mind.`

	tip := ui.NewPar(splashTipFmt)
	tip.Height = 8
	tip.Border = false
	tip.TextFgColor = ui.ColorGreen | ui.AttrBold
	return ui.NewRow(ui.NewCol(12, 0, tip))
}

func rowProjectUrl() *ui.Row {
	const splashAddrFmt  = "  Project address: [https://github.com/aiportal/gowallet](fg-bold,fg-underline)"

	addr := ui.NewPar(splashAddrFmt)
	addr.Height = 2
	addr.Border = false
	addr.TextFgColor = ui.ColorWhite
	return ui.NewRow(ui.NewCol(12, 0, addr))
}

func rowCreateTip() *ui.Row {
	const splashCreateTipFmt = `
  Secret phrase at least 16 characters.
  Salt phrase at least 6 characters.

  Secret phrases should contain uppercase letters, lowercase letters, numbers, and special characters.

  Both secret phrases and salt phrases can use hexadecimal notation such as \xff or \xFF to represent a character.`

	tip := ui.NewPar(splashCreateTipFmt)
	tip.Height = 10
	tip.Border = false
	tip.TextFgColor = ui.ColorGreen | ui.AttrBold
	return ui.NewRow(ui.NewCol(12, 0, tip))
}

func rowCreateEnter() *ui.Row {
	const splashCreateNextFmt  = "  Press <ENTER> to create a wallet."

	next := ui.NewPar(splashCreateNextFmt)
	next.Height = 2
	next.Border = false
	next.TextFgColor = ui.ColorRed | ui.AttrBold
	return ui.NewRow(ui.NewCol(12, 0, next))
}

type SplashView []*ui.Row

var SplashStartView = SplashView{ rowBlank(2), rowTitle(), rowStartTip(), rowProjectUrl(), rowEnter() }
var SplashCreateView = SplashView{ rowBlank(2), rowTitle(), rowCreateTip(), rowCreateEnter() }

func ShowSplashView(rows SplashView) {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// build layout
	ui.Body.AddRows(rows...)

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

	ui.Loop()
}
