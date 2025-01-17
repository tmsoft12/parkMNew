package camera

import "time"

type CapturedEventData struct {
	EventID          string    `json:"event_id"`
	EventDescription string    `json:"event_description"`
	EventComment     string    `json:"event_comment"`
	ChannelName      string    `json:"channel_name"`
	EventData        string    `json:"event_data"`
	CapturedTime     time.Time `json:"captured_time"`
}
