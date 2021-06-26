package roles

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// Category is a role category
type Category struct {
	ID       int64
	ServerID discord.GuildID

	Name        string
	Description string
	Colour      discord.Color

	RequireRole discord.RoleID
	Roles       []uint64
}

func (bot *Bot) guildCategories(id discord.GuildID) (cats []Category, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &cats, "select * from role_categories where server_id = $1 order by name", id)
	return
}

func (bot *Bot) categoryID(guildID discord.GuildID, id int64) (c Category, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &c, "select * from role_categories where server_id = $1 and id = $2", guildID, id)
	return
}

func (bot *Bot) categoryName(id discord.GuildID, name string) (c Category, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &c, "select * from role_categories where server_id = $1 and name ilike $2 limit 1", id, name)
	return
}

func (bot *Bot) newCategory(id discord.GuildID, name, description string, requireRole discord.RoleID, colour discord.Color) (c *Category, err error) {
	c = &Category{}
	err = pgxscan.Get(context.Background(), bot.DB.Pool, c, `insert into role_categories
	(server_id, name, description, require_role, colour) values ($1, $2, $3, $4, $5)
	on conflict (server_id, name) do update
	set description = $3, require_role = $4, colour = $5
	returning *`, id, name, description, requireRole, colour)
	if err != nil {
		return nil, err
	}
	return
}

func (bot *Bot) categoryRole(guildID discord.GuildID, roleID discord.RoleID) (c Category, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &c, "select * from role_categories where server_id = $1 and $2 = any(roles)", guildID, roleID)
	return
}

func (bot *Bot) categoryRoles(id int64, roles []uint64) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "update role_categories set roles = $1 where id = $2", roles, id)
	return
}
