package redirect_test

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/mocks"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "http://url-shortener.com",
		},
		{
			name:      "Not Found",
			alias:     "blabla",
			url:       "http://qwerty.ru",
			mockError: storage.ErrURLNotFound,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			SQLServiceMock := mocks.NewSQLService(t)
			if tc.respError == "" || tc.mockError != nil {
				SQLServiceMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError)
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), SQLServiceMock))
			testServer := httptest.NewServer(r)
			defer testServer.Close()

			redirectedToUrl, err := api.GetRedirect(testServer.URL + "/" + tc.alias)
			if tc.mockError != nil && errors.Is(tc.mockError, storage.ErrURLNotFound) {
				require.Equal(t, api.ErrNotFound, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.url, redirectedToUrl)
			}
		})

	}
}
