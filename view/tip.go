package view

import (
	"fmt"
	"strings"

	ui "github.com/aiportal/termui"
)

func NewTipView(name string) (view *GridView) {

	if strings.HasPrefix(name, "1") {
		view = createTransTipView(name)
		return
	}
	switch name {
	case "vanity":
		view = createVanityTipView()
	case "password":
		view = createPasswordTipView()
	case "spend":
		view = createSpendTipView()
	case "export":
		view = createExportTipView()
	}
	return
}

func createTransTipView(address string) (view *GridView) {
	var splashTipFmt = fmt.Sprintf(`
  List transactions for address: [%s](bg-blue,fg-underline).

  This feature is not yet implemented.

  Your donation can accelerate the progress of project development.`, address)

	tip := ui.NewPar(splashTipFmt)
	tip.Height = 8
	tip.Border = false
	tip.TextFgColor = ui.ColorGreen | ui.AttrBold

	rowTip := ui.NewRow(ui.NewCol(12, 0, tip))
	view = new(GridView)
	view.rows = []*ui.Row{rowBlank(2), rowTitle(), rowTip, rowProjectUrl(), rowEnter()}
	return
}

func createVanityTipView() (view *GridView) {
	const splashTipFmt = `
  Find [vainty address](bg-blue) in child wallets.

  This feature is not yet implemented.

  Your donation can accelerate the progress of project development.`

	tip := ui.NewPar(splashTipFmt)
	tip.Height = 8
	tip.Border = false
	tip.TextFgColor = ui.ColorGreen | ui.AttrBold

	rowTip := ui.NewRow(ui.NewCol(12, 0, tip))
	view = new(GridView)
	view.rows = []*ui.Row{rowBlank(2), rowTitle(), rowTip, rowProjectUrl(), rowEnter()}
	return
}

func createPasswordTipView() (view *GridView) {
	const splashTipFmt = `
  Create or change [login password](bg-blue) for GoWallet.

  This feature is not yet implemented.

  Your donation can accelerate the progress of project development.`

	tip := ui.NewPar(splashTipFmt)
	tip.Height = 8
	tip.Border = false
	tip.TextFgColor = ui.ColorGreen | ui.AttrBold

	rowTip := ui.NewRow(ui.NewCol(12, 0, tip))
	view = new(GridView)
	view.rows = []*ui.Row{rowBlank(2), rowTitle(), rowTip, rowProjectUrl(), rowEnter()}
	return
}

func createSpendTipView() (view *GridView) {
	const splashTipFmt = `
  [Send Bitcoin](bg-blue) to other addresses.

  This feature is not yet implemented.

  Your donation can accelerate the progress of project development.`

	tip := ui.NewPar(splashTipFmt)
	tip.Height = 8
	tip.Border = false
	tip.TextFgColor = ui.ColorGreen | ui.AttrBold

	rowTip := ui.NewRow(ui.NewCol(12, 0, tip))
	view = new(GridView)
	view.rows = []*ui.Row{rowBlank(2), rowTitle(), rowTip, rowProjectUrl(), rowEnter()}
	return
}

func createExportTipView() (view *GridView) {
	const splashTipFmt = `
  Export Wallets [private key and address](bg-blue) to a plain text file.

  This feature is not yet implemented.

  Your donation can accelerate the progress of project development.`

	tip := ui.NewPar(splashTipFmt)
	tip.Height = 8
	tip.Border = false
	tip.TextFgColor = ui.ColorGreen | ui.AttrBold

	rowTip := ui.NewRow(ui.NewCol(12, 0, tip))
	view = new(GridView)
	view.rows = []*ui.Row{rowBlank(2), rowTitle(), rowTip, rowProjectUrl(), rowEnter()}
	return
}
