package miauth

import (
	"flag"
	"github.com/syssam/go-validator"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

type DOSMJEmailPayload struct {
	TemplateName string `yaml:"template_name" json:"template_name"`
	TemplateData struct {
		ResetLink string `yaml:"reset_link"`
	} `yaml:"template_data" json:"template_data"`
	EmailSpecs struct {
		Subject string `json:"subject"`
		To struct {
			Name string `json:"name"`
			Email string `json:"email"`
		} `json:"to"`
	} `yaml:"email_specs" json:"email_specs"`
}
type ConfigMiauth struct {
	Name string
	PublicForgotPasswordURL string `yaml:"public_forgot_password_url"`
	Port string
	BCrypt struct {
		Salt string `yaml:"salt"`
	} `yaml:"bcrypt"`
	AccessToken struct {
		Secret string
		ExpiresIn string `yaml:"expires_in"`
	} `yaml:"access_token"`
	RefreshToken struct {
		Secret string
	} `yaml:"refresh_token"`
	ResetPassword struct {
		ExpiresIn string `yaml:"expires_in"`
		Secret string
		MailService struct {
			DOSMJ struct {
				Method string
				Endpoint string
				Payload DOSMJEmailPayload
			}
		} `yaml:"mail_service"`
	} `yaml:"reset_password"`
	DB struct {
		Postgres string
	}
	FieldValidations struct {
		Username struct {
			Pattern                    string `yaml:"pattern"`
			Len                        []int `yaml:"len"`
			InvalidPatternErrorMessage string `yaml:"invalid_pattern_error_message"`
		}
		Password struct {
			Len                        []int `yaml:"len"`
			InvalidPatternErrorMessage string `yaml:"invalid_pattern_error_message"`
		}
		Email struct {
			Len                        []int `yaml:"len"`
			InvalidPatternErrorMessage string `yaml:"invalid_pattern_error_message"`
		}
	} `yaml:"field_validations"`
}

var Config *ConfigMiauth

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


	validator.MessageMap["miauth_username"] = Config.FieldValidations.Username.InvalidPatternErrorMessage
	validator.MessageMap["miauth_email"] = Config.FieldValidations.Email.InvalidPatternErrorMessage
	usernameMinLen := Config.FieldValidations.Username.Len[0]
	usernameMaxLen := Config.FieldValidations.Username.Len[1]

	emailMinLen := Config.FieldValidations.Email.Len[0]
	emailMaxLen := Config.FieldValidations.Email.Len[1]

	validator.CustomTypeRuleMap.Set("miauth_username", func(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
		switch v.Kind() {
		case reflect.String:
			if len(v.String()) < usernameMinLen || len(v.String()) > usernameMaxLen {
				return false
			}
			return UsernamePatterns[Config.FieldValidations.Username.Pattern].MatchString(v.String())
		}
		return false
	})
	validator.CustomTypeRuleMap.Set("miauth_email", func(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
		switch v.Kind() {
		case reflect.String:
			if len(v.String()) < emailMinLen || len(v.String()) > emailMaxLen {
				return false
			}
			return RxEmail.MatchString(v.String())
		}
		return false
	})

}

func ReadConfig(data string) error {
	if err := yaml.Unmarshal([]byte(data), &Config); err != nil {
		return err
	}
	return nil
}
