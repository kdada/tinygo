package log

func init() {
	Register("console", ConsoleLoggerCreator)
	Register("file", FileLoggerCreator)
}
