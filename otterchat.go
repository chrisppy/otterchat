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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"

	"github.com/chrisppy/otterchat/api"
	"github.com/gdamore/tcell"
)

const pluginPath = "~/.otterchat/plugins"

func printSplashScreen(w io.Writer) {
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
	fmt.Fprint(w, asciiOtter)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func loadPlugins(w io.Writer) (map[string]api.Command, error) {
	if !exists(pluginPath) {
		return nil, fmt.Errorf("path: '%s' does not exist", pluginPath)
	}

	files, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		return nil, err
	}

	commands := make(map[string]api.Command)

f:
	for _, file := range files {
		name := file.Name()
		if filepath.Ext(name) != ".so" {
			continue
		}

		plug, err := plugin.Open(filepath.Join(pluginPath, name))
		if err != nil {
			fmt.Fprintf(w, "failed to open plugin '%s': %v\n", name, err)
			continue
		}

		cmdSymbol, err := plug.Lookup(api.CMDSymbolName)
		if err != nil {
			fmt.Fprintf(w, "plugin %s does not export symbol '%s'\n", name, api.CMDSymbolName)
			continue
		}

		cmds, ok := cmdSymbol.(api.Commands)
		if !ok {
			fmt.Fprintf(w, "Symbol %s (from %s) does not implement Commands interface\n", api.CMDSymbolName, name)
			continue
		}

		for cname, cmd := range cmds.Registry() {
			if _, ok := commands[cname]; ok {
				fmt.Fprintf(w, "connot load plugin: '%s' due to command: '%s' already present\n", name, cname)
				continue f
			}

			commands[cname] = cmd
		}
	}

	if len(commands) == 0 {
		return nil, fmt.Errorf("at least one command must be present")
	}

	return commands, nil
}

func main() {
	ui := api.Init()

	buf := &bytes.Buffer{}

	// TODO: load all plugins
	commands, err := loadPlugins(buf)
	if err != nil {
		fmt.Fprintf(buf, "error: %s\n", err.Error())
	}

	ui.AddPage(api.OtterchatKey, commands)
	ui.Pages.SwitchToPage(api.OtterchatKey)
	page := ui.PageMap[api.OtterchatKey]
	printSplashScreen(page.ChatView)

	if buf != nil {
		fmt.Fprintf(page.ChatView, "%s\n", buf)
	}

	// TODO: check to see if there are any duplicated commands, if so return to error and exit

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
