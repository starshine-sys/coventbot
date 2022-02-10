ids = {"banana_id", "bnirb_id"}
names = {"banana", "bnirb"}

if scheduled_args then
    remove_role(message.author.id, scheduled_args)
    return
end

require.not_role(string.format("You are already %v'd", names[1]), ids[1])
require.not_role(string.format("You are already %v'd", names[2]), ids[2])

to_set = {ids[1], names[1]}
if math.random(100) == 50 then
    to_set = {ids[2], names[2]}
end

add_role(message.author.id, to_set[1])

send_message(nil, string.format("You are now %v for `5` minutes", to_set[2]))

schedule_cc(nil, "5m", to_set[1])
