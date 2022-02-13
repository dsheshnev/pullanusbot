package usecases

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/ailinykh/pullanusbot/v2/core"
)

func CreateYoutubeFlow(l core.ILogger, mediaFactory core.IMediaFactory, videoFactory core.IVideoFactory, videoSplitter core.IVideoSplitter) *YoutubeFlow {
	return &YoutubeFlow{l: l, mediaFactory: mediaFactory, videoFactory: videoFactory, videoSplitter: videoSplitter}
}

type YoutubeFlow struct {
	mutex         sync.Mutex
	l             core.ILogger
	mediaFactory  core.IMediaFactory
	videoFactory  core.IVideoFactory
	videoSplitter core.IVideoSplitter
}

// HandleText is a core.ITextHandler protocol implementation
func (flow *YoutubeFlow) HandleText(message *core.Message, bot core.IBot) error {
	r := regexp.MustCompile(`youtu\.?be(\.com)?\/(watch\?v=)?([\w\-_]+)`)
	match := r.FindStringSubmatch(message.Text)
	if len(match) == 4 {
		err := flow.process(match[3], message, bot)
		if err != nil {
			return err
		}

		if !strings.Contains(message.Text, " ") {
			return bot.Delete(message)
		}
	} else if strings.Contains(message.Text, "youtu") {
		for i, m := range match {
			flow.l.Info(i, " ", m)
		}
		return errors.New("possibble regexp mismatch: " + message.Text)
	}
	return nil
}

func (flow *YoutubeFlow) process(id string, message *core.Message, bot core.IBot) error {
	flow.mutex.Lock()
	defer flow.mutex.Unlock()

	flow.l.Infof("processing %s", id)
	media, err := flow.mediaFactory.CreateMedia(id)
	if err != nil {
		flow.l.Error(err)
		return err
	}

	if !message.IsPrivate && media[0].Duration > 900 {
		flow.l.Infof("skip video in group chat due to duration %d", media[0].Duration)
		return errors.New("skip video in group chat due to duration")
	}

	title := media[0].Title
	flow.l.Infof("downloading %s", id)
	file, err := flow.videoFactory.CreateVideo(id)
	if err != nil {
		return err
	}
	defer file.Dispose()

	caption := fmt.Sprintf(`<a href="https://youtu.be/%s">🎞</a> <b>%s</b> <i>(by %s)</i>`, id, title, message.Sender.DisplayName())
	_, err = bot.SendVideo(file, caption)
	if err != nil {
		flow.l.Error("Can't send video: ", err)
		if err.Error() == "telegram: Request Entity Too Large (400)" {
			flow.l.Info("Fallback to splitting")
			files, err := flow.videoSplitter.Split(file, 50000000)
			if err != nil {
				return err
			}

			for _, file := range files {
				defer file.Dispose()
			}

			for i, file := range files {
				caption := fmt.Sprintf(`<a href="https://youtu.be/%s">🎞</a> <b>[%d/%d] %s</b> <i>(by %s)</i>`, id, i+1, len(files), title, message.Sender.DisplayName())
				_, err := bot.SendVideo(file, caption)
				if err != nil {
					return err
				}
			}

			flow.l.Info("All parts successfully sent")
			return nil
		}
		return err
	}
	return nil
}
