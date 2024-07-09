package go_logger

type Config struct {
	// Mode file or date
	// file use lumberjack
	// date use rotatelogs
	Mode string `json:"mode"`
}
