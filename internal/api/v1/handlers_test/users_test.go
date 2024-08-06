package handlers_test

// import (
//     "bytes"
//     "context"
//     "encoding/json"
//     "net/http"
//     "net/http/httptest"
//     "testing"
//     "time"

//     "github.com/google/uuid"
//     "github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
//     "github.com/jimmyvallejo/gleamspeak-api/internal/database"
//     "github.com/stretchr/testify/assert"
//     "github.com/stretchr/testify/mock"
//     "database/sql"
// )

// type MockDB struct {
//     mock.Mock
// }

// func (m *MockDB) CreateUserStandard(ctx context.Context, params database.CreateUserStandardParams) (database.User, error) {
//     args := m.Called(ctx, params)
//     return args.Get(0).(database.User), args.Error(1)
// }

// func TestCreateUserStandard(t *testing.T) {
//     mockDB := new(MockDB)
//     h := &handlers.Handlers{DB: mockDB}

//     reqBody := `{"email": "test@test.com", "handle": "testuser", "password": "testpassword"}`
//     req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(reqBody))
//     assert.NoError(t, err)

//     rr := httptest.NewRecorder()

//     now := time.Now().UTC().Truncate(time.Millisecond)
    
//     expectedUser := database.User{
//         ID:        uuid.New(),
//         Email:     "test@test.com",
//         Handle:    "testuser",
//         Password:  sql.NullString{String: "hashedpassword", Valid: true},
//         CreatedAt: now,
//         UpdatedAt: now,
//     }

//     mockDB.On("CreateUserStandard", mock.Anything, mock.MatchedBy(func(params database.CreateUserStandardParams) bool {
//         return params.Email == "test@test.com" && 
//                params.Handle == "testuser" && 
//                params.Password.Valid && 
//                len(params.Password.String) > 0
//     })).Return(expectedUser, nil)

//     h.CreateUserStandard(rr, req)

//     assert.Equal(t, http.StatusCreated, rr.Code)

//     var response handlers.CreateUserResponse
//     err = json.Unmarshal(rr.Body.Bytes(), &response)
//     assert.NoError(t, err)

//     assert.Equal(t, expectedUser.ID, response.ID)
//     assert.Equal(t, expectedUser.Email, response.Email)
//     assert.Equal(t, expectedUser.Handle, response.Handle)

//     mockDB.AssertExpectations(t)
// }