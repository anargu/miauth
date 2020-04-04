package miauth

//  importing deps
import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

var (
	WelcomeTemplate = "welcome_template.templ"
	client          = resty.New()
)

type DM map[string]interface{}
type OkDOSMJResponse struct {
	Message string
}

type ErrorDOSMJResponse struct {
	Error string
}
type MailPayloadData struct {
	ResetLink string `json:"reset_link"`
}

func Send(
	data *MailPayloadData,
	email string,
	name string,
	subject string) error {

	payload := Config.ResetPassword.MailService.DOSMJ.Payload
	payload.EmailSpecs.To.Email = email
	payload.EmailSpecs.To.Name = name
	payload.EmailSpecs.Subject = subject
	payload.TemplateData.ResetLink = data.ResetLink

	rawResult, err := client.R().
		SetBody(payload).
		SetResult(OkDOSMJResponse{}).
		SetError(ErrorDOSMJResponse{}).
		Post(fmt.Sprintf("%s", Config.ResetPassword.MailService.DOSMJ.Endpoint))
	if err != nil {
		return err
	}
	if rawResult.IsError() {
		dosmjError := rawResult.Error().(*ErrorDOSMJResponse)
		return errors.New(dosmjError.Error)
	}
	return nil
}

func SendResetPassword(user *User, resetLink string, subject string) error {
	err := Send(&MailPayloadData{ResetLink: resetLink}, user.Email, fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		subject,
	)
	return err
}
