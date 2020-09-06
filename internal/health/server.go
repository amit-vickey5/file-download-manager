package health

import (
	"encoding/json"
	"github.com/amit/file-download-manager/pkg/db"
	"github.com/amit/file-download-manager/pkg/logger"
	"net/http"
	"strconv"
)


type StatusHandler struct {}

func (s *StatusHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.LogStatement("health check request...", nil)
	ok, err := db.DatabaseHealthCheck()
	responseContent := make(map[string]string)
	if ok && err == nil {
		logger.LogStatement("status check :: server up and running", nil)
		responseContent["status"] = "file-download-manager server is up and running"
		resp.WriteHeader(http.StatusOK)
	} else {
		logger.LogStatement("STATUS CHECK ERROR ::", err)
		responseContent["status"] = "file-download-manager server is down"
		resp.WriteHeader(503)
	}
	respBytes, _ := json.Marshal(responseContent)
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.Write(respBytes)
}
