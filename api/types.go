//
// This file is part of otterchat.
//
// otterchat is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// otterchat is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with otterchat.  If not, see <https://www.gnu.org/licenses/>.
//

package api

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// OtterchatKey key for the main page
const (
	OtterchatKey  = "otterchat"
	CMDSymbolName = "Commands"
)

// Command represents an executable a command
type Command interface {
	Name() string
	Usage() string
	Desc() string
	Exec(*UI, *Page, []string) error
}

// Commands a plugin that contains one or more command
type Commands interface {
	Registry() map[string]Command
}

// Page The elements needed for the page
type Page struct {
	ChannelInfoView *tview.TextView
	ChatView        *tview.TextView
	InputInfoView   *tview.TextView
	UserInput       *tview.InputField
	UserList        *tview.List
}

// UI The User Inteface elements
type UI struct {
	App            *tview.Application
	AppView        *tview.Flex
	ConnectionList *tview.List
	Pages          *tview.Pages
	PageMap        map[string]*Page
}

// Init the User Interface
func Init() *UI {
	ui := new(UI)

	ui.PageMap = make(map[string]*Page)

	ui.App = tview.NewApplication()

	ui.Pages = tview.NewPages()

	ui.ConnectionList = tview.NewList()

	ui.ConnectionList.SetBorder(true).
		SetBorderColor(tcell.ColorBlack).
		SetTitle("Connections")

	ui.AppView = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ui.ConnectionList, 20, 1, false).
		AddItem(ui.Pages, 0, 2, true)

	return ui
}

func (ui *UI) handleInput(id string, page *Page, commands map[string]Command) error {
	if commands == nil || len(commands) == 0 {
		return fmt.Errorf("commands are not present")
	}

	txt := page.UserInput.GetText()
	if txt == "" {
		return nil
	}
	defer page.UserInput.SetText("")

	if strings.HasPrefix(txt, "/") {
		args := strings.Split(txt, " ")
		if cmd, ok := commands[args[0]]; ok {
			if err := cmd.Exec(ui, page, args); err != nil {
				return err
			}
		}
	}

	if id != OtterchatKey {
		fmt.Fprintf(page.ChatView, "%s\n", txt)
	}

	return nil
}

// AddPage adds a next context switcher between channels
func (ui *UI) AddPage(id string, commands map[string]Command) {
	ui.ConnectionList.AddItem(id, "", '-', nil)
	ui.ConnectionList.SetShortcutColor(tcell.ColorBlack).
		SetSelectedTextColor(tcell.ColorCornflowerBlue).
		SetSelectedBackgroundColor(tcell.ColorBlack)

	page := &Page{}

	page.ChatView = tview.NewTextView().
		SetScrollable(true).
		SetTextColor(tcell.ColorWhite).
		SetChangedFunc(func() {
			ui.App.Draw()
		})
	page.ChatView.SetBorder(true).
		SetBorderColor(tcell.ColorCornflowerBlue)

	page.UserList = tview.NewList()
	page.UserList.SetBorder(true).
		SetBorderColor(tcell.ColorBlack).
		SetTitle("Users")

	mflex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(page.ChatView, 0, 1, false).
		AddItem(page.UserList, 20, 1, false)

	page.ChannelInfoView = tview.NewTextView().
		SetScrollable(false).
		SetTextColor(tcell.ColorWhite)
	page.ChannelInfoView.SetBackgroundColor(tcell.ColorCornflowerBlue)

	page.InputInfoView = tview.NewTextView().
		SetScrollable(false).
		SetTextColor(tcell.ColorWhite)
	page.InputInfoView.SetBackgroundColor(tcell.ColorCornflowerBlue)

	page.UserInput = tview.NewInputField().
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabel("usr >>").
		SetPlaceholder("message").
		SetAcceptanceFunc(tview.InputFieldMaxLength(512))

	rflex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(page.ChannelInfoView, 1, 1, false).
		AddItem(mflex, 0, 2, false).
		AddItem(page.InputInfoView, 1, 1, false).
		AddItem(page.UserInput, 1, 1, true)

	ui.Pages.AddPage(id, rflex, true, true)
	ui.PageMap[id] = page

	fmt.Fprint(page.ChannelInfoView, id)
	fmt.Fprint(page.InputInfoView, id)

	page.UserInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if err := ui.handleInput(id, page, commands); err != nil {
				fmt.Fprintf(page.ChatView, "%s\n", err.Error())
			}
		}
	})

}

// Run the User Interface
func (ui *UI) Run() error {
	if err := ui.App.SetRoot(ui.AppView, true).Run(); err != nil {
		return err
	}
	return nil
}
