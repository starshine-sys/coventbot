package db

import (
	"time"

	"github.com/starshine-sys/snowflake/v2"
)

var sfGen = snowflake.NewSmallGen(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
