package db_logger

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func Red(str string) string    { return colorRed + str + colorReset }
func Yellow(str string) string { return colorYellow + str + colorReset }
func Blue(str string) string   { return colorBlue + str + colorReset }
