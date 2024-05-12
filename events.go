package gordie

type EventRaw struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		Content   string `json:"content"`
		GuildId   string `json:"guild_id"`
		ChannelId string `json:"channel_id"`
		User      struct {
			Username      string `json:"username"`
			Id            string `json:"id"`
			Discriminator string `json:"discriminator"`
		} `json:"user"`
	} `json:"d"`
}

func (e EventRaw) ToEvent() Event {
	return Event{
		T:         e.T,
		S:         e.S,
		Op:        e.Op,
		Content:   e.D.Content,
		GuildId:   e.D.GuildId,
		ChannelId: e.D.ChannelId,
		User: struct {
			Username      string
			Id            string
			Discriminator string
		}{
			Username:      e.D.User.Username,
			Id:            e.D.User.Id,
			Discriminator: e.D.User.Discriminator,
		},
	}
}

type Event struct {
	T         string
	S         int
	Op        int
	Content   string
	GuildId   string
	ChannelId string
	User      struct {
		Username      string
		Id            string
		Discriminator string
	}
}

type HelloEventRaw struct {
	Op int `json:"op"`
	D  struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}

func (e HelloEventRaw) ToHelloEvent() HelloEvent {
	return HelloEvent{
		Op:                e.Op,
		HeartBeatInterval: e.D.HeartbeatInterval,
	}

}

type HelloEvent struct {
	Op                int
	HeartBeatInterval int
}
