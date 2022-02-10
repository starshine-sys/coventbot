## Functions

```lua
-- arguments:
-- 1. channel ID to send to (nil for invoking channel) (string, number)
-- 2. content to send (string, or table with "content" and "embed" key)
send_message(channel.id, content) -- returns message ID

-- arguments:
-- 1. channel ID to react in (nil for invoking channel) (string, number)
-- 2. message ID to react to (nil for invoking message) (string, number)
-- 3. API emoji to react with (string)
react(channel.id, message.id, reaction)
```

## Global variables

```lua
guild -- invoking guild with id, name, icon, icon_url properties

channel -- invoking channel with id, guild_id, parent_id, type, nsfw, position, name, topic properties

message -- invoking message with id, channel_id, guild_id, pinned, mention_everyone, author, webhook_id, content properties
```

## Examples

Send a message to a target channel:

```lua
-- get target channel, don't raise error on failure
target, ok = args.next_channel_check()
if not ok then
    send_message(nil, "Channel not found.")
    return
end

-- send message to target channel
send_message(target.id, string.format("heya from <#%s>!", channel.id))
-- react to original message
react(nil, nil, "✅")
```

Reply with "hello":

```lua
send_message(nil, string.format("Hello there, %s!", message.author.tag))
```