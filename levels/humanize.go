// SPDX-License-Identifier: AGPL-3.0-only
package levels

import (
	"fmt"
	"strconv"
)

func HumanizeInt64(i int64) string {
	if i >= -1000 && i <= 1000 {
		return strconv.FormatInt(i, 10)
	}

	s := fmt.Sprintf("%.1fk", float64(i)/1000)
	if i < 0 {
		return "-" + s
	}
	return s
}
