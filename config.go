package alog

type LoggerConfig struct {
	Level  string
	Format string

	EnableConsole bool
	FilePath      string

	PanicBehavior string
	FatalBehavior string
}

func DefaultConfig() LoggerConfig {
	return LoggerConfig{
		Level:         "debug",
		Format:        "[{time}] [{level}] {msg}",
		EnableConsole: true,
		FilePath:      "",
		PanicBehavior: "panic",
		FatalBehavior: "exit",
	}
}
