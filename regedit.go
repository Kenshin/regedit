/*
 Regedit - Windows Regedit CLI by Golang. Vesion 0.0.1

 Website https://github.com/kenshin/regedit

 Copyright (c) 2016 Kenshin Wang <kenshin@ksria.com>
*/

package regedit

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	Add = iota
	Query
	Delete
	Copy
	Save
	Restore
	Load
	Unload
	Compare
	Export
	Import
)

const (
	HKCR = iota
	HKLM
	HKU
	HKCU
	HKCC
)

const (
	SZ = iota
	BINARY
	MULTI
	EXPAND
)

type (
	Regedit struct {
		Action string
		Field  string
		Reg
	}
	Reg struct {
		Key   string
		Type  string
		Value string
	}
	RegCmd struct {
		cmd *exec.Cmd
		reg *Regedit
	}
)

var (
	Actions = map[int]string{
		0:  "add",
		1:  "query",
		2:  "delete",
		3:  "copy",
		4:  "save",
		5:  "restore",
		6:  "load",
		7:  "unload",
		8:  "compare",
		9:  "export",
		10: "import",
	}
	Fields = map[int]string{
		0: "HKEY_CLASSES_ROOT",
		1: "HKEY_LOCAL_MACHINE",
		2: "HKEY_USERS",
		3: "HKEY_CURRENT_USER",
		4: "HKEY_CURRENT_CONFIG",
	}
	Types = map[int]string{
		0: "REG_SZ",
		1: "REG_BINARY",
		2: "REG_MULTI_SZ",
		3: "REG_EXPAND_SZ",
	}
)

/*
  Create Regedit struct
*/
func New(action, filed int, path string) *Regedit {
	return &Regedit{Action: Actions[action], Field: Fields[filed] + path}
}

/*
  Add
*/
func (this *Regedit) Add(reg Reg) RegCmd {
	(*this).Reg = reg
	return RegCmd{exec.Command("cmd", "/c", "reg", this.Action, this.Field, "/v", this.Key, "/t", this.Type, "/d", this.Value), this}
}

/*
  Search
*/
func (this *Regedit) Search(reg Reg) RegCmd {
	(*this).Reg = reg
	return RegCmd{exec.Command("cmd", "/c", "reg", this.Action, this.Field, "/s"), this}
}

/*
  Execute Regedit

  Support:
	- Add
	- Search

  Return:
  	- []Reg
  	- error
*/
func (this RegCmd) Exec() ([]Reg, error) {
	if this.reg.Action == Actions[Add] {
		this.cmd.Stdout = os.Stdout
		this.cmd.Stderr = os.Stderr
		this.cmd.Stdin = os.Stdin
		if err := this.cmd.Run(); err != nil {
			return nil, err
		}
	} else if this.reg.Action == Actions[Query] {
		if out, err := this.cmd.Output(); err != nil {
			return nil, err
		} else {
			buff, regList := bytes.NewBuffer(out), make([]Reg, 0)
			for {
				content, err := buff.ReadString('\n')
				content = strings.TrimSpace(content)
				if err != nil || err == io.EOF {
					return regList, nil
				}
				if arr := strings.Fields(content); len(arr) == 3 {
					if this.reg.Key == "" {
						regList = append(regList, Reg{arr[0], arr[1], arr[2]})
					} else if this.reg.Key != "" && strings.ToLower(arr[0]) == strings.ToLower(this.reg.Key) {
						regList = append(regList, Reg{arr[0], arr[1], arr[2]})
					}
				}
			}
		}
	} else {
		return nil, errors.New("Regedit only support Add and Search.")
	}
	return nil, nil
}
