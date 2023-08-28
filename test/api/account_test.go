package test_api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-backend-practice/api"
	"github.com/stretchr/testify/require"
)

func Test_GetAccountAPI(t *testing.T) {
	arg := api.CreateAccountParams{
		Owner:    "test",
		Currency: "USD",
	}

	s, _ := json.Marshal(arg);

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/account", bytes.NewReader(s))

	require.NoError(t, err)
	
	server.Serve(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
}