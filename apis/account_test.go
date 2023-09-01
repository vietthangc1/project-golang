package apis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vietthangc1/simple_bank/db/mock"
	db "github.com/vietthangc1/simple_bank/db/sqlc"
	"github.com/vietthangc1/simple_bank/pkg/randomx"
)

func TestGetAccounByID(t *testing.T) {
	account := createRandomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStore    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, server *Server, recorder *httptest.ResponseRecorder, req *http.Request)
	}{
		{
			name:      "StatusOK",
			accountID: account.ID,
			buildStore: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountByID(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, server *Server, recorder *httptest.ResponseRecorder, req *http.Request) {
				server.router.ServeHTTP(recorder, req)
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var actualAccount db.Account
				err = json.Unmarshal(data, &actualAccount)
				require.NoError(t, err)

				require.Equal(t, account, actualAccount)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStore: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountByID(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, server *Server, recorder *httptest.ResponseRecorder, req *http.Request) {
				server.router.ServeHTTP(recorder, req)
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			buildStore: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, server *Server, recorder *httptest.ResponseRecorder, req *http.Request) {
				server.router.ServeHTTP(recorder, req)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStore: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountByID(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, server *Server, recorder *httptest.ResponseRecorder, req *http.Request) {
				server.router.ServeHTTP(recorder, req)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

		url := fmt.Sprintf("/account/get/%d", test.accountID)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		test.checkResponse(t, server, recorder, req)
	}
}

func createRandomAccount() db.Account {
	Rand := randomx.NewRandom()
	return db.Account{
		ID:       Rand.RandomInt(1, 100),
		Owner:    Rand.RandomString(6),
		Balance:  0,
		Currency: "VND",
	}
}

func TestCreateAccount(t *testing.T) {
	
}