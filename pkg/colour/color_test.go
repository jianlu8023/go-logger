package colour

import (
	"testing"

	"github.com/labstack/gommon/color"
)

func TestBlue(t *testing.T) {
	str := "==> 执行语句: SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE 'basic%' ORDER BY SCHEMA_NAME='basic' DESC,SCHEMA_NAME limit 1 \n==> 影响行数: 1 \n==> 执行时间: 0.975ms"
	blue := Blue(str)
	color.Println(blue)

}
