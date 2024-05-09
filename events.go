package gordie

type Event struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		Content   string `json:"content"`
		GuildId   string `json:"guild_id"`
		ChannelId string `json:"channel_id"`
	} `json:"d"`
}

type HelloEvent struct {
	Op int `json:"op"`
	D  struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}
