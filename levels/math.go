package levels

import "math"

func (sc Server) CalculateExp(level int64) (xp int64) {
	return int64(5.0 / 6.0 * float64(level) * (2*math.Pow(float64(level), 2) + 27*float64(level) + 91))
}

func (sc Server) CalculateLevel(xp int64) (lvl int64) {
	x := float64(xp + 1)
	pow := math.Cbrt(
		math.Sqrt(3)*math.Sqrt(3888.0*math.Pow(x, 2)+(291600.0*x)-207025.0) - 108.0*x - 4050.0,
	)

	res := (-pow/(2.0*math.Pow(3.0, 2.0/3.0)*math.Pow(5.0, 1.0/3.0)) -
		(61.0*math.Cbrt(5.0/3.0))/(2.0*pow) - (9.0 / 2.0))

	return int64(res)
}

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
