package groups

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Praveenkusuluri08/bootstrap"
	"github.com/gin-gonic/gin"
)

func TestCreateGroup(t *testing.T) {
	//test case to check the db is working correctly
	bootstrap.DBConnect()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	group := &Group{
		GroupName: "test_group",
		Users:     []map[string]string{{"email": "test@example.com"}},
		Type:      "test_type",
	}
	groupJSON, _ := json.Marshal(group)
	c.Request = httptest.NewRequest("POST", "/api/v1/groups/creategroup", strings.NewReader(string(groupJSON)))
	handler := CreateGroup()
	handler(c)
	if c.Writer.Status() != http.StatusCreated {
		t.Errorf("Expected status %d, but got %d", http.StatusCreated, c.Writer.Status())
	}
	// test case the invalid body
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	groupJson, _ := json.Marshal(map[string]string{"invalid": "body"})
	c.Request = httptest.NewRequest("POST", "/api/v1/groups/creategroup", strings.NewReader(string(groupJson)))
	handler = CreateGroup()
	handler(c)
	if c.Writer.Status() != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, c.Writer.Status())
	}

	//test case to check the group_name is exists
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	groupJson_group_exists, _ := json.Marshal(group)

	c.Request = httptest.NewRequest("POST", "/api/v1/groups/creategroup", strings.NewReader(string(groupJson_group_exists)))
	handler = CreateGroup()
	handler(c)
	if c.Writer.Status() != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, c.Writer.Status())
	}

}
