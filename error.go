package miauth

import "fmt"

type ErrorMessage struct {
	Message     string `json:"error_description"`
	UserMessage string `json:"user_message"`
	Name        string `json:"error"`
	Code        int    `json:"code"`
}

func (err *ErrorMessage) Error() string {
	return fmt.Sprintf("Error: %s. Descripcion: %s.", err.Name, err.Message)
}