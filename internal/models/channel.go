package models

type Channel string

const (
	ChannelEmail     Channel = "email"
	ChannelWebsocket Channel = "websocket"
)

var ValidChannels = map[Channel]struct{}{
	ChannelEmail:     {},
	ChannelWebsocket: {},
}

func IsValidChannel(ch string) bool {
	_, ok := ValidChannels[Channel(ch)]
	return ok
}
