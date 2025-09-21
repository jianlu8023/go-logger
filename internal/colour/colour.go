package colour

import (
	"github.com/labstack/gommon/color"
	// "github.com/fatih/color"
)

var c *color.Color

func init() {
	c = color.New()
	c.Enable()
}
func Magenta(str string) string   { return c.Magenta(str) }
func Blue(str string) string      { return c.Blue(str) }
func Yellow(str string) string    { return c.Yellow(str) }
func Red(str string) string       { return c.Red(str) }
func Black(str string) string     { return c.Black(str) }
func Green(str string) string     { return c.Green(str) }
func Cyan(str string) string      { return c.Cyan(str) }
func White(str string) string     { return c.White(str) }
func Grey(str string) string      { return c.Grey(str) }
func BlackBg(str string) string   { return c.BlackBg(str) }
func RedBg(str string) string     { return c.RedBg(str) }
func GreenBg(str string) string   { return c.GreenBg(str) }
func YellowBg(str string) string  { return c.YellowBg(str) }
func BlueBg(str string) string    { return c.BlueBg(str) }
func MagentaBg(str string) string { return c.MagentaBg(str) }
func CyanBg(str string) string    { return c.CyanBg(str) }
func WhiteBg(str string) string   { return c.WhiteBg(str) }
func Reset(str string) string     { return c.Reset(str) }
func Bold(str string) string      { return c.Bold(str) }
func Dim(str string) string       { return c.Dim(str) }
func Italic(str string) string    { return c.Italic(str) }
func Underline(str string) string { return c.Underline(str) }
func Inverse(str string) string   { return c.Inverse(str) }
func Hidden(str string) string    { return c.Hidden(str) }
func Strikeout(str string) string { return c.Strikeout(str) }
