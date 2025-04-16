/*
This file contains code adapted from the tomlv tool by Andrew Gallant (BurntSushi)
Original code: https://github.com/BurntSushi/toml/blob/master/cmd/tomlv/main.go

The MIT License (MIT)

Copyright (c) 2013 TOML authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type ConfigFile string

func (cf ConfigFile) Parse() (*Config, error) {
	_, err := os.Stat(string(cf))
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config file '%s' does not exist", cf)
	}

	config := &Config{
		Workspaces: make(map[string]*WorkspaceConfig),
		Apps:       make(map[string]*AppConfig),
	}

	_, err = toml.DecodeFile(string(cf), config)
	if err != nil {
		var perr toml.ParseError
		if errors.As(err, &perr) {
			return nil, fmt.Errorf("error in '%s': %s", cf, perr.ErrorWithPosition())
		}
		return nil, fmt.Errorf("error in '%s': %s", cf, err)
	}

	for name, workspace := range config.Workspaces {
		workspace.Name = name
	}

	for name, app := range config.Apps {
		app.Name = name
	}

	return config, nil
}
