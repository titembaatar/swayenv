package sway

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type SwayCmd struct {
	Command string
	Type    string
	Raw     bool
}

func NewSwayCmd(command string) *SwayCmd {
	return &SwayCmd{
		Command: command,
		Type:    "command",
		Raw:     false,
	}
}

func NewSwayCmdType(msgType string) *SwayCmd {
	return &SwayCmd{
		Command: "",
		Type:    msgType,
		Raw:     true,
	}
}

func (sc *SwayCmd) Output() ([]byte, error) {
	var args []string

	if sc.Raw {
		args = append(args, "-r")
	}

	args = append(args, "-t", sc.Type)

	if sc.Type == "command" && sc.Command != "" {
		args = append(args, sc.Command)
	}

	cmd := exec.Command("swaymsg", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command '%s' failed: %w", cmd, err)
	}

	return output, nil
}

func (sc *SwayCmd) Run() error {
	_, err := sc.Output()
	return err
}

func (sc *SwayCmd) GetJSON(v any) error {
	if !sc.Raw {
		sc.Raw = true
	}

	output, err := sc.Output()
	if err != nil {
		return err
	}

	return json.Unmarshal(output, v)
}
