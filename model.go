package miauth

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
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
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.NewV4()
	return scope.SetColumn("ID", id)
}

type User struct {
	Base
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Role Role `gorm:"not null" json:"role"`

	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`

	Credentials []LoginCredential `gorm:"foreignkey:UserID" json:"credentials"`
	Sessions    []Session         `gorm:"foreignkey:UserID" json:"sessions"`
}

const (
	MiauthLC = iota + 1
	FacebookLC
	GoogleLC
	AppleLC
)

type Role struct {
	Base
	Name string `gorm:"unique;not null"`
}

type LoginCredential struct {
	Base
	KindLoginCredential int       `gorm:"not null" json:"kind_login_credential"`
	LoginCredentialID   uuid.UUID `gorm:"not null" json:"login_credential_id"`
	UserID              uuid.UUID `gorm:"not null" json:"user_id"`
}

type KindLoginCredential interface {
	Kind() int
}

type AppleLoginCredential struct {
	Base
	AccountID string `gorm:"unique;not null" valid:"required" json:"account_id"`
}

type FacebookLoginCredential struct {
	Base
	AccountID string `gorm:"unique;not null" valid:"required" json:"account_id"`
}

type GoogleLoginCredential struct {
	Base
	AccountID string `gorm:"unique;not null" valid:"required" json:"account_id"`
}

type MiauthLoginCredential struct {
	Base
	Hash string `gorm:"not null" valid:"max=52,min=4" json:"hash"`
}

func (flc FacebookLoginCredential) Kind() int {
	return FacebookLC
}
func (glc GoogleLoginCredential) Kind() int {
	return GoogleLC
}
func (alc AppleLoginCredential) Kind() int {
	return AppleLC
}
func (mlc MiauthLoginCredential) Kind() int {
	return MiauthLC
}

type UserCounter struct {
	ID uint64 `gorm:"primaryKey"`
}

type Session struct {
	Base
	AccessToken  string    `gorm:"not null" json:"access_token"`
	RefreshToken string    `gorm:"not null" json:"refresh_token"`
	Scope        *string   `json:"scope"`
	ExpiresIn    string    `gorm:"not null" json:"expires_in"`
	UserID       uuid.UUID `json:"user_id"`
}

func (user *User) FindCredentialType(kind int) (*KindLoginCredential, error) {
	if err := DB.Preload("Credentials").Find(&user).Error; err != nil {
		return nil, err
	}

	var loginCredential KindLoginCredential

	var miauthLoginCredential MiauthLoginCredential
	var facebookLoginCredential FacebookLoginCredential
	var googleLoginCredential GoogleLoginCredential
	var appleLoginCredential AppleLoginCredential

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
			case AppleLC:
				DB.Where(&AppleLoginCredential{
					Base: Base{ID: credential.LoginCredentialID},
				}).First(&appleLoginCredential)
				loginCredential = appleLoginCredential
				return &loginCredential, nil
			}
			break
		}
	}

	return nil, errors.New("no LoginCredential found")
}

func GenerateGenericUsername() (*string, error) {
	uc := UserCounter{}
	result := DB.Create(&uc)
	if result.Error != nil{
		return nil, result.Error
	}
	newUsername := fmt.Sprintf("%s_%d", "user", uc.ID)
	return &newUsername, nil
}
