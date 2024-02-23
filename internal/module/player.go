package module

import (
	"context"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"k8s.io/klog"
	"os"
	"sync"
	"time"
)

type player struct {
	sync.Once
	ctx             context.Context
	cancel          context.CancelFunc
	file            *os.File
	streamer        beep.StreamSeekCloser
	format          beep.Format
	ctrl            *beep.Ctrl      // 控制器
	volume          *effects.Volume // 带音量的控制
	processCallback func(cur time.Duration)
	doneCallBack    func()
}

func newPlayer(ctx context.Context, music model.Music, volume float64, processCb func(cur time.Duration), doneCb func()) (*player, error) {
	var (
		p        player
		streamer beep.StreamSeekCloser
		format   beep.Format
		err      error
	)
	p.ctx, p.cancel = context.WithCancel(ctx)
	
	p.file, err = os.Open(music.Path)
	if err != nil {
		p.Close()
		return nil, err
	}
	switch music.Type {
	case model.MusicTypeMP3:
		streamer, format, err = mp3.Decode(p.file)
	case model.MusicTypeFLAC:
		streamer, format, err = flac.Decode(p.file)
	case model.MusicTypeWAV:
		streamer, format, err = wav.Decode(p.file)
	default:
		panic("unhandled default case")
	}
	if err != nil {
		return nil, err
	}
	p.streamer = streamer
	p.format = format
	p.ctrl = &beep.Ctrl{Streamer: beep.Seq(p.streamer), Paused: false}
	p.volume = &effects.Volume{
		Streamer: p.ctrl,
		Base:     2,
		Volume:   volume,
		Silent:   false,
	}
	
	p.processCallback = processCb
	p.doneCallBack = doneCb
	go p.Runner()
	return &p, nil
}

func (p *player) Play() {
	speaker.Lock()
	p.ctrl.Paused = false
	speaker.Unlock()
}
func (p *player) Pause() {
	speaker.Lock()
	p.ctrl.Paused = true
	speaker.Unlock()
}
func (p *player) Stop() {
	p.Close()
}
func (p *player) SetVolume(volume float64) {
	speaker.Lock()
	p.volume.Volume = volume
	speaker.Unlock()
}
func (p *player) Runner() {
	speaker.Init(p.format.SampleRate, p.format.SampleRate.N(time.Second/10))
	speaker.Play(p.volume)
	defer speaker.Clear()
	defer p.Close()
	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			speaker.Lock()
			p.processCallback(p.format.SampleRate.D(p.streamer.Position()).Round(time.Second))
			speaker.Unlock()
		}
	}
	
}
func (p *player) EndTime() time.Duration {
	return p.format.SampleRate.D(p.streamer.Len()).Round(time.Second)
}
func (p *player) Seek(ptr float64) {
	speaker.Lock()
	defer speaker.Unlock()
	seek := float64(p.streamer.Len()) * (ptr / 100)
	err := p.streamer.Seek(int(seek))
	if err != nil {
		klog.Error(err)
	}
	
}
func (p *player) Close() {
	p.Once.Do(func() {
		p.cancel()
		if p.streamer != nil {
			_ = p.streamer.Close()
		}
		if p.file != nil {
			_ = p.file.Close()
		}
		if p.doneCallBack != nil && p.streamer.Position() == p.streamer.Len() {
			p.doneCallBack()
		}
	})
}
