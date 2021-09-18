package color

import (
	"strings"

	"github.com/mgutz/ansi"
)

type Color struct {
	FgColor string
	FgAttr  []string
	BgColor string
	BgAttr  []string
}

func (c *Color) Load(color string) {
	c.reset()

	parts := strings.Split(color, "_")
	if len(parts) == 0 {
		return
	}
	flag := false
	for idx := 0; idx < len(parts); idx++ {
		part := parts[idx]
		if idx == 0 && part == "c" {
			continue
		}
		if !flag {
			if len(part) == 1 {
				c.FgAttr = append(c.FgAttr, part)
			} else if part == "light" {
				c.FgAttr = append(c.FgAttr, "h")
			} else if part == "dark" {
				c.FgAttr = append(c.FgAttr, "d")
			} else if strings.HasPrefix(part, "lt") {
				c.FgAttr = append(c.FgAttr, "h")
				c.FgColor = strings.TrimPrefix(part, "lt")
				flag = !flag
			} else if strings.HasPrefix(part, "dk") {
				c.FgAttr = append(c.FgAttr, "d")
				c.FgColor = strings.TrimPrefix(part, "dk")
				flag = !flag
			} else {
				c.FgColor = part
				flag = !flag
			}
		} else {
			if len(part) == 1 {
				c.BgAttr = append(c.BgAttr, part)
			} else if part == "light" {
				c.BgAttr = append(c.BgAttr, "h")
			} else if part == "dark" {
				c.BgAttr = append(c.BgAttr, "d")
			} else if strings.HasPrefix(part, "lt") {
				c.BgAttr = append(c.BgAttr, "h")
				c.BgColor = strings.TrimPrefix(part, "lt")
			} else if strings.HasPrefix(part, "dk") {
				c.BgAttr = append(c.BgAttr, "d")
				c.BgColor = strings.TrimPrefix(part, "dk")
			} else {
				c.BgColor = part
			}
		}
	}
}

func (c *Color) Colorized(str string) string {
	return ansi.Color(str, c.convert())
}

func (c *Color) convert() string {
	res := ""
	res += c.FgColor
	if len(c.FgAttr) > 0 {
		res += "+"
		for _, a := range c.FgAttr {
			res += a
		}
	}

	if c.BgColor == "" {
		return res
	}

	res += ":"
	res += c.BgColor
	if len(c.BgAttr) > 0 {
		res += "+"
		for _, a := range c.BgAttr {
			res += a
		}
	}
	return res
}

func (c *Color) reset() {
	c.BgAttr = make([]string, 0)
	c.FgAttr = make([]string, 0)
	c.BgColor = ""
	c.FgColor = ""
}
