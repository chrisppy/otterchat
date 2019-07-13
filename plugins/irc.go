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
	"io"
	"strings"

	"github.com/chrisppy/otterchat/api"
)

type serverCMD string

func (c serverCMD) Name() string {
	return string(c)
}

func (c serverCMD) Usage() string {
	return fmt.Sprintf("/%s [add/delete] [name] [url]", c.Name())
}

func (c serverCMD) Desc() string {
	return "Add or Remove an IRC server"
}

func addServer(i *api.UI, name, url string) {
	fmt.Fprintf(i.PageMap["otterchat"].ChatView, "Adding IRC service named:'%s' with url:'%s'\n", name, url)
}

func deleteServer(i *api.UI, name string) {
	fmt.Fprintf(i.PageMap["otterchat"].ChatView, "Deleting IRC service named:'%s'\n", name)
}

func (c serverCMD) Exec(i *api.UI, page *api.Page, args []string) error {
	len := len(args)
	if len < 3 {
		return fmt.Errorf(c.Usage())
	}
	action := args[1]
	name := args[2]
	switch action {
	case "add":
		if len != 4 {
			return fmt.Errorf(c.Usage())
		}
		addServer(i, name, args[3])
		return nil
	case "delete":
		if len != 3 {
			return fmt.Errorf(c.Usage())
		}
		deleteServer(i, name)
		return nil
	}
	return nil
}

type joinCMD string

func (c joinCMD) Name() string {
	return string(c)
}

func (c joinCMD) Usage() string {
	return fmt.Sprintf("/%s [name] [channel]", c.Name())
}

func (c joinCMD) Desc() string {
	return "Join an IRC channel"
}

func joinChannel(i *api.UI, name, channel string) {
	i.AddPage(channel, Commands.Registry())

	i.ConnectionList.SetCurrentItem(i.ConnectionList.GetItemCount() - 1)

	fmt.Fprintf(i.PageMap["otterchat"].ChatView, "joining channel:'%s' on the:'%s' server\n", channel, name)
}

func (c joinCMD) Exec(i *api.UI, page *api.Page, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("must be in the form: /join [name] [channel]")
	}

	name := args[1]
	channel := args[2]

	if !strings.HasPrefix(channel, "#") {
		return fmt.Errorf("channel must start with #")
	}

	joinChannel(i, name, channel)

	return nil
}

// ircCommands represents a collection of commands supported by this
// command module.
type ircCommands struct {
	stdout io.Writer
}

func (c *ircCommands) Registry() map[string]api.Command {
	return map[string]api.Command{
		"/server": serverCMD("/server"),
		"/join":   joinCMD("/join"),
	}
}

// Commands plugin entry point
var Commands ircCommands
