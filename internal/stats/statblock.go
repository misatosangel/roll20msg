// Copyright 2020 misatos.angel@gmail.com.  All rights reserved.

package stats

import (
	"sort"
	"fmt"
	"strings"
	"strconv"
)

type StatBlock struct {
	Count uint
	Median float64
	Mean float64
	Mode int
	HasMode bool
	Min int
	Max int
	Total int64
	ByRoll map[int]int
	OrderedByRoll []int
	OrderedByTime []int
}

// takes a set of date-results, orders them by date and then
// fills out all the other stats, returning the stat block
func NewStatBlock( vals []DatedResult ) StatBlock {
	count := len(vals)
	sb := StatBlock{
		Count: uint(count),
		ByRoll: make( map[int]int ),
		OrderedByRoll: make( []int, count, count),
		OrderedByTime: make( []int, count, count),
	}
	if count == 0 {
		return sb
	}
	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i].Date.Before(vals[j].Date)
	})

	allVals := make( []int, 0, 10)
	for i, dateVal := range vals {
		val := dateVal.Result
		if i == 0 {
			sb.Min = val
			sb.Max = val
		} else if ( val < sb.Min ) {
			sb.Min = val
		} else if ( val > sb.Max ) {
			sb.Max = val
		}
		sb.OrderedByTime[i] = val
		sb.Total += int64(val)
		if _, exists := sb.ByRoll[val] ; ! exists {
			sb.ByRoll[val] = 1
			allVals = append(allVals, val)
		} else {
			sb.ByRoll[val] = sb.ByRoll[val] + 1
		}
	}
	sort.Ints(allVals)
	sb.Mean = float64(sb.Total) / float64(sb.Count)
	pos := 0
	largestCnt := 0
	for _, val := range allVals {
		cnt := sb.ByRoll[val]
		if largestCnt == cnt {
			sb.HasMode = false
		} else if cnt > largestCnt {
			largestCnt = cnt
			sb.Mode = val
			sb.HasMode = true
		}
		for i := 0 ; i < cnt ; i++ {
			sb.OrderedByRoll[pos] = val
			pos++
		}
	}

	// finally find the median based on odd/even count
    if count & 1 == 1 {
        half := (count+1) >> 1;
        sb.Median = float64(sb.OrderedByRoll[half-1]);
    } else {
        half := count >> 1;
        sb.Median = float64(sb.OrderedByRoll[half-1] + sb.OrderedByRoll[half]) / 2;
    }
    return sb
}

// Creates a multi-line format string for this stat block
// suitable for sending to discord (with markup)
func (self *StatBlock) FormatResultsDiscord() string {
	mode := "<none>"
	if self.HasMode {
		mode = fmt.Sprintf("%d", self.Mode)
	}

	return fmt.Sprintf(
          "**Count:** %d\n"+
          "**By time:** %s\n"+
          "**By roll:** %s\n"+
          "**Median:** %.2f\n"+
          "**Mode:** %s\n"+
          "**Mean:** %.2f\n"+
          "**Min:** %d (%d)\n"+
          "**Max:** %d (%d)\n",
          self.Count, JoinIntSlice(self.OrderedByTime), JoinIntSlice(self.OrderedByRoll),
          self.Median, mode, self.Mean, self.Min, self.ByRoll[self.Min], self.Max, self.ByRoll[self.Max]);
}

// why golang makes this so hard I have no idea
func JoinIntSlice(ints []int) string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = strconv.Itoa(int(v))
	}

	return strings.Join(strs, ", ")
}
