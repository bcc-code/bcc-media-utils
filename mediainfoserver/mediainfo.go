package mediainfoserver

import (
	"fmt"
	"log"
	"strconv"
)

func (m MediaInfoResult) HasAudio() bool {
	for _, track := range m.Media.Track {
		if track.Type == "Audio" {
			return true
		}
	}
	return false
}

func (m MediaInfoResult) HasVideo() bool {
	for _, track := range m.Media.Track {
		if track.Type == "Video" {
			return true
		}
	}
	return false
}

func (m MediaInfoResult) HasStillImage() bool {
	for _, track := range m.Media.Track {
		if track.Type == "Image" {
			return true
		}
	}
	return false
}

func (m MediaInfoResult) MediaWidth() int {
	for _, track := range m.Media.Track {
		if track.Type == "Video" || track.Type == "Image" {
			w, err := strconv.Atoi(track.Width)
			if err != nil {
				log.Printf("Error: %s", err.Error())
				continue
			}

			return w
		}
	}
	return 0
}

func (m MediaInfoResult) MediaHeight() int {
	for _, track := range m.Media.Track {
		if track.Type == "Video" || track.Type == "Image" {
			h, err := strconv.Atoi(track.Height)
			if err != nil {
				log.Printf("Error: %s", err.Error())
				continue
			}

			return h
		}
	}
	return 0
}

func (m MediaInfoResult) VideoTCstart() string {
	for _, track := range m.Media.Track {
		if track.Type == "Other" &&
			track.SubType == "Time code" {
			return fmt.Sprintf("%s@%s", track.TimeCodeFirstFrame, track.FrameRate)
		}
	}
	return "00:00:00:00@25"
}

type AspectRatio string

const (
	AspectRatio43    AspectRatio = "4:3"
	AspectRatio169   AspectRatio = "16:9"
	AspectRatioOther AspectRatio = "other"
)

func (m MediaInfoResult) AspectRatio() AspectRatio {
	for _, track := range m.Media.Track {
		if track.Type == "Video" {
			switch track.DisplayAspectRatio {
			case "1.33":
				return AspectRatio43
			case "1.778":
				return AspectRatio169
			}
		}

	}
	return AspectRatioOther
}

func (m MediaInfoResult) AudioChannels() int {
	audioCount := 0
	for _, track := range m.Media.Track {
		if track.Type == "Audio" {
			count, err := strconv.Atoi(track.Channels)
			if err == nil {
				audioCount += count
			}
		}
	}
	return audioCount
}
