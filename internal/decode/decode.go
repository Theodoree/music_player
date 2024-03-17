package decode

import (
	"context"
	"io"
	
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/music"
)

func newMP3Decoder(ctx context.Context, reader io.ReadSeekCloser, volume float64, cb *music.Callback) (Decoder, error) {
	return newBeepDecoder(ctx, reader, volume, cb, model.MusicTypeMP3)
}

func newFLACDecoder(ctx context.Context, reader io.ReadSeekCloser, volume float64, cb *music.Callback) (Decoder, error) {
	return newBeepDecoder(ctx, reader, volume, cb, model.MusicTypeFLAC)
}
func newWAVDecoder(ctx context.Context, reader io.ReadSeekCloser, volume float64, cb *music.Callback) (Decoder, error) {
	return newBeepDecoder(ctx, reader, volume, cb, model.MusicTypeWAV)
}
