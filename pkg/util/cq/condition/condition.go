package condition

import (
	"github.com/tidwall/gjson"
	"reflect"
	"strconv"
)

type Condition struct {
	mp map[string]string
}

func NewCondition(mp map[string]string) *Condition {

	return &Condition{mp: mp}
}

//判断某个消息是否符合该条件
func (c *Condition) Fit(msg gjson.Result) bool {
	for k, v := range c.mp {
		value := msg.Get(k).Value()
		if value == nil {
			return false
		}
		switch t := reflect.TypeOf(value); t.Kind() {
		case reflect.Float64:
			intv := int64(value.(float64))
			if strconv.FormatInt(intv, 10) == v {
				return true
			}
		}
		return v == value
	}
	return false
}
