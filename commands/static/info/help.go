package info

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) help(ctx *bcr.Context) (err error) {
	if ctx.RawArgs != "" {
		return ctx.Help(ctx.Args)
	}

	showAll, _ := ctx.Flags.GetBool("all")
	cmds := bot.Router.Commands()
	sort.Sort(bcr.Commands(cmds))

	prefix := ctx.Prefix
	if strings.Contains(prefix, "@") {
		prefix = "@" + ctx.Router.Bot.Username + " "
	}
	footer := fmt.Sprintf("For help with a command, use %vhelp <command> | ", prefix)

	descs := make([]string, 0, len(cmds))
	if bot.Config.Branding.SupportServer != "" {
		descs = append(descs, "For help and feedback, join the support server: "+bot.Config.Branding.SupportServer+"\n\n")
	}

	// user perms
	var perms discord.Permissions
	if ctx.Guild != nil && ctx.Member != nil {
		for _, gr := range ctx.Guild.Roles {
			for _, ur := range ctx.Member.RoleIDs {
				if gr.ID == ur {
					perms |= gr.Permissions
					if gr.Permissions.Has(discord.PermissionAdministrator) {
						perms |= discord.PermissionAll
					}
				}
			}
		}
	} else {
		// dm permissions
		perms = discord.PermissionAddReactions | discord.PermissionAttachFiles | discord.PermissionEmbedLinks | discord.PermissionReadMessageHistory | discord.PermissionSendMessages | discord.PermissionUseExternalEmojis | discord.PermissionUseExternalStickers | discord.PermissionUseSlashCommands | discord.PermissionViewChannel
	}

	userLevel := 0
	if isAdmin, _ := bot.AdminRole.Check(ctx); isAdmin {
		userLevel = 3
	} else if isMod, _ := bot.ManagerRole.Check(ctx); isMod {
		userLevel = 2
	} else if isHelper, _ := bot.ModeratorRole.Check(ctx); isHelper {
		userLevel = 1
	}

	for _, cmd := range cmds {
		if !cmd.Hidden && (showAll || (bot.permLevel(cmd) <= userLevel && perms.Has(cmd.Permissions|cmd.GuildPermissions))) {
			descs = append(descs, fmt.Sprintf("`%v`: %v\n", cmd.Name, cmd.Summary))
		}
	}

	embeds := bcr.StringPaginator("Help", bcr.ColourBlurple, descs, 15)
	for i := range embeds {
		embeds[i].Footer.Text = footer + embeds[i].Footer.Text
	}

	_, _ = bot.PagedEmbed(ctx, embeds, 15*time.Minute)
	return err
}

func (bot *Bot) permLevel(cmd *bcr.Command) int {
	if cmd.OwnerOnly {
		return 4
	}

	switch cmd.CustomPermissions {
	case bot.AdminRole:
		return 3
	case bot.ManagerRole:
		return 2
	case bot.ModeratorRole:
		return 1
	}

	return 0
}
