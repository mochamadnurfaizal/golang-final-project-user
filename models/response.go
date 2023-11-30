package models

type Response struct {
	Messages string      `swagger:"description(Message)" example:"string"`
	Success  bool        `swagger:"description(Success)" example:"true"`
	Data     interface{} `swagger:"description(Data)"`
}
