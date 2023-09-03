package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vietthangc1/simple_bank/db/mock"
	db "github.com/vietthangc1/simple_bank/db/sqlc"
)

var (
	wordLength        = 6
	passwordLength    = 36
	emailLength       = 10
	emailDomainLength = 5
)

type createUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e createUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := passwordManager.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e createUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v, password %s", e.arg, e.password)
}

func eqCreateUser(arg db.CreateUserParams, password string) gomock.Matcher {
	return createUserParamsMatcher{
		arg:      arg,
		password: password,
	}
}

func TestCreateUser(t *testing.T) {
	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStore    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, server *Server, recorder *httptest.ResponseRecorder, req *http.Request)
	}{
		{
			name: "StatusOK",
			body: gin.H{
				"user_name": user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStore: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(
						gomock.Any(),
						// gomock.Eq(arg),
						eqCreateUser(arg, password),
					).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, server *Server, recorder *httptest.ResponseRecorder, req *http.Request) {
				server.router.ServeHTTP(recorder, req)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, test := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		test.buildStore(store)

		server := NewServer(store)
		recorder := httptest.NewRecorder()

		data, err := json.Marshal(test.body)
		require.NoError(t, err)

		url := "/user/create"
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(data)))
		require.NoError(t, err)

		test.checkResponse(t, server, recorder, req)
	}
}

func createRandomUser(t *testing.T) (db.User, string) {
	password := GeneratePassword()
	hashedPassword, err := passwordManager.HashPassword(password)
	require.NoError(t, err)

	return db.User{
		Username:       GenerateUsername(),
		FullName:       GenerateFullname(),
		Email:          GenerateEmail(),
		HashedPassword: hashedPassword,
	}, password
}

func GenerateUsername() string {
	return randomEntity.RandomString(wordLength)
}

func GenerateFullname() string {
	return fmt.Sprintf("%s %s %s", randomEntity.RandomString(wordLength), randomEntity.RandomString(wordLength), randomEntity.RandomString(wordLength))
}

func GenerateEmail() string {
	return fmt.Sprintf("%s@%s.com", randomEntity.RandomString(emailLength), randomEntity.RandomString(emailDomainLength))
}

func GeneratePassword() string {
	return randomEntity.RandomString(passwordLength)
}
