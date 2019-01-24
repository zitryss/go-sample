package log

import (
	"strings"
	"testing"
)

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level int
		want  string
	}{
		{
			name:  "",
			level: debug,
			want:  "debug: message1\ninfo: message2\nwarn: message3\nerror: message4\ncritical: message5\n",
		},
		{
			name:  "",
			level: info,
			want:  "info: message2\nwarn: message3\nerror: message4\ncritical: message5\n",
		},
		{
			name:  "",
			level: warn,
			want:  "warn: message3\nerror: message4\ncritical: message5\n",
		},
		{
			name:  "",
			level: error_,
			want:  "error: message4\ncritical: message5\n",
		},
		{
			name:  "",
			level: critical,
			want:  "critical: message5\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			l.SetOutput(w)
			l.SetFlags(0)
			l.level = tt.level
			Debug("message1")
			Info("message2")
			Warn("message3")
			Error("message4")
			Critical("message5")
			got := w.String()
			if got != tt.want {
				t.Errorf("%v = %v; want %v", tt.level, tt.want, got)
			}
		})
	}
}
