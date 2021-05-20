package levels

// expForNextLevel gets the XP needed for the *NEXT* level after the one given
func expForNextLevel(level int64) (xp int64) {
	for i := level; i >= 0; i-- {
		if i != 0 {
			xp += 20 * i
		} else {
			xp += 25
		}
	}

	return
}

func currentLevel(xp int64) (lvl int64) {
	for i := int64(0); ; i++ {
		if expForNextLevel(i) > xp {
			lvl = i
			break
		}
	}

	return
}
