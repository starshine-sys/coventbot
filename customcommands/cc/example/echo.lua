require.role("You're not a mod, you can't do that!", {"Mods", "Admin"})

target, ok = args.next_channel_check()
if not ok then
    target = channel
end

if not args.has_next() then
    send_message(nil, "You must give something to echo.")
    return
end

send_message(target.id, args.remainder())
react(nil, nil, "✅")
