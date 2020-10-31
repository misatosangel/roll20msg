// Copyright 2020 misatos.angel@gmail.com.  All rights reserved.

package stats

import (
	"time"
)

// simple date-result tuple
type DatedResult struct {
	Date time.Time
	Result int
}
