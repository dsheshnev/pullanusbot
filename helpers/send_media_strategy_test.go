package helpers_test

import (
	"testing"

	"github.com/ailinykh/pullanusbot/v2/core"
	"github.com/ailinykh/pullanusbot/v2/helpers"
	"github.com/ailinykh/pullanusbot/v2/test_helpers"
	"github.com/stretchr/testify/assert"
)

func Test_SendMedia_DoesNotFailOnEmptyMedia(t *testing.T) {
	strategy, bot := makeMediaStrategySUT()
	media := []*core.Media{}

	strategy.SendMedia(media, bot)

	assert.Equal(t, []string{}, bot.SentMedias)
}

func Test_SendMedia_SendsASingleMediaTroughABot(t *testing.T) {
	strategy, bot := makeMediaStrategySUT()
	media := []*core.Media{{ResourceURL: "https://a-url.com"}}

	strategy.SendMedia(media, bot)

	assert.Equal(t, []string{"https://a-url.com"}, bot.SentMedias)
}

func Test_SendMedia_SendsAGroupMediaTroughABot(t *testing.T) {
	strategy, bot := makeMediaStrategySUT()
	media := []*core.Media{{ResourceURL: "https://a-url.com"}, {ResourceURL: "https://another-url.com"}}

	strategy.SendMedia(media, bot)

	assert.Equal(t, []string{"https://a-url.com", "https://another-url.com"}, bot.SentMedias)
}

// Helpers
func makeMediaStrategySUT() (core.ISendMediaStrategy, *test_helpers.FakeBot) {
	logger := test_helpers.CreateFakeLogger()
	strategy := helpers.CreateSendMediaStrategy(logger)
	bot := test_helpers.CreateFakeBot()
	return strategy, bot
}
