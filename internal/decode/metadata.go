package decode

import (
	"encoding/base64"
	"errors"
	"github.com/dhowden/tag"
	"github.com/go-audio/wav"
	"io"
)

// support MP3 M4A M4B M4P ALAC FLAC OGG DSF
func getTagMetadata(r io.ReadSeeker) (Metadata, error) {
	decoder, err := tag.ReadFrom(r)
	if err != nil {
		return Metadata{}, err
	}
	
	var m Metadata
	m.Title = decoder.Title()
	m.Artist = decoder.Artist()
	m.Album = decoder.Album()
	pic := decoder.Picture()
	if pic != nil {
		m.AlbumPicture = base64.StdEncoding.EncodeToString(pic.Data)
		m.MimeType = pic.MIMEType
	}
	return m, nil
}

func getFLACMetadata(r io.ReadSeeker) (Metadata, error) {
	return getTagMetadata(r)
}

func getMP3Metadata(r io.ReadSeeker) (Metadata, error) {
	return getTagMetadata(r)
}

func getWAVMetadata(r io.ReadSeeker) (Metadata, error) {
	decoder := wav.NewDecoder(r)
	if decoder == nil {
		return Metadata{}, errors.New("wav decoder is nil")
	}
	
	decoder.ReadMetadata()
	var m Metadata
	m.Title = decoder.Metadata.Title
	m.Artist = decoder.Metadata.Artist
	m.Album = decoder.Metadata.Product
	return m, nil
}
