package gordie

type SlashCommand struct {
	Name        string
	Type        int
	Description string
	Options     []SlashCommandOption
}

type SlashCommandOption struct {
	Type        int
	Name        string
	Description string
	Required    bool
}
