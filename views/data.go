package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
	// AlertMsgGeneric is displayed when any random error is encountered by backend.
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

type Data struct {
	Alert *Alert
	Yield interface{}
}

type Alert struct {
	Level   string
	Message string
}
