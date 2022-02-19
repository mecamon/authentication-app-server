package auth

import "testing"

var validEmailSetup = []struct {
	testName       string
	email          string
	password       string
	expectedResult bool
}{
	{"Invalid email 1", "notvalid.com", "it-does-not-matter", false},
	{"Invalid email 2", "not@validcom", "it-does-not-matter", false},
	{"valid email", "valid@mail.com", "it-does-not-matter", true},
}

func TestAuth_ValidEmail(t *testing.T) {
	for _, tt := range validEmailSetup {
		t.Log(tt.testName)
		{
			auth := Auth{
				Email:    tt.email,
				Password: tt.password,
			}

			output := auth.ValidEmail()
			if output != tt.expectedResult {
				t.Errorf("Got %v when it should be %v", output, tt.expectedResult)
			}
		}
	}
}

var validPasswordSetup = []struct {
	testName       string
	email          string
	password       string
	expectedResult bool
}{
	{"invalid pass 1", "valid@mail.com", "invalid", false},
	{"invalid pass 2", "valid@mail.com", "L4i", false},
	{"invalid pass 3", "valid@mail.com", "INVALID234", false},
	{"invalid pass 4", "valid@mail.com", "123456789", false},
	{"valid pass", "valid@mail.com", "Validpass12", true},
}

func TestAuth_ValidPassword(t *testing.T) {
	for _, tt := range validPasswordSetup {
		t.Log(tt.testName)
		{
			auth := Auth{
				Email:    tt.email,
				Password: tt.password,
			}

			output := auth.ValidPassword()
			if output != tt.expectedResult {
				t.Errorf("Got %v when it should be %v", output, tt.expectedResult)
			}
		}
	}
}

func TestAuth_EvaluateNewUserCredentials(t *testing.T) {
	t.Log("invalid user credentials")
	{
		auth := Auth{
			Email:    "invalidmail.co",
			Password: "123145",
		}
		errResponse := auth.EvaluateNewUserCredentials()

		if len(errResponse.Message) != 2 {
			t.Error("has not the expected quantity of errors")
		}
	}
}
