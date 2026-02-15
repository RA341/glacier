package assets

type Config struct {
	UseYtDlp      bool   `yaml:"useYtDlp" env:"USE_YTDLP" default:"false" help:"use ytdlp to download game metadata assets"`
	YTRelayUrl    string `yaml:"YTRelayUrl" env:"YT_RELAY_URL" default:"-" help:"url for the YtRelay instance"`
	YTRelayApiKey string `yaml:"YTRelayApiKey" env:"YT_RELAY_API_KEY" default:"-" help:"optionally api key to be sent with requests" hide:"true"`
}
