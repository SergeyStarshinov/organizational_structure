package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"orgstructure/internal/config"
	"orgstructure/internal/domain/model"
	"orgstructure/internal/logger"
	"orgstructure/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	body   string
	status int
}

func TestCreateDepartment(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Address: "localhost:8080",
		},
		Database: config.DBConfig{
			User:     "hitalent",
			Password: "hitalent",
			Host:     "localhost",
			Port:     5432,
			DBName:   "hitalent",
		},
		Logger: config.LogConfig{
			LogFile:  "stdout",
			LogLevel: "debug",
		},
	}
	log, logFile := logger.New(cfg.Logger)
	defer logFile.Close()
	repo, err := storage.New(cfg)
	require.NoError(t, err)

	handler := NewBaseHandler(repo, log)

	testCases := []testCase{
		{"{\"name\":\"test1\"}", http.StatusOK},
		{"{\"name\":\"test2\", \"parent_id\":2}", http.StatusBadRequest},
		{"{\"name\":\"test2\", \"parent_id\":\"2\"}", http.StatusBadRequest},
		{"{\"name\":\"test2\", \"parent_id\":0}", http.StatusOK},
		{"{\"name\":\"test2\", \"parent_id\":0}", http.StatusConflict},
	}

	for _, test := range testCases {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/departments/", bytes.NewReader([]byte(test.body)))
		require.NoError(t, err)
		handler.CreateDepartment(rr, req)
		assert.Equal(t, rr.Code, test.status)
	}

	repo.DB.Where("name LIKE ?", "test%").Delete(&model.Department{})

}
