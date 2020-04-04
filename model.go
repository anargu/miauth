package miauth

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/satori/go.uuid"
	"regexp"
	"time"
)

var OnlyAlphanumericNoSpaceValues = "only_alphanumeric_no_space_values"
var UsernamePatterns map[string]*regexp.Regexp
var RxEmail *regexp.Regexp

func init() {
	UsernamePatterns = make(map[string]*regexp.Regexp, 1)
	onlyAlphanumericNoSpaceValuesRegex := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
	UsernamePatterns[OnlyAlphanumericNoSpaceValues] = onlyAlphanumericNoSpaceValuesRegex

	Email := "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	RxEmail = regexp.MustCompile(Email)
}

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.NewV4()
	return scope.SetColumn("ID", id)
}

type User struct {
	Base
	FirstName string
	LastName  string

	Role Role `gorm:"not null"`

	Username string `gorm:"unique;not null";valid:"max=72,min=6"`
	Email string `gorm:"unique;not null";valid:"required,email,max=72,min=3"`

	Credentials []LoginCredential `gorm:"foreignkey:UserID"`
	Sessions []Session `gorm:"foreignkey:UserID"`
}

const (
	MiauthLC = iota + 1
	FacebookLC
	GoogleLC
)

type Role struct {
	Base
	Name string `gorm:"unique;not null"`
}

type LoginCredential struct {
	Base
	KindLoginCredential int `gorm:"not null"`
	LoginCredentialID uuid.UUID `gorm:"not null"`
	UserID uuid.UUID `gorm:"not null"`
}

type KindLoginCredential interface {
	Kind() int
}

type FacebookLoginCredential struct {
	Base
	AccountID string `valid:"required"`
}

type GoogleLoginCredential struct {
	Base
	AccountID string `valid:"required"`
}

type MiauthLoginCredential struct {
	Base
	Hash     string `gorm:"not null";valid:"max=52,min=4"`
}

func (flc FacebookLoginCredential) Kind() int {
	return FacebookLC
}
func (glc GoogleLoginCredential) Kind() int {
	return GoogleLC
}
func (mlc MiauthLoginCredential) Kind() int {
	return MiauthLC
}


type Session struct {
	Base
	AccessToken string `gorm:"not null"`
	RefreshToken string `gorm:"not null"`
	Scope *string
	ExpiresIn string `gorm:"not null"`
	UserID uuid.UUID
}

func (user *User) FindCredentialType(kind int) (*KindLoginCredential, error)  {
	if err := DB.Preload("Credentials").Find(&user).Error; err != nil {
		return nil, err
	}

	var loginCredential KindLoginCredential

	var miauthLoginCredential MiauthLoginCredential
	var facebookLoginCredential FacebookLoginCredential
	var googleLoginCredential GoogleLoginCredential

	for _, credential := range user.Credentials {
		if credential.KindLoginCredential == kind {
			switch kind {
			case MiauthLC:
				if err := DB.Where(&MiauthLoginCredential{
					Base: Base{ID: credential.LoginCredentialID},
				}).First(&miauthLoginCredential).Error; err != nil {
					return nil, err
				}
				loginCredential = miauthLoginCredential
				return &loginCredential, nil
			case FacebookLC:
				DB.Where(&FacebookLoginCredential{
					Base: Base{ID: credential.LoginCredentialID},
				}).First(&facebookLoginCredential)
				loginCredential = facebookLoginCredential
				return &loginCredential, nil
			case GoogleLC:
				DB.Where(&GoogleLoginCredential{
					Base: Base{ID: credential.LoginCredentialID},
				}).First(&googleLoginCredential)
				loginCredential = googleLoginCredential
				return &loginCredential, nil
			}
			break
		}
	}

	return nil, errors.New("no LoginCredential found")
}
