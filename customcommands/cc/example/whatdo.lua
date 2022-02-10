-- configuration
local vote_channel = "757193930084319393"
local offset = "18h"

-- code

-- vote reminder
if scheduled_args then
    send_message(nil, scheduled_args)
    return
end

require.role("You're not a mod, you can't do that!", {"Mods", "Admin"})

if not args.has_next() then
    send_message(nil, "No vote text given.")
    return
end

local vote_text = args.remainder()

local msg_id = send_message(vote_channel, string.format("*See <#%s>/#%s for context*\nLink to original message: <https://discord.com/channels/%s/%s/%s>\n**Vote for:** %s", channel.id, channel.name, guild.id, channel.id, message.id, vote_text))

send_message(nil, string.format("Vote created in <#%s> <3\nLink to vote: <https://discord.com/channels/%s/%s/%s>", vote_channel, guild.id, vote_channel, msg_id))

schedule_cc(nil, offset, string.format("The voting period for the following vote ends in 6 hours!\nLink to vote: <https://discord.com/channels/%s/%s/%s>\n**Vote for:** %s", guild.id, vote_channel, msg_id, vote_text))

react(nil, nil, "✅")
