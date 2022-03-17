package bot

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func (bot *Bot) Commands() []api.CreateCommandData { return Commands }

var Commands = []api.CreateCommandData{
	// chat input commands
	{
		Name:        "bubble",
		Description: "Bubble wrap!",
		Type:        discord.ChatInputCommand,
		Options: discord.CommandOptions{
			&discord.IntegerOption{
				OptionName:  "size",
				Description: "Size of the bubble wrap (default 10)",
				Min:         option.NewInt(1),
				Max:         option.NewInt(13),
			},
			&discord.BooleanOption{
				OptionName:  "prepop",
				Description: "Whether to pre-pop some bubbles",
			},
			&discord.BooleanOption{
				OptionName:  "ephemeral",
				Description: "Whether to send the bubble wrap as a message only visible to you",
			},
		},
	},
	{
		Name:        "linkto",
		Description: "Move a conversation to another channel.",
		Type:        discord.ChatInputCommand,
		Options: discord.CommandOptions{
			&discord.ChannelOption{
				OptionName:   "channel",
				ChannelTypes: []discord.ChannelType{discord.GuildText, discord.GuildNews, discord.GuildPublicThread, discord.GuildPrivateThread, discord.GuildNewsThread},
				Required:     true,
				Description:  "The channel to link to.",
			},
			&discord.StringOption{
				OptionName:  "topic",
				Required:    false,
				Description: "The topic.",
			},
		},
	},
	{
		Name:        "pride",
		Description: "Add a pride flag circle to your profile picture!",
		Type:        discord.ChatInputCommand,
		Options: discord.CommandOptions{
			&discord.StringOption{
				OptionName:  "flag",
				Description: "Which flag to use.",
				Required:    false,
			},
			&discord.StringOption{
				OptionName:  "pk-member",
				Description: "Which PluralKit member's avatar to add a pride flag to.",
				Required:    false,
			},
		},
	},
	{
		Name:        "sampa",
		Description: "Convert X-SAMPA to IPA.",
		Type:        discord.ChatInputCommand,
		Options: discord.CommandOptions{
			&discord.StringOption{
				OptionName:  "text",
				Description: "The text to convert to IPA.",
				Required:    true,
			},
		},
	},
	{
		Name:        "reminders",
		Description: "Show your reminders.",
		Type:        discord.ChatInputCommand,
	},
	{
		Name:        "remindme",
		Description: "Show your reminders.",
		Type:        discord.ChatInputCommand,
		Options: discord.CommandOptions{
			&discord.StringOption{
				OptionName:  "when",
				Description: "When or in how long to remind you.",
				Required:    true,
			},
			&discord.StringOption{
				OptionName:  "text",
				Description: "What to remind you of.",
			},
		},
	},

	// user context menu commands
	{
		Name: "Show user avatar",
		Type: discord.UserCommand,
	},
}
