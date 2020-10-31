// Copyright 2020 misatos.angel@gmail.com.  All rights reserved.

package roll20msg

import (
	"fmt"
	"math"
	"time"
	"encoding/json"
)

type DiceResult struct {
	Value         int         `json:"v,omitempty"`
}

// there may well be other modifiers that exist, but for now we parse these we've seen
type CustomCritMod struct {
	Comparator    string      `json:"comp,omitempty"`
	Point         int64       `json:"point,omitempty"`
}
type KeepMod struct {
	End           string      `json:"end,omitempty"`
	Count         int64       `json:"count,omitempty"`
}

// modifiers hash map. Currently we only know about these
type Mods struct {
	CustomCrit    []CustomCritMod `json:"customCrit,omitempty"`
	Keep          KeepMod     `json:"keep,omitempty"`
}

// Most of the below often don't exist. Type is the only one that is guaranteed
// Most have Expression, but some have text instead.
type Roll struct {
	// the roll expression might be an integer or a string, because ????
	Expression    interface{} `json:"expr,omitempty"`
	Text          string      `json:"text,omitempty"`
	Type          string      `json:"type,omitempty"`
	Dice          int64       `json:"dice,omitempty"`
	Sides         int64       `json:"sides,omitempty"`
	Mods          Mods        `json:"mods,omitempty"`
	Results      []DiceResult `json:"results,omitempty"`
}

type RollResult struct {
	ResultType    string      `json:"resultType,omitempty"`
	Type          string      `json:"type,omitempty"`
	Total         int64       `json:"total,omitempty"`
	Rolls       []Roll        `json:"rolls,omitempty"`
}

type InlineRoll struct {
	Expression    string      `json:"expression,omitempty"`
	Results       RollResult  `json:"results,omitempty"`

	// the signature field is false if not defined rather than empty or missing or null because ????
	RollId        string      `json:"rollid,omitempty"`
	Signature     interface{} `json:"signature,omitempty"`
}

type Msg struct {
	R20DateStamp  float64     `json:".priority,omitempty"`
	// the avatar signature field is false if not defined rather than empty or missing or null because ????
	Avatar        interface{} `json:"avatar,omitempty"`
	Content       string      `json:"content,omitempty"`
	ListenerId    string      `json:"listenerid,omitempty"`
	PlayerId      string      `json:"playerid,omitempty"`
	RollTemplate  string      `json:"rolltemplate,omitempty"`
	Type          string      `json:"type,omitempty"`
	Who           string      `json:"who,omitempty"`
	InlineRolls []InlineRoll  `json:"inlinerolls,omitempty"`
	Target        string      `json:"target,omitempty"`
	TargetName    string      `json:"target_name,omitempty"`
	OriginalRoll  string      `json:"origRoll,omitempty"`

}

type MsgBlock map[string]Msg
type MsgStream []MsgBlock


func (self *Msg) BriefDesc() string {
	return fmt.Sprintf( "%s type %s by %s", self.TimeStamp().Format(time.RFC1123Z), self.Type, self.Who);
}

func (self *Msg) TimeStamp() time.Time {
	// the ".priority" (named R20DateStamp) here seems to be a java powered timestamp in milliseconds,
	// complete with a decimal point because ????
	// so convert with Unix foramt
	i, f := math.Modf(self.R20DateStamp/1000) // gives us seconds and fractional seconds
	nanoSecs, _ := math.Modf(f*1000000000) // convert to nanoseconds
	return time.Unix(int64(i),int64(nanoSecs))
}

func (self *Msg) HasRollResults() (bool, error) {
	err := self.UnpackRolls()
	if err != nil {
		return true, err
	}
	return len(self.InlineRolls) > 0, nil
}


// mesages with "rollresult" type actually have their roll result embedded inside their
// content information. roll20 in next level of sigh here.
func (self *Msg) UnpackRolls() error {
	if len(self.InlineRolls) > 0 || (self.Type != "rollresult" && self.Type != "gmrollresult") {
		return nil
	}
	ir := InlineRoll{
		Expression: self.OriginalRoll,
		Signature: false,
	}
	err := json.Unmarshal([]byte(self.Content), &ir.Results)
	if err != nil {
		return err
	}
	self.InlineRolls = append(self.InlineRolls, ir)
	return nil
}

// will pass each actual dice roll made to the interating function
// function should return true to continue or false to abort
//
// This function will return true if all rolls were iterated,
// or false if it was aborted by the calling funciton.
func (self *Msg) IterateRawDiceRolls(f func(r Roll) bool) (bool, error) {
	hasRoles, err := self.HasRollResults() // will unpack 'rollresult' type if required
	if err != nil || ! hasRoles {
		return false, err
	}
	for _, ir := range self.InlineRolls {
		if len(ir.Results.Rolls) == 0 {
			continue
		}
		for _, roll := range ir.Results.Rolls {
			if ! f(roll) {
				return false, nil
			}
		}
	}
	return true, nil
}

// Some helpers to deal with r20's stupidity of <string> || FALSE
// who programs this garbage?
func (self *Msg) GetAvatar() string {
	if s, ok := self.Avatar.(string) ; ok {
		return s
	}
	return ""
}

func (self *InlineRoll) GetSignature() string {
	if s, ok := self.Signature.(string) ; ok {
		return s
	}
	return ""
}

func (self *Roll) GetExpression() string {
	if s, ok := self.Expression.(string) ; ok {
		return s
	}
	return ""
}
