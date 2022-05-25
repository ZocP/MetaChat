package eventBridge

type CQEvent struct {
}

func UnmarshalCQEvent(data []byte) (CQEvent, error) {
	var event CQEvent

	return event, nil
}
