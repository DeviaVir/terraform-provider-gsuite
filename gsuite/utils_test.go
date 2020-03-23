package gsuite

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {

	testCases := []struct {
		email   string
		success bool
	}{
		{"myemail@domain.com", true},
		{"Alice <alice@example.com>", false},
		{"much-much-much-much-much-much-much-much-much-much-much-much-too-long@domain.com", false},
		{"\"some@much-much-much-much-much-much-much-much-much-much-much-much-too-long\"@domain.com", false},
	}

	for _, testCase := range testCases {
		_, errs := validateEmail(testCase.email, "")
		if len(errs) > 0 && testCase.success {
			t.Log(errs)
			t.Errorf("expected a valid email for %s", testCase.email)
		} else if len(errs) == 0 && !testCase.success {
			t.Errorf("expected an invalid email for %s", testCase.email)
		}
	}
}
