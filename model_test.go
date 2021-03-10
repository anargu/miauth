package miauth_test

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	miauthv2 "github.com/anargu/miauth"
	"github.com/anargu/miauth/server"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/gorm"
	"github.com/k0kubun/pp"
)

var (
	DB *gorm.DB
	//t1, t2, t3, t4, t5 time.Time
)

var withDB *bool

func TestMain(m *testing.M) {
	withDB = flag.Bool("with-db", false, "tests made connecting with db")
	flag.Parse()
	fmt.Printf(">> flag with-db: %v\n", *withDB)
	Init()
	os.Exit(m.Run())
}

func Init() {
	var err error

	if *withDB {
		LoadConfig()
		miauthv2.Config.DB.Postgres = "user=miauth password=miauth DB.name=miauth port=9910 sslmode=disable"
		if DB, err = OpenTestConnection(); err != nil {
			panic(fmt.Sprintf("No error should happen when connecting to test database, but got err=%+v", err))
		}

		runTestMigration()
	}

}

func OpenTestConnection() (*gorm.DB, error) {
	fmt.Println("connecting to db...")
	miauthv2.InitDB()
	return miauthv2.DB, nil
}

var tables = []interface{}{
	miauthv2.User{},
	miauthv2.Session{},
	miauthv2.Role{},
	miauthv2.LoginCredential{},
	miauthv2.MiauthLoginCredential{},
	miauthv2.GoogleLoginCredential{},
	miauthv2.FacebookLoginCredential{},
	miauthv2.UserCounter{},
}

func runTestMigration() {
	if err := DB.DropTableIfExists(tables...).Error; err != nil {
		panic(fmt.Sprintf("got error when droping tables %+v", err))
	}
	if err := DB.AutoMigrate(tables...).Error; err != nil {
		panic(fmt.Sprintf("No error should happen when create table, but got %+v", err))
	}
}

func TestHasCreatedTables(t *testing.T) {
	for _, table := range tables {
		if hasTable := DB.HasTable(table); !hasTable {
			t.Error("After migration, DB has no USER table")
		}
	}
}

func TestCreateUser(t *testing.T) {
	var user *miauthv2.User
	user = &miauthv2.User{
		Email:    "joe@mail.com",
		Username: "joe1234",
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Error(err)
	}
	if _, err := pp.Printf("::: USER CREATED :::\n%+v\n", user); err != nil {
		panic(err)
	}
}

func TestEmptyLoginCredential(t *testing.T) {
	user := &miauthv2.User{
		Email:    "joe1@mail.com",
		Username: "joe1",
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Error(err)
	}
	if user == nil {
		t.Error("user not created")
		return
	}
	_, err := user.FindCredentialType(miauthv2.MiauthLC)
	if err != nil {
		if err.Error() == "no LoginCredential found" {
			pp.Printf("Test throws correct error: %s\n", err.Error())
		}
	} else {
		t.Error("it should throw Error")
	}
}

func TestMiauthLoginCredential(t *testing.T) {
	user := &miauthv2.User{
		Email:    "joe2@mail.com",
		Username: "joe2",
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Error(err)
	}
	if user == nil {
		t.Error("user not created")
		return
	}
	mlc := &miauthv2.MiauthLoginCredential{
		Hash: "passwordOfJoe",
	}
	if err := DB.Create(&mlc).Error; err != nil {
		t.Error(err)
	}
	lc := &miauthv2.LoginCredential{
		KindLoginCredential: miauthv2.MiauthLC,
		LoginCredentialID:   mlc.ID,
		UserID:              user.ID,
	}
	if err := DB.Create(&lc).Error; err != nil {
		t.Error(err)
	}
	mlcRawResult, err := user.FindCredentialType(miauthv2.MiauthLC)
	if err != nil {
		t.Error(err.Error())
	} else {
		mlcResult, _ := (*mlcRawResult).(miauthv2.MiauthLoginCredential)
		if _, err := pp.Printf("%+v\n", mlcResult); err != nil {
			panic(err)
		}
	}
}

func TestNoFacebookLoginCredential(t *testing.T) {
	user := &miauthv2.User{
		Email:    "joe4@mail.com",
		Username: "joe4",
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Error(err)
	}
	if user == nil {
		t.Error("user not created")
		return
	}
	userSameEmail := &miauthv2.User{
		Email:    "joe4@mail.com",
		Username: "joe5",
	}
	if err := DB.Create(&userSameEmail).Error; err != nil {
		_, _ = pp.Printf("::: userSameEmail ::: it throws error: %v\n...it's fine", err)
	} else {
		t.Error("userSameEmail should throw error")
	}
	userSameUsername := &miauthv2.User{
		Email:    "joe5@mail.com",
		Username: "joe4",
	}
	if err := DB.Create(&userSameUsername).Error; err != nil {
		_, _ = pp.Printf("::: userSameUsername ::: it throws error: %v\n...it's fine", err)
	} else {
		t.Error("userSameUsername should throw error")
	}
	userSameBoth := &miauthv2.User{
		Email:    "joe4@mail.com",
		Username: "joe4",
	}
	if err := DB.Create(&userSameBoth).Error; err != nil {
		_, _ = pp.Printf("::: userSameBoth ::: it throws error: %v\n...it's fine", err)
	} else {
		t.Error("userSameBoth should throw error")
	}

}

