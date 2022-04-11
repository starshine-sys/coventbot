package common

var DefaultPermissions = Nodes{
	{"level.*", EveryoneLevel},
	{"level.config.*", ManagerLevel},
	{"permissions.*", AdminLevel},
}

var defaultPermsAreSorted = false
