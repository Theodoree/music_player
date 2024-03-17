package mp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"k8s.io/klog"
	"os"
	"regexp"
	"strings"
	"time"
)

// LrcLine 结构体表示一行LRC歌词信息
type LrcLine struct {
	Time   time.Time
	Lyrics string
}

func parseLrcLine(line string) (LrcLine, error) {
	re := regexp.MustCompile(`^\[(\d{2}:\d{2}\.\d{2})\](.*)$`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return LrcLine{}, errors.New("")
	}
	
	timeTag := matches[1]
	lyrics := matches[2]
	ts, _ := time.Parse("04:05.999999999", timeTag)
	return LrcLine{Time: ts, Lyrics: lyrics}, nil
}

func decodeLrc(str string) ([]LrcLine, error) {
	var lrcLines []LrcLine
	r := bufio.NewReader(strings.NewReader(str))
	for {
		buf, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			klog.Error(err)
			return nil, err
		}
		lrcLine, err := parseLrcLine(string(buf))
		if err == nil {
			lrcLines = append(lrcLines, lrcLine)
		} else {
			fmt.Fprintf(os.Stderr, "Error parsing line: %v\n", err)
		}
	}
	return lrcLines, nil
}
