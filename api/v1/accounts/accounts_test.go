package accounts

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jorgepiresg/ChallangePismo/app"
	mocksApp "github.com/jorgepiresg/ChallangePismo/mocks/app"
	modelAccounts "github.com/jorgepiresg/ChallangePismo/model/accounts"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	t.Run("register group", func(t *testing.T) {
		Register(echo.New().Group(""), app.App{})
	})
}

func TestCreate(t *testing.T) {

	type fields struct {
		accounts *mocksApp.MockIAccounts
	}

	type expected struct {
		Status   int
		Response string
	}

	tests := map[string]struct {
		input    string
		expected expected
		err      error
		prepare  func(f *fields)
	}{
		"should be able to create a new account": {
			input: `{"document_number":"111.111.111-11"}`,
			prepare: func(f *fields) {
				f.accounts.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(modelAccounts.Account{ID: "id", DocumentNumber: "111111111", CreatedAt: time.Now()}, nil)
			},
			expected: expected{
				Status:   201,
				Response: `{"account_id":"id"}`,
			},
		},
		"should not be able to create a new account with payload invalid": {
			input: `{"document_number":111.111.111-11}`,
			prepare: func(f *fields) {
			},
			err: fmt.Errorf("any error"),
		},
		"should not be able to create a new account with error in app.create": {
			input: `{"document_number":"111.111.111-11"}`,
			prepare: func(f *fields) {
				f.accounts.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(modelAccounts.Account{}, fmt.Errorf("any"))
			},
			err: fmt.Errorf("any"),
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			accountsMock := mocksApp.NewMockIAccounts(ctrl)

			tt.prepare(&fields{
				accounts: accountsMock,
			})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.input))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &handler{
				app: app.App{
					Accounts: accountsMock,
				},
			}

			if tt.err == nil && assert.NoError(t, h.create(c)) {
				assert.Equal(t, tt.expected.Status, rec.Code)
				assert.Equal(t, tt.expected.Response+"\n", rec.Body.String())
			}

			if tt.err != nil && !assert.Error(t, h.create(c)) {
				t.Errorf(`Expected err: "%s"`, tt.err)
			}
		})
	}
}

func TestGetByAccountID(t *testing.T) {

	type fields struct {
		accounts *mocksApp.MockIAccounts
	}

	type expected struct {
		Status   int
		Response string
	}

	tests := map[string]struct {
		input    string
		expected expected
		err      error
		prepare  func(f *fields)
	}{
		"success: status 200": {
			input: `id`,
			prepare: func(f *fields) {
				f.accounts.EXPECT().GetByAccountID(gomock.Any(), "id").Times(1).Return(modelAccounts.Account{ID: "id", DocumentNumber: "11111111111", CreatedAt: time.Now()}, nil)
			},
			expected: expected{
				Status:   200,
				Response: `{"account_id":"id","document_number":"11111111111"}`,
			},
		},
		"error: status 400 error any": {
			input: `invalid_id`,
			prepare: func(f *fields) {
				f.accounts.EXPECT().GetByAccountID(gomock.Any(), "invalid_id").Times(1).Return(modelAccounts.Account{}, fmt.Errorf("any"))
			},
			err: fmt.Errorf("any"),
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			accountsMock := mocksApp.NewMockIAccounts(ctrl)

			tt.prepare(&fields{
				accounts: accountsMock,
			})

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/accounts/:account_id")
			c.SetParamNames("account_id")
			c.SetParamValues(tt.input)

			h := &handler{
				app: app.App{
					Accounts: accountsMock,
				},
			}

			if tt.err == nil && assert.NoError(t, h.getByAccountID(c)) {
				assert.Equal(t, tt.expected.Status, rec.Code)
				assert.Equal(t, tt.expected.Response+"\n", rec.Body.String())
			}

			if tt.err != nil && !assert.Error(t, h.getByAccountID(c)) {
				t.Errorf(`Expected err: "%s"`, tt.err)
			}
		})
	}
}
