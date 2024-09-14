package state

type Error struct{
	Message string
}

func (err *Error) Set(message string) {
	err.Message = message
}

func (err *Error) Error() string {
	return err.Message
}