package common

var DefaultPermissions = Nodes{
	{"level.*", EveryoneLevel},
	{"level.config.*", ManagerLevel},
}
