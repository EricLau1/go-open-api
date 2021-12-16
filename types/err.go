package types

type Error struct {
	Message string `json:"error"`
}

func NewError(err error) interface{} {
	e := &Error{Message: "error"}
	if err != nil {
		e.Message = err.Error()
	}
	return e
}
