package common

var DefaultPermissions = Nodes{
	// level commands
	{"level.*", EveryoneLevel},
	{"level.setxp", ManagerLevel},
	{"level.setlvl", ManagerLevel},
	{"level.config.*", ManagerLevel},
	{"level.config.import", AdminLevel},
	{"level.background.server", ManagerLevel},
	{"nolevels.*", ModeratorLevel},
	{"leaderboard", EveryoneLevel},

	// tickets
	{"tickets.*", EveryoneLevel},
	{"tickets.config.*", ManagerLevel},
	{"tickets.list", ManagerLevel},

	// roles
	{"roles", EveryoneLevel},
	{"roles.config.*", ManagerLevel},
	{"role", EveryoneLevel},
	{"role.info", EveryoneLevel},
	{"role.dump", ModeratorLevel},
	{"role.create", ManagerLevel},
	{"derole", EveryoneLevel},

	// quotes
	{"quote.*", EveryoneLevel},
	{"quote.delete", ModeratorLevel},
	{"quotes", EveryoneLevel},
	{"quotes.leaderboard", EveryoneLevel},
	{"quotes.toggle", ManagerLevel},
	{"quotes.messages", ManagerLevel},

	// self-mute/pause
	{"muteme", EveryoneLevel},
	{"pauseme", EveryoneLevel},
	{"muteme.message", ManagerLevel},

	// moderation
	{"warn", ModeratorLevel},
	{"unban", ManagerLevel},
	{"setnote", ModeratorLevel},
	{"reason", ModeratorLevel},
	{"purge", ModeratorLevel},
	{"notes", ModeratorLevel},
	{"modlog", ModeratorLevel},
	{"members", ModeratorLevel},
	{"massban", ManagerLevel},
	{"makeinvite", ManagerLevel},
	{"lockdown", ManagerLevel},
	{"exportemotes", ManagerLevel},
	{"embed.*", ManagerLevel},
	{"echo.*", ManagerLevel},
	{"delnote", ModeratorLevel},
	{"bgc", ModeratorLevel},
	{"ban", ManagerLevel},
	{"approve", ManagerLevel},
	{"addemoji", ManagerLevel},
	{"channelban", ManagerLevel},
	{"unchannelban", ManagerLevel},

	// configuration
	{"muterole", ManagerLevel},
	{"pauserole", ManagerLevel},
	{"reactroles.*", ManagerLevel},
	{"prefix.*", ManagerLevel},
	{"permissions.*", AdminLevel},
	{"watchlist.*", ManagerLevel},
	{"triggers.*", ManagerLevel},
	{"starboard.*", ManagerLevel},
	{"slowmode.*", ManagerLevel},
	{"modlog.import", ManagerLevel},
	{"modlog.export", ManagerLevel},
	{"modlog.channel", ManagerLevel},
	{"keyrole.*", ManagerLevel},
	{"cc.*", ManagerLevel},
	{"approval.*", ManagerLevel},

	// gatekeeper
	{"captcha.*", ManagerLevel},
	{"agree", EveryoneLevel},

	// user commands
	{"todo.*", EveryoneLevel},
	{"remindme.*", EveryoneLevel},
	{"userinfo", EveryoneLevel},
	{"user-cfg", EveryoneLevel},
	{"transcript", ManagerLevel},
	{"serverinfo", EveryoneLevel},
	{"sampa", EveryoneLevel},
	{"roll", EveryoneLevel},
	{"roleinfo", EveryoneLevel},
	{"reminders", EveryoneLevel},
	{"quickpoll", EveryoneLevel},
	{"pride", EveryoneLevel},
	{"poll", EveryoneLevel},
	{"ping", EveryoneLevel},
	{"message", EveryoneLevel},
	{"meow", EveryoneLevel},
	{"linkto", EveryoneLevel},
	{"invite", EveryoneLevel},
	{"idtime", EveryoneLevel},
	{"help", EveryoneLevel},
	{"getinvite", EveryoneLevel},
	{"enlarge", EveryoneLevel},
	{"embedsource", EveryoneLevel},
	{"delreminder", EveryoneLevel},
	{"complete", EveryoneLevel},
	{"colour", EveryoneLevel},
	{"bubble", EveryoneLevel},
	{"avatar", EveryoneLevel},
	{"about", EveryoneLevel},
}

var defaultPermsAreSorted = false
