package tests

import (
	"Server/models"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRegistration(t *testing.T) {
	// clear up col before the test
	cleanupCollections()

	tests := []struct {
		name           string
		payload        models.CreateUser
		expectedStatus int
		shouldContain  []string
	}{
		{
			name: "Valid Registration",
			payload: models.CreateUser{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedStatus: 200,
			shouldContain:  []string{"result", "token"},
		},
		{
			name: "Missing Reqired Fields",
			payload: models.CreateUser{
				Email: "missing@example.com",
			},
			expectedStatus: 400,
		},
		{
			name: "Duplicate Email Registration",
			payload: models.CreateUser{
				Email:     "duplicate@example.com",
				Password:  "password456",
				FirstName: "jane",
				LastName:  "smith",
			},
			expectedStatus: 502,
			shouldContain:  []string{"Alraedy Exist"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// for duplicate test we need first to register the original user
			if tt.name == "Duplicate Email Registration" {
				firstpayload := models.CreateUser{
					Email:     "duplicate@example.com",
					Password:  "passwo2226",
					FirstName: "ahmed",
					LastName:  "khalaf",
				}
				registerUser(t, firstpayload, 200)
			}

			// makeing the request
			jsonPaylod, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/user/signup", bytes.NewBuffer(jsonPaylod))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			// chek status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// parese res
			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			if err != nil {
				t.Logf("could not parse response as json: %v", err)
				return
			}

			// convert to string for contains check
			responseStr, _ := json.Marshal(responseBody)
			// check exppected content
			for _, contain := range tt.shouldContain {
				assert.Contains(t, string(responseStr), contain)
			}

			// For susccfully registrain verify basic structer
			if tt.expectedStatus == 200 {
				if respMap, ok := responseBody.(map[string]interface{}); ok {
					assert.Contains(t, respMap, "token")
					assert.Contains(t, respMap, "result")

					if result, ok := respMap["result"].(map[string]interface{}); ok {
						assert.Equal(t, tt.payload.Email, result["email"])
						expactedName := tt.payload.FirstName + " " + tt.payload.LastName
						assert.Equal(t, expactedName, result["name"])
					}
				}
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	cleanupCollections()

	// regitsr a user first
	registerPayload := models.CreateUser{
		Email:     "login@example.com",
		Password:  "password123",
		FirstName: "Login",
		LastName:  "Test",
	}

	registerUser(t, registerPayload, 200)

	tests := []struct {
		name           string
		payload        models.LoginUser
		expectedStatus int
		shouldContain  []string
	}{
		{
			name: "Valid Login",
			payload: models.LoginUser{
				Email:    "login@example.com",
				Password: "password123",
			},
			expectedStatus: 200,
			shouldContain:  []string{"result", "token"},
		},
		{
			name: "Invalid Email",
			payload: models.LoginUser{
				Email:    "noneexistent@example.com",
				Password: "password123",
			},
			expectedStatus: 502,
		},
		{
			name: "Invalid Password",
			payload: models.LoginUser{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: 502,
		},
		{
			name: "Empty",
			payload: models.LoginUser{
				Email:    "",
				Password: "",
			},
			expectedStatus: 400,
		},
	}
	// loop for cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// make the req
			// makeing the request
			jsonPaylod, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/user/signin", bytes.NewBuffer(jsonPaylod))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()
			// check statsu code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// parese res
			var responseBody interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			if err != nil {
				t.Logf("could not parse response as json: %v", err)
				return
			}

			// convert to string for contains check
			responseStr, _ := json.Marshal(responseBody)
			// check exppected content
			for _, contain := range tt.shouldContain {
				assert.Contains(t, string(responseStr), contain)
			}

			// for succ login verify base strucre
			if tt.expectedStatus == 200 {
				if respMap, ok := responseBody.(map[string]interface{}); ok {
					assert.Contains(t, respMap, "token")
					assert.Contains(t, respMap, "result")

					if result, ok := respMap["result"].(map[string]interface{}); ok {
						assert.Equal(t, tt.payload.Email, result["email"])
					}
				}
			}

		})
	}
}

// helper func to create user for testing
func registerUser(t *testing.T, payload models.CreateUser, expeactedStatus int) {
	jsonPaylod, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/user/signup", bytes.NewBuffer(jsonPaylod))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, expeactedStatus, resp.StatusCode)

}
