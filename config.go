package miauth

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	//"github.com/gin-gonic/gin/binding"
	"gopkg.in/yaml.v2"
)

type DOSMJEmailToInputPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DOSMJEmailPayload struct {
	TemplateName string `yaml:"template_name" json:"template_name"`
	TemplateData struct {
		ResetLink string `yaml:"reset_link" json:"reset_link"`
	} `yaml:"template_data" json:"template_data"`
	EmailSpecs struct {
		Subject string                     `json:"subject"`
		To      []DOSMJEmailToInputPayload `json:"to"`
	} `yaml:"email_specs" json:"email_specs"`
}
type ConfigMiauth struct {
	Name                    string
	PublicForgotPasswordURL string `yaml:"public_forgot_password_url"`
	Port                    string
	BCrypt                  struct {
		Salt string `yaml:"salt"`
	} `yaml:"bcrypt"`
	AccessToken struct {
		Secret    string
		ExpiresIn string `yaml:"expires_in"`
	} `yaml:"access_token"`
	RefreshToken struct {
		Secret string
	} `yaml:"refresh_token"`
	ResetPassword struct {
		ExpiresIn   string `yaml:"expires_in"`
		Secret      string
		MailService struct {
			DOSMJ struct {
				Method   string
				Endpoint string
				Payload  DOSMJEmailPayload
			}
		} `yaml:"mail_service"`
	} `yaml:"reset_password"`
	DB struct {
		Postgres string
	}
	FieldValidations struct {
		Username struct {
			Pattern                    string `yaml:"pattern"`
			Len                        []int  `yaml:"len"`
			InvalidPatternErrorMessage string `yaml:"invalid_pattern_error_message"`
		}
		Password struct {
			Len                        []int  `yaml:"len"`
			InvalidPatternErrorMessage string `yaml:"invalid_pattern_error_message"`
		}
		Email struct {
			Len                        []int  `yaml:"len"`
			InvalidPatternErrorMessage string `yaml:"invalid_pattern_error_message"`
		}
	} `yaml:"field_validations"`
}

var Config *ConfigMiauth
var UsernameValidator validator.Func
var EmailValidator validator.Func

func InitConfig() {
	configFilePath := flag.String("config", "miauth.config.v2.yml", "miauth yaml config file")
	flag.Parse()

	bytes, err := ioutil.ReadFile(*configFilePath)
	if err == nil {
		if err := ReadConfig(string(bytes)); err != nil {
			log.Fatal(err)
		}
	} else {
		dataConfig := os.Getenv("MIAUTH_CONFIG")
		if err := ReadConfig(dataConfig); err != nil {
			log.Fatal(err)
		}
	}
	
	usernameMinLen := Config.FieldValidations.Username.Len[0]
	usernameMaxLen := Config.FieldValidations.Username.Len[1]
	UsernameValidator = func (fl validator.FieldLevel) bool  {
		if !fl.Field().IsValid() {
			return false
		}
		parent := fl.Parent()
		var payload reflect.Value
		if parent.Kind() == reflect.Ptr {
			payload =	parent.Elem()
		} else if parent.Kind() == reflect.Struct {
			payload = parent
		} else {
			return false
		}
		if kindAuth := payload.FieldByName("Kind").String(); kindAuth == "miauth" {
			switch fl.Field().Kind() {
			case reflect.String:
				if len(fl.Field().String()) < usernameMinLen || len(fl.Field().String()) > usernameMaxLen {
					return false
				}
				return UsernamePatterns[Config.FieldValidations.Username.Pattern].MatchString(fl.Field().String())
			}
			return false
		} else {
			return true
		}
	}

	emailMinLen := Config.FieldValidations.Email.Len[0]
	emailMaxLen := Config.FieldValidations.Email.Len[1]
	EmailValidator = func (fl validator.FieldLevel) bool  {
		if !fl.Field().IsValid() {
			return false
		}
		switch fl.Field().Kind() {
		case reflect.String:
			if len(fl.Field().String()) < emailMinLen || len(fl.Field().String()) > emailMaxLen {
				return false
			}
			return RxEmail.MatchString(fl.Field().String())
		}
		return false
	}

	
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("miauth_username", UsernameValidator)
		v.RegisterValidation("miauth_email", EmailValidator)
	}
}

func ReadConfig(data string) error {
	if err := yaml.Unmarshal([]byte(data), &Config); err != nil {
		return err
	}
	return nil
}