func TestValidRegexForUsername(t *testing.T) {
	//[a-zA-Z0-9\-\_]
	var validUsername = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
	//var validUsername = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)

	assert.Equal(t, validUsername.MatchString("anargu"), true)
	assert.Equal(t, validUsername.MatchString("34324"), true)
	assert.Equal(t, validUsername.MatchString("343_24"), true)
	assert.Equal(t, validUsername.MatchString("3asdasda-43-24"), true)

	assert.Equal(t, validUsername.MatchString("anar gu"), false)
	assert.Equal(t, validUsername.MatchString(""), false)
	assert.Equal(t, validUsername.MatchString("3asdasda-43-24@"), false)
	assert.Equal(t, validUsername.MatchString("3asdasda-43-24$"), false)
	assert.Equal(t, validUsername.MatchString("3asdasda-43-24%"), false)
}

//func TestFirstAndLast(t *testing.T) {
//	DB.Save(&User{Name: "user1", Emails: []Email{{Email: "user1@example.com"}}})
//	DB.Save(&User{Name: "user2", Emails: []Email{{Email: "user2@example.com"}}})
//
//	var user1, user2, user3, user4 User
//	DB.First(&user1)
//	DB.Order("id").Limit(1).Find(&user2)
//
//	ptrOfUser3 := &user3
//	DB.Last(&ptrOfUser3)
//	DB.Order("id desc").Limit(1).Find(&user4)
//	if user1.Id != user2.Id || user3.Id != user4.Id {
//		t.Errorf("First and Last should by order by primary key")
//	}
//
//	var users []User
//	DB.First(&users)
//	if len(users) != 1 {
//		t.Errorf("Find first record as slice")
//	}
//
//	var user User
//	if DB.Joins("left join emails on emails.user_id = users.id").First(&user).Error != nil {
//		t.Errorf("Should not raise any error when order with Join table")
//	}
//
//	if user.Email != "" {
//		t.Errorf("User's Email should be blank as no one set it")
//	}
//}

func TestUsernameGenerator(t *testing.T) {
	for i := 0; i < 20; i++ {
		username, err := miauthv2.GenerateGenericUsername()
		if err != nil {
			t.Fail()
		}
		if *username != fmt.Sprintf("user_%d", i+1) {
			t.Fail()
		}

	}
}

func TestUsernameRulesAtSignup(t *testing.T) {

	// first initialize config params
	miauthv2.InitConfig()

	testCases := map[string]bool{
		`{ 
				"kind": "miauth",
				"role": "user",
				"email": "juan@abc.com",
				"username": "anargu",
				"credential": {
					"password": "anargu"
				}
		}`: true,

		`{ 
				"kind": "apple",
				"role": "user",
				"email": "juan@abc.com",
				"username": "",
				"credential": {
					"password": "anargu"
				}
		}`: true,

		`{ 
				"kind": "facebook",
				"role": "user",
				"email": "juan@abc.com",
				"username": "",
				"credential": {
					"password": "anargu"
				}
		}`: true,

		`{ 
				"kind": "google",
				"role": "user",
				"email": "juan@abc.com",
				"username": "",
				"credential": {
					"password": "anargu"
				}
		}`: true,

		`{ 
				"kind": "apple",
				"role": "user",
				"email": "juan@abc.com",
				"username": "anargu",
				"credential": {
					"password": "anargu"
				}
		}`: true,

		`{ 
				"kind": "miauth",
				"role": "user",
				"email": "juan@abc.com",
				"username": "",
				"credential": {
					"password": "anargu"
				}
		}`: false,
	}
	for payload, okResponse := range testCases {

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(payload))
		var input server.SignupInputPayload
		err := c.ShouldBindJSON(&input)
		if okResponse {
			assert.Equal(t, nil, err)
			assert.Equal(t, 200, w.Code)
		} else {
			assert.NotEqual(t, nil, err)
			assert.Equal(t, 200, w.Code)
		}
	}

}

func TestUsernameGenerated4ThirdParties(t *testing.T) {
	// type UserInput map[string]string
	// userCases := []UserInput{
	// 	{"username": "", "email": "abc@maggie.com", "facebookId": "12345"},
	// 	{"username": "", "email": "abc@maggie.com", "facebookId": "12345"},
	// }

	// for i, userCase := range userCases {
	// 	t.Logf("test case %v", i)

	// }

}
