package pkg

import "testing"

var validEmailSetup = []struct {
	testName       string
	email          string
	expectedResult bool
}{
	{"Invalid email 1", "notvalid.com", false},
	{"Invalid email 2", "not@validcom", false},
	{"valid email", "valid@mail.com", true},
}

func TestValidEmail(t *testing.T) {
	for _, tt := range validEmailSetup {
		t.Log(tt.testName)
		{
			output := ValidEmail(tt.email)
			if output != tt.expectedResult {
				t.Errorf("Got %v when it should be %v", output, tt.expectedResult)
			}
		}
	}
}

var validPasswordSetup = []struct {
	testName       string
	password       string
	expectedResult bool
}{
	{"invalid pass 1", "invalid", false},
	{"invalid pass 2", "L4i", false},
	{"invalid pass 3", "INVALID234", false},
	{"invalid pass 4", "123456789", false},
	{"valid pass", "Validpass12", true},
}

func TestValidPassword(t *testing.T) {
	for _, tt := range validPasswordSetup {
		t.Log(tt.testName)
		{
			output := ValidPassword(tt.password)
			if output != tt.expectedResult {
				t.Errorf("Got %v when it should be %v", output, tt.expectedResult)
			}
		}
	}
}

var telephoneTests = []struct {
	testName       string
	telephone      string
	expectedResult bool
}{
	{"Phone number lower than 10", "123456789", false},
	{"Phone number with letters", "12345678daf", false},
	{"Valid phone number", "12014326578", true},
}

func TestValidTelephone(t *testing.T) {

	for _, tt := range telephoneTests {
		t.Log(tt.testName)
		{
			output := ValidTelephone(tt.telephone)
			if output != tt.expectedResult {
				t.Errorf("Got %v but expected %v", output, tt.expectedResult)
			}
		}
	}
}
