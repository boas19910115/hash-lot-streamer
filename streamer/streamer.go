package streamer

import (
	"bot/helpers"
	"bot/processor"
	"bot/ytdl"
	"errors"
	"fmt"
	"net/url"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Item struct {
	Input string
	User  *gotgbot.User
	Video *ytdl.Video
}

var (
	now             Item
	errSignalKilled = errors.New("signal: killed")
)

func Stream(b *gotgbot.Bot, input string, user *gotgbot.User) error {
	var video *ytdl.Video
	origInput := input
	UrlList := helpers.GetVideos()
	var VideoPathList []string
	var err error
	for _, URL := range UrlList {
		_, err := url.ParseRequestURI(URL)
		if err == nil {
			video, err = ytdl.Download(URL)
			if err == nil {
				VideoPathList = append(VideoPathList, video.Url)
			}
		}
	}
	errc := make(chan error)
	go func() {
		for {
			err, ok := <-errc
			if !ok {
				break
			}
			if errSignalKilled.Error() == err.Error() {
				continue
			}
			b.SendMessage(user.Id, fmt.Sprintf("Failed to process: %s", err.Error()), nil)
		}
	}()
	err = processor.Process(VideoPathList, errc)
	if err == nil {
		now.Input = origInput
		now.User = user
		now.Video = video
	} else {
		close(errc)
	}
	return err
}

func Now() (bool, Item) {
	return processor.Processing(), now
}
