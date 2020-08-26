package view

import (
	"fmt"
	"sync"
	"time"

	"github.com/aiportal/gowallet/wallet"
	ui "github.com/aiportal/termui"
)

type AccountView struct {
	GridView
	walletActiveRow int
	Data            interface{}
}

func NewAccountView(ws []*wallet.Wallet) (v *AccountView) {
	v = new(AccountView)
	v.controls = []*GridViewControl{
		v.ctrlTitleBar(),
		v.ctrlWallets(ws),
		v.ctrlMenu(),
	}
	v.rows = []*ui.Row{
		ui.NewRow(ui.NewCol(12, 0, v.controls[0].Element.(ui.GridBufferer))),
		ui.NewRow(ui.NewCol(12, 0, v.controls[1].Element.(ui.GridBufferer))),
		ui.NewRow(ui.NewCol(12, 0, v.controls[2].Element.(ui.GridBufferer))),
	}
	return
}

func (v *AccountView) ctrlTitleBar() *GridViewControl {
	c := ui.NewGauge()
	c.Label = "GoWallet"
	c.Percent = 100
	c.Height = 3
	c.Border = false
	c.BarColor = ui.ColorBlue
	c.PercentColorHighlighted = ui.ColorYellow | ui.AttrBold
	return NewGridViewControl("title", c, nil)
}

// {Type:keyboard Path:/sys/kbd/v From:termbox To: Data:{KeyStr:v} Time:1501936747}
// {Type:window Path:/sys/wnd/resize From:termbox To: Data:{Width:97 Height:18} Time:1501937946}
func (v *AccountView) ctrlWallets(ws []*wallet.Wallet) *GridViewControl {
	rows := make([][]string, len(ws)+2)
	rows[0] = []string{"No", "Address", "Balance"}
	for i, w := range ws {
		rows[i+1] = []string{
			fmt.Sprintf("%d", w.No),
			w.Address,
			"loading...",
		}
	}
	maxRow := len(rows) - 1
	rows[maxRow] = []string{"", "Total", "computing..."}
	time.AfterFunc(time.Second, func() {
		v.loadBalance(ws, rows[1:])
		v.Update()
	})

	activeRow := 1
	activeRowFg := ui.ColorWhite | ui.AttrBold
	activeRowBg := ui.ColorGreen

	t := ui.NewTable()
	t.Rows = rows
	t.Height = 14
	t.FgColor = ui.ColorGreen
	t.BgColor = ui.ColorDefault
	t.TextAlign = ui.AlignRight
	t.Separator = false
	t.Analysis()
	t.FgColors[0] = ui.ColorGreen | ui.AttrBold | ui.AttrUnderline
	t.FgColors[1] = activeRowFg
	t.BgColors[1] = activeRowBg
	t.FgColors[maxRow] = ui.ColorGreen | ui.AttrBold

	m := make(HandlerMap)
	m["sys/kbd/<up>"] = func(ui.Event) {
		count := len(t.FgColors) - 1
		for i := 1; i < count; i++ {
			t.FgColors[i] = t.FgColor
			t.BgColors[i] = t.BgColor
		}
		if activeRow > 1 {
			activeRow -= 1
		}
		t.FgColors[activeRow] = activeRowFg
		t.BgColors[activeRow] = activeRowBg
		v.Update()
	}
	m["sys/kbd/<down>"] = func(ui.Event) {
		count := len(t.FgColors) - 2
		for i := 1; i < count; i++ {
			t.FgColors[i] = t.FgColor
			t.BgColors[i] = t.BgColor
		}
		if activeRow < count {
			activeRow += 1
		}
		t.FgColors[activeRow] = activeRowFg
		t.BgColors[activeRow] = activeRowBg
		v.Update()
	}
	m["sys/kbd/<enter>"] = func(ui.Event) {
		w := ws[activeRow-1]
		v.Data = w.Address
		ui.StopLoop()
	}
	return NewGridViewControl("wallets", t, m)
}

func (v *AccountView) loadBalance(ws []*wallet.Wallet, rows [][]string) {
	var wg sync.WaitGroup
	for i, w := range ws {
		wg.Add(1)
		go func(i int, w *wallet.Wallet) {
			defer wg.Done()
			w.LoadBalance()
			if w.Balance != nil {
				if w.Balance.Value() > 0 {
					rows[i][2] = fmt.Sprintf("%s", w.Balance.FmtBtc())
				} else {
					rows[i][2] = fmt.Sprintf("%s", w.Balance.FmtBtc())
				}
			} else {
				rows[i][2] = "Unkown"
			}
			v.Update()
		}(i, w)
	}
	wg.Wait()

	total := uint64(0)
	for _, w := range ws {
		if w.Balance != nil {
			total += w.Balance.Value()
		}
	}
	balance := wallet.NewBalance(total)
	rows[len(rows)-1][2] = balance.FmtBtc()
	v.Update()
}

func (v *AccountView) ctrlMenu() *GridViewControl {

	const AccountMenuFmt = "" +
		"  [V](fg-bold,fg-underline)anity  " +
		"  [S](fg-bold,fg-underline)pend  " +
		"  [P](fg-bold,fg-underline)assword  " +
		"  [E](fg-bold,fg-underline)xport  " +
		"  [Q](fg-bold,fg-underline)uit  "

	c := ui.NewPar(AccountMenuFmt)
	c.BorderTop = false
	c.BorderBottom = false
	c.Height = 1
	c.TextFgColor = ui.ColorMagenta
	c.TextBgColor = ui.ColorBlue
	c.Bg = ui.ColorBlue

	m := make(HandlerMap)
	m["sys/kbd/v"] = func(ui.Event) {
		v.Data = "vanity"
		ui.StopLoop()
	}
	m["sys/kbd/s"] = func(ui.Event) {
		v.Data = "spend"
		ui.StopLoop()
	}
	m["sys/kbd/p"] = func(ui.Event) {
		v.Data = "password"
		ui.StopLoop()
	}
	m["sys/kbd/e"] = func(ui.Event) {
		v.Data = "export"
		ui.StopLoop()
	}
	m["sys/kbd/q"] = func(ui.Event) {
		v.Data = "quit"
		ui.StopLoop()
	}
	return NewGridViewControl("menu", c, m)
}
