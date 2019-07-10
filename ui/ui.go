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

package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const asciiOtter = `
                                                                               
                 ::/oooosoosoo::-                                              
             .+hhhysssoooooooossshhhy+.                                        
           +ddo//////////////////////sdh/                                      
     .yyyhmh///////////////////////////+dhhyho.      @-------------------@     
     M+//++///////o+////////////s+///////o///sM     @                     @    
     ymo////////oNhMd/////////sNdMh/////////yN+     |       |  |  *       |    
      -Ms///////oMMMm/////////sMMMd////////yN.      |       |--|  |       |    
      yd/////////+so//+yhhhhs+/+so//////////mo      |       |  |  |       |    
      Ms//////+o/////yMMMMMMMM//////o///////hN      @                     @    
      my/////sds//////ydNMMmh+//////smo/////hy       @-------.  .--------@     
      :No+yys/Nss+/+yhmmNNNmmmhs+/s+ym+sys+oN.               | |               
       ymd//sysddyhdNdyyyyyyyydMhhsddoy+//mN/               / /                
     /o..ymdo/sy////m/--------+h////yo/ydmo -/.            //                  
   .:.   ++odhd/////+ho------sh//////dhd+o:   :+          .                    
        o.  .yohdy+///oyysysy+///+ydhoy   :+    .                              
       o.   y    /ohhhhyssssyhhhho:   .y   -/                                  
       .   -:         .::::::.         o                                       
           :                            .                                      
                                                                               
`

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
	UIPages        map[string]*Page
}

// Init the User Interface
func Init() *UI {
	ui := new(UI)

	ui.UIPages = make(map[string]*Page)

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

// AddPage adds a next context switcher between channels
func (ui *UI) AddPage(id string) {
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
	ui.UIPages[id] = page

	if id == "otterchat" {
		fmt.Fprint(page.ChatView, asciiOtter)
	}
	fmt.Fprint(page.ChannelInfoView, id)
	fmt.Fprint(page.InputInfoView, id)

	page.UserInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if err := ui.handleInput(page, id); err != nil {
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

func (ui *UI) addServer(name, url string) {
	fmt.Fprintf(ui.UIPages["otterchat"].ChatView, "Adding IRC service named:'%s' with url:'%s'\n", name, url)
}

func (ui *UI) deleteServer(name string) {
	fmt.Fprintf(ui.UIPages["otterchat"].ChatView, "Deleting IRC service named:'%s'\n", name)
}

func (ui *UI) handleServerCMD(strs []string) error {
	len := len(strs)
	if len < 3 {
		return fmt.Errorf("must be in the form: /server [action] [name] [other params]")
	}
	action := strs[1]
	name := strs[2]

	switch action {
	case "add":
		if len != 4 {
			return fmt.Errorf("must be in the form: /server add [name] [url]")
		}
		ui.addServer(name, strs[3])
		return nil
	case "delete":
		if len != 3 {
			return fmt.Errorf("must be in the form: /server delete [name]")
		}
		ui.deleteServer(name)
		return nil
	}
	return nil
}

func (ui *UI) addChannel(name, channel string) {

	ui.AddPage(channel)

	ui.ConnectionList.SetCurrentItem(ui.ConnectionList.GetItemCount() - 1)

	fmt.Fprintf(ui.UIPages["otterchat"].ChatView, "joining channel:'%s' on the:'%s' server\n", channel, name)
}

func (ui *UI) handleJoinCMD(strs []string) error {
	if len(strs) != 3 {
		return fmt.Errorf("must be in the form: /join [name] [channel]")
	}

	name := strs[1]
	channel := strs[2]

	if !strings.HasPrefix(channel, "#") {
		return fmt.Errorf("channel must start with #")
	}

	ui.addChannel(name, channel)

	return nil
}

func (ui *UI) handleInput(page *Page, id string) error {
	txt := page.UserInput.GetText()
	if txt == "" {
		return nil
	}
	defer page.UserInput.SetText("")

	if strings.HasPrefix(txt, "/") {
		strs := strings.Split(txt, " ")
		switch strs[0] {
		case "/server":
			return ui.handleServerCMD(strs)
		case "/join":
			return ui.handleJoinCMD(strs)
		case "/connect":
		}
	}

	if id != "otterchat" {
		fmt.Fprintf(page.ChatView, "%s\n", txt)
	}

	return nil
}
