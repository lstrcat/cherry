package cherryLogger

import (
	ctime "github.com/cherry-game/cherry/extend/time"
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	config := defaultConsoleConfig()
	config.EnableConsole = false
	config.EnableWriteFile = true
	config.FileLinkPath = "logs/log1.log"
	config.FilePathFormat = "logs/log1_%Y%m%d%H%M.log"

	log1 := NewConfigLogger(config)

	for i := 0; i < b.N; i++ {
		log1.Debug(ctime.Now().ToDateTimeFormat())
	}
}
