package gordie

const (
	OP_Dispatch              uint8 = 0
	OP_Heartbeat             uint8 = 1
	OP_Identify              uint8 = 2
	OP_Presence_Update       uint8 = 3
	OP_Voice_State_Update    uint8 = 4
	OP_Resume                uint8 = 6
	OP_Reconnect             uint8 = 7
	OP_Request_Guild_Members uint8 = 8
	OP_Invalid_Session       uint8 = 9
	OP_Hello                 uint8 = 10
	OP_Heartbeat_ACK         uint8 = 11
)

type EventRaw struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		Type              int    `json:"type"`
		TTS               bool   `json:"tts"`
		Timestamp         string `json:"timestamp"`
		ReferencedMessage struct {
			Content string `json:"content"`
		} `json:"referenced_message"`
		Pinned   bool   `json:"pinned"`
		Nonce    string `json:"nonce"`
		Mentions []struct {
			Username    string `json:"username"`
			PublicFlags int    `json:"public_flags"`
			Member      struct {
				Roles                      []string `json:"roles"`
				PremiumSince               string   `json:"premium_since"`
				Pending                    bool     `json:"pending"`
				Nick                       string   `json:"nick"`
				Mute                       bool     `json:"mute"`
				JoinedAt                   string   `json:"joined_at"`
				Flags                      int      `json:"flags"`
				Deaf                       bool     `json:"deaf"`
				CommunicationDisabledUntil string   `json:"communication_disabled_until"`
				Avatar                     string   `json:"avatar"`
			} `json:"member"`
			Id                   string `json:"id"`
			GlobalName           string `json:"global_name"`
			Discriminator        string `json:"discriminator"`
			Clan                 string `json:"clan"`
			Bot                  bool   `json:"bot"`
			AvatarDecorationData string `json:"avatar_decoration_data"`
			Avatar               string `json:"avatar"`
		} `json:"mentions"`
		MentionRoles    []string `json:"mention_roles"`
		MentionEveryone bool     `json:"mention_everyone"`
		Member          struct {
			Roles                      []string `json:"roles"`
			PremiumSince               string   `json:"premium_since"`
			Pending                    bool     `json:"pending"`
			Nick                       string   `json:"nick"`
			Mute                       bool     `json:"mute"`
			JoinedAt                   string   `json:"joined_at"`
			Flags                      int      `json:"flags"`
			Deaf                       bool     `json:"deaf"`
			CommunicationDisabledUntil string   `json:"communication_disabled_until"`
			Avatar                     string   `json:"avatar"`
		} `json:"member"`
		Id              string     `json:"id"`
		Flags           int        `json:"flags"`
		Embeds          []struct{} `json:"embeds"`
		EditedTimestamp string     `json:"edited_timestamp"`
		Content         string     `json:"content"`
		Components      []struct{} `json:"components"`
		ChannelId       string     `json:"channel_id"`
		Author          struct {
			Username             string `json:"username"`
			PublicFlags          int    `json:"public_flags"`
			Id                   string `json:"id"`
			GlobalName           string `json:"global_name"`
			Discriminator        string `json:"discriminator"`
			Clan                 string `json:"clan"`
			AvatarDecorationData string `json:"avatar_decoration_data"`
			Avatar               string `json:"avatar"`
			Bot                  bool   `json:"bot"`
		} `json:"author"`
		User struct {
			Verified      bool   `json:"verified"`
			Username      string `json:"username"`
			MFAEnabled    bool   `json:"mfa_enabled"`
			Id            string `json:"id"`
			GlobalName    string `json:"global_name"`
			Flags         int    `json:"flags"`
			Email         string `json:"email"`
			Discriminator string `json:"discriminator"`
			Clan          string `json:"clan"`
			Avatar        string `json:"avatar"`
		} `json:"user"`
		Attachments []struct{} `json:"attachments"`
		GuildId     string     `json:"guild_id"`
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
		Author: struct {
			Username             string
			PublicFlags          int
			Id                   string
			GlobalName           string
			Discriminator        string
			Clan                 string
			AvatarDecorationData string
			Avatar               string
			Bot                  bool
		}(e.D.Author),
		User: struct {
			Verified      bool
			Username      string
			MFAEnabled    bool
			Id            string
			GlobalName    string
			Flags         int
			Email         string
			Discriminator string
			Clan          string
			Avatar        string
		}(e.D.User),
	}
}

type Event struct {
	T         string
	S         int
	Op        int
	Content   string
	GuildId   string
	ChannelId string
	Author    struct {
		Username             string
		PublicFlags          int
		Id                   string
		GlobalName           string
		Discriminator        string
		Clan                 string
		AvatarDecorationData string
		Avatar               string
		Bot                  bool
	}
	User struct {
		Verified      bool
		Username      string
		MFAEnabled    bool
		Id            string
		GlobalName    string
		Flags         int
		Email         string
		Discriminator string
		Clan          string
		Avatar        string
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
