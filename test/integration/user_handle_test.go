package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"microblog/domain/user/domain"
	testDB "microblog/infrastructure/database/test"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type Map map[string]interface{}

func dataUSer() []domain.User {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	return []domain.User{
		{
			ID:        uint(1),
			FirstName: "Daniel",
			LastName:  "De La Pava Suarez",
			Username:  "daniel.delapava",
			Email:     "daniel.delapava@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uint(1),
			FirstName: "Rebecca",
			LastName:  "Romero",
			Username:  "rebecca.romero",
			Email:     "rebecca.romero@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func TestIntegration_GetAllUser(t *testing.T) {

	t.Run("No Content (no seed data)", func(tt *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)

		if e, a := http.StatusOK, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var users []domain.User
		var valueResponse = Map{"user": users}
		if err := json.Unmarshal([]byte(w.Body.String()), &valueResponse); err != nil {
			tt.Errorf("error decoding response body: %v", err)
		}

		if len(users) > 0 {
			tt.Errorf("expected no lists to be returned, got %v lists", len(users))
		}
	})

	t.Run("Ok (database has been seeded)", func(tt *testing.T) {
		defer func() {
			if err := testDB.Truncate(d.DB); err != nil {
				tt.Errorf("error truncating test database tables: %v", err)
			}
		}()

		expectedUsers, err := testDB.SeedUsers(d.DB)
		if err != nil {
			tt.Fatalf("error seeding users: %v", err)
		}

		req, err := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)

		if e, a := http.StatusOK, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var users []domain.User
		if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
			tt.Errorf("error decoding response body: %v", err)
		}

		if d := cmp.Diff(expectedUsers[0].ID, users[0].ID); d != "" {
			tt.Errorf("unexpected difference in response body:\n%v", d)
		}
	})
}

func TestIntegration_GetOneHandler(t *testing.T) {
	defer func() {
		if err := testDB.Truncate(d.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedLists, err := testDB.SeedUsers(d.DB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
	}

	tests := []struct {
		Name         string
		UserID       uint
		ExpectedBody domain.User
		ExpectedCode int
	}{
		{
			Name:         "Get One User Successful",
			UserID:       expectedLists[0].ID,
			ExpectedBody: expectedLists[0],
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "User Not Found",
			UserID:       0,
			ExpectedBody: domain.User{},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d", test.UserID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			s.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusNotFound {
				var response domain.User

				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if e, a := test.ExpectedBody.ID, response.ID; e != a {
					t.Errorf("expected user ID: %v, got user ID: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func TestIntegration_CreateHandler(t *testing.T) {

	defer func() {
		if err := testDB.Truncate(d.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	tests := []struct {
		Name         string
		RequestBody  domain.User
		ExpectedCode int
	}{
		{
			Name:         "Create User Successful",
			RequestBody:  dataUSer()[0],
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:         "Break Unique UserName Constraint",
			RequestBody:  dataUSer()[0],
			ExpectedCode: http.StatusConflict,
		},
		{
			Name:         "No Data User",
			RequestBody:  domain.User{},
			ExpectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users/", &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			w := httptest.NewRecorder()
			s.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusConflict {
				var response domain.User

				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				if e, a := test.RequestBody.Username, response.Username; e != a {
					t.Errorf("expected user UserName: %v, got user UserName: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func TestIntegration_UpdateHandler(t *testing.T) {

	defer func() {
		if err := testDB.Truncate(d.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedUsers, err := testDB.SeedUsers(d.DB)
	if err != nil {
		t.Fatalf("error seeding users: %v", err)
	}

	tests := []struct {
		Name         string
		UserID       uint
		RequestBody  domain.User
		ExpectedCode int
	}{
		{
			Name:         "Update User Successful",
			UserID:       expectedUsers[0].ID,
			RequestBody:  expectedUsers[0],
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "Break Unique UserName Constraint",
			UserID:       expectedUsers[1].ID,
			RequestBody:  expectedUsers[0],
			ExpectedCode: http.StatusConflict,
		},
		{
			Name:         "No Data User",
			UserID:       expectedUsers[0].ID,
			RequestBody:  domain.User{},
			ExpectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%d", test.UserID), &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			w := httptest.NewRecorder()
			s.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}
		}

		t.Run(test.Name, fn)
	}
}

func TestIntegration_DeleteHandler(t *testing.T) {
	defer func() {
		if err := testDB.Truncate(d.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedLists, err := testDB.SeedUsers(d.DB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
	}

	tests := []struct {
		Name         string
		UserID       uint
		ExpectedCode int
	}{
		{
			Name:         "Delete User Successful",
			UserID:       expectedLists[0].ID,
			ExpectedCode: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", test.UserID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			s.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}
		}

		t.Run(test.Name, fn)
	}
}
