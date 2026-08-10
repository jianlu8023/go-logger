package colour

import (
	"os"

	"github.com/jianlu8023/go-logger/v2/internal/colour/colours"
	"github.com/mattn/go-isatty"
)

var c *colours.Color

func init() {
	c = colours.New()
	// 直接检测 os.Stdout 是否为终端，避免 SetOutput 中对 colorable.Wrapper 类型断言失败的问题
	// SetOutput 使用 colorable.NewColorableStdout()，其返回类型不是 *os.File，
	// 导致 SetOutput 内的 isatty 判断始终以非终端处理（disabled=true）
	// 因此这里需要独立基于 os.Stdout 做判断
	if isatty.IsTerminal(os.Stdout.Fd()) {
		c.Enable()
	} else {
		c.Disable()
	}
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
