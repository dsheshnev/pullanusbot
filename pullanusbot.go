package main

import (
	"os"
	"path"

	"github.com/ailinykh/pullanusbot/v2/api"
	"github.com/ailinykh/pullanusbot/v2/core"
	"github.com/ailinykh/pullanusbot/v2/helpers"
	"github.com/ailinykh/pullanusbot/v2/infrastructure"
	"github.com/ailinykh/pullanusbot/v2/usecases"
	"github.com/google/logger"
)

func main() {
	logger, close := createLogger()
	defer close()

	telebot := api.CreateTelebot(os.Getenv("BOT_TOKEN"), logger)
	telebot.SetupInfo()

	localizer := infrastructure.GameLocalizer{}
	dbFile := path.Join(getWorkingDir(), "pullanusbot.db")
	settingsStorage := infrastructure.CreateSettingsStorage(dbFile, logger)
	gameStorage := infrastructure.CreateGameStorage(dbFile)
	rand := infrastructure.CreateMathRand()
	gameFlow := usecases.CreateGameFlow(logger, localizer, gameStorage, rand)
	telebot.AddHandler("/pidorules", gameFlow.Rules)
	telebot.AddHandler("/pidoreg", gameFlow.Add)
	telebot.AddHandler("/pidor", gameFlow.Play)
	telebot.AddHandler("/pidorstats", gameFlow.Stats)
	telebot.AddHandler("/pidorall", gameFlow.All)
	telebot.AddHandler("/pidorme", gameFlow.Me)

	converter := infrastructure.CreateFfmpegConverter(logger)
	videoFlow := usecases.CreateVideoFlow(logger, converter, converter)
	telebot.AddHandler(videoFlow)

	fileDownloader := infrastructure.CreateFileDownloader()
	remoteMediaSender := helpers.CreateSendMediaStrategy(logger)
	localMediaSender := helpers.CreateUploadMediaStrategy(logger, remoteMediaSender, fileDownloader, converter)

	rabbit, close := infrastructure.CreateRabbitFactory(logger, os.Getenv("AMQP_URL"))
	defer close()
	task := rabbit.NewTask("twitter_queue")

	twitterMediaFactory := api.CreateTwitterMediaFactory(logger, task)
	twitterFlow := usecases.CreateTwitterFlow(logger, twitterMediaFactory, localMediaSender)
	twitterTimeout := usecases.CreateTwitterTimeout(logger, twitterFlow)
	twitterParser := usecases.CreateTwitterParser(logger, twitterTimeout)
	twitterRemoveSourceDecorator := usecases.CreateRemoveSourceDecorator(logger, twitterParser, settingsStorage)
	telebot.AddHandler(twitterRemoveSourceDecorator)

	httpClient := api.CreateHttpClient()
	convertMediaSender := helpers.CreateConvertMediaStrategy(logger, localMediaSender, fileDownloader, converter, converter)
	linkFlow := usecases.CreateLinkFlow(logger, httpClient, converter, convertMediaSender)
	removeLinkSourceDecorator := usecases.CreateRemoveSourceDecorator(logger, linkFlow, settingsStorage)
	telebot.AddHandler(removeLinkSourceDecorator)

	tiktokHttpClient := api.CreateHttpClient() // domain specific headers and cookies
	tiktokJsonApi := api.CreateTikTokJsonAPI(logger, tiktokHttpClient, rand)
	tiktokHtmlApi := api.CreateTikTokHTMLAPI(logger, tiktokHttpClient, rand)
	tiktokApiDecorator := api.CreateTikTokAPIDecorator(tiktokJsonApi, tiktokHtmlApi)
	tiktokMediaFactory := api.CreateTikTokMediaFactory(logger, tiktokApiDecorator)
	tiktokFlow := usecases.CreateTikTokFlow(logger, tiktokHttpClient, tiktokMediaFactory, localMediaSender)
	telebot.AddHandler(tiktokFlow)

	fileUploader := api.CreateTelegraphAPI()
	//TODO: image_downloader := api.CreateTelebotImageDownloader()
	imageFlow := usecases.CreateImageFlow(logger, fileUploader, telebot)
	telebot.AddHandler(imageFlow)

	publisherFlow := usecases.CreatePublisherFlow(logger)
	telebot.AddHandler(publisherFlow)
	telebot.AddHandler("/loh666", publisherFlow.HandleRequest)

	youtubeAPI := api.CreateYoutubeAPI(logger, fileDownloader)
	sendVideoStrategy := helpers.CreateSendVideoStrategy(logger)
	sendVideoStrategySplitDecorator := helpers.CreateSendVideoStrategySplitDecorator(logger, sendVideoStrategy, converter)
	youtubeFlow := usecases.CreateYoutubeFlow(logger, youtubeAPI, youtubeAPI, sendVideoStrategySplitDecorator)
	removeYoutubeSourceDecorator := usecases.CreateRemoveSourceDecorator(logger, youtubeFlow, settingsStorage)
	telebot.AddHandler(removeYoutubeSourceDecorator)

	telebot.AddHandler("/proxy", func(m *core.Message, bot core.IBot) error {
		_, err := bot.SendText("tg://proxy?server=proxy.ailinykh.com&port=443&secret=dd71ce3b5bf1b7015dc62a76dc244c5aec")
		return err
	})

	iDoNotCare := usecases.CreateIDoNotCare()
	telebot.AddHandler(iDoNotCare)
	// Start endless loop
	telebot.Run()
}

func createLogger() (core.ILogger, func()) {
	logFilePath := path.Join(getWorkingDir(), "pullanusbot.log")
	lf, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}

	l := logger.Init("pullanusbot", true, false, lf)
	close := func() {
		lf.Close()
		l.Close()
	}
	return l, close
}

func getWorkingDir() string {
	workingDir := os.Getenv("WORKING_DIR")
	if len(workingDir) == 0 {
		return "pullanusbot-data"
	}
	return workingDir
}
