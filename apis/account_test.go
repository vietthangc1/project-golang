package apis

import (
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
	r := require.New(t)
	account := createRandomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		GetAccountByID(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/account/get/%d", account.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	r.NoError(err)

	server.router.ServeHTTP(recorder, req)
	r.Equal(http.StatusOK, recorder.Code)

	data, err := io.ReadAll(recorder.Body)
	r.NoError(err)

	var actualAccount db.Account
	err = json.Unmarshal(data, &actualAccount)
	r.NoError(err)

	r.Equal(account, actualAccount)
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
