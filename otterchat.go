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

package main

import (
	"fmt"
	"os"

	"github.com/chrisppy/otterchat/ui"
	"github.com/gdamore/tcell"
)

func main() {
	ui := ui.Init()

	ui.AddPage("otterchat")
	ui.Pages.SwitchToPage("otterchat")

	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt {
			if event.Key() == tcell.KeyUp {
				size := ui.ConnectionList.GetItemCount()
				currItemIndex := ui.ConnectionList.GetCurrentItem()

				newIndex := (currItemIndex - 1)
				if newIndex < 0 {
					newIndex = size - 1
				}

				id, _ := ui.ConnectionList.GetItemText(newIndex)
				ui.ConnectionList.SetCurrentItem(newIndex)

				ui.Pages.SwitchToPage(id)
			} else if event.Key() == tcell.KeyDown {
				size := ui.ConnectionList.GetItemCount()
				currItemIndex := ui.ConnectionList.GetCurrentItem()

				newIndex := (currItemIndex + 1)
				if newIndex >= size {
					newIndex = 0
				}

				id, _ := ui.ConnectionList.GetItemText(newIndex)
				ui.ConnectionList.SetCurrentItem(newIndex)
				ui.Pages.SwitchToPage(id)
			}
		}
		return event
	})

	if err := ui.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
