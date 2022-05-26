package eventBridge

import "encoding/json"

func UnmarshalCQEvent(data []byte) (CQEvent, error) {
	var event CQEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return event, err
	}
	return event, nil
}
