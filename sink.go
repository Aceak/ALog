package alog

type MultiSink struct {
	sinks []Sink
}

func NewMultiSink(sinks ...Sink) *MultiSink {
	return &MultiSink{sinks: sinks}
}

func (s *MultiSink) Write(line string) {
	for _, sink := range s.sinks {
		sink.Write(line)
	}
}

func (l *Logger) SetSink(sink Sink) {
	l.sink = sink
}
