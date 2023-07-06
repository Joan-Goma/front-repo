package views

import (
	"html/template"
)

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random error is encountered
	AlertMsgGeneric = "Something went wrong." + "\n" + "Please try again, and contact us if the problem persists."
)

type Alert struct {
	Level   string
	Message template.HTML
}

type Data struct {
	Active string
	Alert  *Alert
	CSRF   template.HTML
	Yield  interface{}
}

func (d *Data) SetAlert(err string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: template.HTML(err),
	}
}
