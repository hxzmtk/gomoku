package errex

type Item struct {
	Code    int
	Message string
}

func (item Item) Error() string {
	return item.Message
}

func Create(code int,msg string) Item {
	return Item{
		Code: code,
		Message: msg,
	}
}
