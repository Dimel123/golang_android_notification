package models

type Notification struct {
	Account_id      string
	Event_type      string
	Is_popup        bool
	Event_details   string
	Entity          string
	Event_uid       string
}

func NewNotification(accountId string, eventType string, isPopup bool, eventDetails string, entity string, eventUid string) *Notification {
	return &Notification{accountId, eventType, isPopup, eventDetails, entity, eventUid}
}