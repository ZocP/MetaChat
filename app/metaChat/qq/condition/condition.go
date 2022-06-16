package condition

import "github.com/tidwall/gjson"

type Condition struct {
	mp map[string]string
}

func NewCondition(mp map[string]string) *Condition {
	return &Condition{mp: mp}
}

func (c *Condition) Fit(msg gjson.Result) bool {
	for k, v := range c.mp {
		if msg.Get(k).String() != v {
			return false
		}
	}
	return true
}
