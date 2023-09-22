package transactions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jorgepiresg/ChallangePismo/app"
	mocksApp "github.com/jorgepiresg/ChallangePismo/mocks/app"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Run("register group", func(t *testing.T) {
		Register(echo.New().Group(""), app.App{})
	})
}

func TestMake(t *testing.T) {

	type fields struct {
		transactions *mocksApp.MockITransactions
	}

	tests := map[string]struct {
		input    string
		expected int
		err      error
		prepare  func(f *fields)
	}{
		"success: status 201 created": {
			input: `{"account_id":"id", "operation_type_id": 1, "amount": 1}`,
			prepare: func(f *fields) {
				f.transactions.EXPECT().Make(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			expected: 201,
		},
		"error: status 400 payload invalid": {
			input: `{"account_id": 123}`,
			prepare: func(f *fields) {
			},
			err: fmt.Errorf("any error"),
		},
		"error: status 400 error transaction": {
			input: `{"account_id":"id", "operation_type_id": 1, "amount": 1}`,
			prepare: func(f *fields) {
				f.transactions.EXPECT().Make(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("any"))
			},
			err: fmt.Errorf("any"),
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			transactionsMock := mocksApp.NewMockITransactions(ctrl)

			tt.prepare(&fields{
				transactions: transactionsMock,
			})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.input))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &handler{
				app: app.App{
					Transactions: transactionsMock,
				},
			}

			if tt.err == nil && assert.NoError(t, h.make(c)) {
				assert.Equal(t, tt.expected, rec.Code)
			}

			if tt.err != nil && !assert.Error(t, h.make(c)) {
				t.Errorf(`Expected err: "%s"`, tt.err)
			}
		})
	}
}
