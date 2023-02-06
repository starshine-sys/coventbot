// SPDX-License-Identifier: AGPL-3.0-only
package db

import (
	"context"
	"sort"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/tribble/common"
)

// Permissions returns the given guild's permission overrides.
func (db *DB) Permissions(guildID discord.GuildID) (ns common.Nodes, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &ns, "select name, level from permission_nodes where guild_id = $1", guildID)
	if err != nil {
		return nil, errors.Cause(err)
	}

	// sort nodes before returning
	sort.Sort(ns)
	return ns, nil
}

// SetPermissions adds a guild-level override for the given node.
func (db *DB) SetPermissions(guildID discord.GuildID, node string, level common.PermissionLevel) (err error) {
	_, err = db.Pool.Exec(context.Background(), `insert into permission_nodes
(guild_id, name, level) values ($1, $2, $3)
on conflict (guild_id, name) do update
set level = $3`, guildID, node, level)
	return err
}

// ResetPermissions removes the override for the given node.
func (db *DB) ResetPermissions(guildID discord.GuildID, node string) (err error) {
	_, err = db.Pool.Exec(context.Background(), `delete from permission_nodes where guild_id = $1 and name = $2`, guildID, node)
	return err
}
