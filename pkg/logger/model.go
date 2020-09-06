package logger

type LoggerConfig struct {
	//LogLevel ... the default log level
	LogLevel string
	//ContextString ... in case of logging extra context information to logs, this should be enabled
	ContextString string
}
