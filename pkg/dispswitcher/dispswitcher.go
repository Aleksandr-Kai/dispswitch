package dispswitcher

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (
	outInfo = regexp.MustCompile(`([a-zA-Z]{3}[\w\-]+)( \w+)( \w+)?( \d+x\d+)?`)
	outMode = regexp.MustCompile(`\d+x\d+`)
)

type Output struct {
	Name      string
	Connected bool
	Primary   bool
	Enabled   bool
}

func (o *Output) IsEmpty() bool {
	return o.Name == ""
}

func (o *Output) Enable(options ...string) {
	options = append([]string{"--output", o.Name, "--auto"}, options...)
	cmd := exec.Command("xrandr", options...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(cmd.String())
		fmt.Println(string(out))
		fmt.Println(err)
	}
}

func (o *Output) Disable() {
	cmd := exec.Command("xrandr", "--output", o.Name, "--off")
	cmd.Run()
}

type OutputList struct {
	list []Output
}

func (ol *OutputList) Add(out Output) {
	ol.list = append(ol.list, out)
}

func (ol *OutputList) Scan() {
	cmd := exec.Command("xrandr", "-q")
	list, _ := cmd.Output()
	data := outInfo.FindAll(list, -1)
	for _, s := range data {
		ol.Add(ParseInfo(string(s)))
	}
}

func (ol *OutputList) Outputs() []Output {
	return ol.list
}

func (ol *OutputList) Connected() *OutputList {
	res := &OutputList{}
	for _, out := range ol.list {
		if out.Connected {
			res.Add(out)
		}
	}
	return res
}

func (ol *OutputList) Enabled() *OutputList {
	res := &OutputList{}
	for _, out := range ol.list {
		if out.Enabled {
			res.Add(out)
		}
	}
	return res
}

func (ol *OutputList) Count() int {
	return len(ol.list)
}

func (ol *OutputList) First() *Output {
	return &ol.list[0]
}

func (ol *OutputList) Primary() *Output {
	for _, out := range ol.list {
		if out.Primary {
			return &out
		}
	}
	return &Output{}
}

func (ol *OutputList) NotPrimary() *OutputList {
	res := &OutputList{}
	for _, out := range ol.list {
		if !out.Primary {
			res.Add(out)
		}
	}
	return res
}

func ParseInfo(str string) Output {
	out := Output{
		Name:      strings.Split(str, " ")[0],
		Connected: strings.Contains(str, " connected"),
		Primary:   strings.Contains(str, "primary"),
		Enabled:   outMode.MatchString(str),
	}
	return out
}
