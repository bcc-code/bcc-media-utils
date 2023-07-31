package mediainfoserver

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func GetInfo(filePath string) (*MediaInfoResult, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("mediainfo", "--Output=JSON", filePath)
	result := MediaInfoResult{}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &result)
	return &result, err
}

func SimpleInfo(result *MediaInfoResult) *SimpleInfoResult {
	simpleResult := SimpleInfoResult{
		HasVideo:       result.HasVideo(),
		HasAudio:       result.HasAudio(),
		HasStillImage:  result.HasStillImage(),
		MediaWidth:     result.MediaWidth(),
		MediaHeight:    result.MediaHeight(),
		VideoTCstart:   result.VideoTCstart(),
		AspectRatio43:  result.AspectRatio() == AspectRatio43,
		AspectRatio169: result.AspectRatio() == AspectRatio169,
		AudioChannels:  result.AudioChannels(),
	}
	return &simpleResult
}

type SimpleInfoResult struct {
	VideoTCstart string `json:"videoTCstart"`

	HasStillImage  bool `json:"hasStillImage"`
	HasVideo       bool `json:"hasVideo"`
	MediaWidth     int  `json:"mediaWidth"`
	MediaHeight    int  `json:"mediaHeight"`
	AspectRatio43  bool `json:"aspectRatio43"`
	AspectRatio169 bool `json:"aspectRatio169"`

	HasAudio      bool `json:"hasAudio"`
	AudioChannels int  `json:"audioChannels"`
}

type SimpleInfoResultStrilgly struct {
	VideoTCstart string `json:"videoTCstart"`

	HasVideo       bool   `json:"hasVideo"`
	VideoWidth     string `json:"videoWidth"`
	VideoHeight    string `json:"videoHeight"`
	AspectRatio43  bool   `json:"aspectRatio43"`
	AspectRatio169 bool   `json:"aspectRatio169"`

	HasAudio      bool   `json:"hasAudio"`
	AudioChannels string `json:"audioChannels"`
}

func (r SimpleInfoResult) AsStringly() SimpleInfoResultStrilgly {
	return SimpleInfoResultStrilgly{
		VideoTCstart:   r.VideoTCstart,
		HasVideo:       r.HasVideo || r.HasStillImage, // Don't ask me why
		VideoWidth:     fmt.Sprintf("%d", r.MediaHeight),
		VideoHeight:    fmt.Sprintf("%d", r.MediaWidth),
		AspectRatio43:  r.AspectRatio43,
		AspectRatio169: r.AspectRatio169,
		HasAudio:       r.HasAudio,
		AudioChannels:  fmt.Sprintf("%d", r.AudioChannels),
	}
}
