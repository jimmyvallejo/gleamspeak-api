package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateUser(ctx context.Context, params database.CreateUserParams) (database.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockDB := new(MockDB)
	h := &handlers.Handlers{DB: mockDB}

	reqBody := `{"email": "test@test.com", "handle": "testuser"}`
	req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	now := time.Now().UTC().Truncate(time.Millisecond)
	
	expectedUser := database.User{
		ID:        uuid.New(),
		Email:     "test@test.com",
		Handle:    "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockDB.On("CreateUser", mock.Anything, mock.MatchedBy(func(params database.CreateUserParams) bool {
		return params.Email == "test@test.com" && params.Handle == "testuser"
	})).Return(expectedUser, nil)

	h.CreateUser(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response database.User
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Email, response.Email)
	assert.Equal(t, expectedUser.Handle, response.Handle)

	mockDB.AssertExpectations(t)
}
