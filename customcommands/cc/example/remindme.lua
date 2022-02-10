-- a bad reminder cc
local offset = args.pop()
local text = args.remainder()

if scheduled_args then
    send_message(nil, scheduled_args)
    return
end

local id = schedule_cc(nil, offset, string.format("<@!%s>: %s", message.author.id, text))
send_message(nil, string.format("Ok %s, I'll remind you about %s in %s. (%d)", message.author.tag, text, offset, id))
