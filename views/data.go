package views

import (
	"github.com/gitalek/gogal/models"
	"log"
)

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
	// AlertMsgGeneric is displayed when any random error is encountered by backend.
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

type PublicError interface {
	error
	Public() string
}

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

type Alert struct {
	Level   string
	Message string
}

func (d *Data) SetAlert(err error) {
	var msg string
	if pubErr, ok := err.(PublicError); ok {
		msg = pubErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}
