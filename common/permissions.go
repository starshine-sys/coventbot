package common

import (
	"sort"
	"strings"
)

type PermissionLevel int

const (
	InvalidLevel   PermissionLevel = 0
	EveryoneLevel  PermissionLevel = 1
	ModeratorLevel PermissionLevel = 2
	ManagerLevel   PermissionLevel = 3
	AdminLevel     PermissionLevel = 4
)

// Node is a permission node.
type Node struct {
	Name  string
	Level PermissionLevel
}

// Matches returns if s matches this node.
// s matches if n.Name is identical to s, or if n is a wildcard node and s starts with n.Name.
func (n Node) Matches(s string) bool {
	if n.Name == s {
		return true
	}

	if !n.IsWildcard() {
		return false
	}

	unwildcard := strings.TrimSuffix(n.Name, "*")
	return unwildcard != "" && (strings.HasPrefix(s, unwildcard) || s == strings.TrimSuffix(unwildcard, "."))
}

// IsWildcard returns true if this node is a wildcard node.
func (n Node) IsWildcard() bool { return strings.HasSuffix(n.Name, ".*") }

func (n Node) len() int {
	unwildcard := strings.TrimSuffix(n.Name, ".*")
	return len(strings.Split(unwildcard, "."))
}

type Nodes []Node

var _ sort.Interface = Nodes(nil)

func (ns Nodes) Len() int { return len(ns) }

func (ns Nodes) Less(i, j int) bool {
	return ns[i].len() > ns[j].len()
}

func (ns Nodes) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

// NodeFor returns the most specific permission node for the given command.
// It assumes ns has already been sorted with sort.Sort.
// If no valid nodes are found, returns InvalidLevel.
func (ns Nodes) NodeFor(s string) Node {
	for _, node := range ns {
		if node.Matches(s) {
			return node
		}
	}

	// return the default permission level
	// (doesn't call DefaultPermissions.NodeFor to prevent an infinite loop in case of invalid input)
	return defaultNodeFor(s)
}

func defaultNodeFor(s string) Node {
	for _, node := range DefaultPermissions {
		if node.Matches(s) {
			return node
		}
	}

	return Node{Name: s, Level: InvalidLevel}
}
