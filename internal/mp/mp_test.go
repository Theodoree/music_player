package mp

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	buf, _ := os.ReadFile("/Users/ted/workspace/go/music_player/save/周杰伦 温岚-屋顶-5257138.lrc")
	_, _ = decodeLrc(string(buf))
}
