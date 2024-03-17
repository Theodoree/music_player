package decode

import (
	"context"
	"io"
	"sync"
	"time"
	
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/music"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"k8s.io/klog"
)

type beepDecoder struct {
	// ctx is the context of the decoder.
	ctx context.Context
	// cancel is the cancel function of the decoder.
	cancel context.CancelFunc
	// r is the reader of the decoder.
	r io.ReadSeekCloser
	// streamer is the streamer of the decoder.
	streamer beep.StreamSeekCloser
	// format is the format of the decoder.
	format beep.Format
	// ctrl is the controller of the decoder.
	ctrl *beep.Ctrl
	// volume is the volume of the decoder.
	volume *effects.Volume
	// cb is when action is done, call this callback.
	cb        *music.Callback
	musicType model.MusicType
	// stats is the status of the decoder.
	status    model.StatusController
	playOnce  sync.Once
	closeOnce sync.Once
}

func newBeepDecoder(ctx context.Context, reader io.ReadSeekCloser, volume float64, cb *music.Callback, Type model.MusicType) (Decoder, error) {
	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
		err      error
		decoder  beepDecoder
	)
	switch Type {
	case model.MusicTypeMP3:
		streamer, format, err = mp3.Decode(reader)
	case model.MusicTypeFLAC:
		streamer, format, err = flac.Decode(reader)
	case model.MusicTypeWAV:
		streamer, format, err = wav.Decode(reader)
	default:
		panic("unhandled default case")
	}
	if err != nil {
		return nil, err
	}
	decoder.ctx, decoder.cancel = context.WithCancel(ctx)
	
	decoder.r = reader
	decoder.streamer = streamer
	decoder.format = format
	decoder.ctrl = &beep.Ctrl{Streamer: beep.Seq(streamer, beep.Callback(func() {
		decoder.Close()
	})), Paused: false}
	decoder.volume = &effects.Volume{
		Streamer: decoder.ctrl,
		Base:     2,
		Silent:   false,
	}
	decoder.cb = cb
	decoder.musicType = Type
	decoder.SetVolume(volume)
	
	return &decoder, nil
}

func (d *beepDecoder) Play() {
	speaker.Lock()
	defer speaker.Unlock()
	d.ctrl.Paused = false
	d.playOnce.Do(func() {
		go d.Runner()
	})
	return
}

func (d *beepDecoder) Pause() {
	speaker.Lock()
	defer speaker.Unlock()
	d.ctrl.Paused = true
	return
}

func (d *beepDecoder) Stop() {
	d.status.SetStop()
	d.Close()
}

var _volume = [101]float64{-6.60, -6.60, -5.64, -5.05, -4.64, -4.32, -4.05, -3.83, -3.64, -3.47, -3.32, -3.18, -3.05, -2.94, -2.83, -2.73, -2.64, -2.55, -2.47, -2.39, -2.32, -2.25, -2.18, -2.12, -2.05, -2.00, -1.94, -1.88, -1.83, -1.78, -1.73, -1.68, -1.64, -1.59, -1.55, -1.51, -1.47, -1.43, -1.39, -1.35, -1.32, -1.28, -1.25, -1.21, -1.18, -1.15, -1.12, -1.08, -1.05, -1.02, -1.00, -0.97, -0.94, -0.91, -0.88, -0.86, -0.83, -0.81, -0.78, -0.76, -0.73, -0.71, -0.68, -0.66, -0.64, -0.62, -0.59, -0.57, -0.55, -0.53, -0.51, -0.49, -0.47, -0.45, -0.43, -0.41, -0.39, -0.37, -0.35, -0.34, -0.32, -0.30, -0.28, -0.26, -0.25, -0.23, -0.21, -0.20, -0.18, -0.16, -0.15, -0.13, -0.12, -0.10, -0.08, -0.07, -0.05, -0.04, -0.02, -0.01, -0.00}

func (d *beepDecoder) SetVolume(f float64) {
	speaker.Lock()
	defer speaker.Unlock()
	if f > 1 {
		f = 1
	}
	
	d.volume.Volume = _volume[int(f*100)]
}

func (d *beepDecoder) Seek(f float64) {
	if f > 100 {
		f = 100
	}
	total := d.streamer.Len()
	seek := int(float64(total) * f / 100)
	speaker.Lock()
	defer speaker.Unlock()
	if err := d.streamer.Seek(seek); err != nil {
		klog.Error(err)
		return
	}
}

func (d *beepDecoder) Metadata() (Metadata, error) {
	var err error
	var m Metadata
	switch d.musicType {
	case model.MusicTypeMP3:
		_, _ = d.r.Seek(0, 0)
		m, err = getMP3Metadata(d.r)
	case model.MusicTypeFLAC:
		_, _ = d.r.Seek(0, 0)
		m, err = getFLACMetadata(d.r)
	case model.MusicTypeWAV:
		_, _ = d.r.Seek(0, 0)
		m, err = getWAVMetadata(d.r)
	}
	if err != nil {
		return Metadata{}, err
	}
	
	m.EndTime = d.format.SampleRate.D(d.streamer.Len()).Round(time.Second)
	return m, nil
}
func (d *beepDecoder) CurTime() time.Duration {
	return d.format.SampleRate.D(d.streamer.Position()).Round(time.Second)
}
func (d *beepDecoder) EndTime() time.Duration {
	return d.format.SampleRate.D(d.streamer.Len()).Round(time.Second)
}

func (d *beepDecoder) Runner() {
	speaker.Init(d.format.SampleRate, d.format.SampleRate.N(time.Second/10))
	speaker.Play(d.volume)
	defer d.Close()
	defer speaker.Clear()
	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			if d.cb != nil {
				d.cb.CurTime(d.format.SampleRate.D(d.streamer.Position()).Round(time.Second))
			}
		}
	}
}

func (d *beepDecoder) Close() {
	d.closeOnce.Do(func() {
		d.cancel()
		if d.streamer != nil {
			_ = d.streamer.Close()
		}
		if d.r != nil {
			_ = d.r.Close()
		}
		d.status.SetPlayDone()
		if d.cb != nil {
			go d.cb.DoneFn(d.status.Load())
		}
	})
}
