package pkg

import "regexp"

type RegExps struct {
	// Advanced regex
	NameRegex     *regexp.Regexp
	PasswordRegex *regexp.Regexp
	OtpRegex      *regexp.Regexp
}

var Rex = &RegExps{}

func InitRegex() error {

	// Regular expression for FirstName, MiddleName, LastName
	nameRegex, err := regexp.Compile(`^[a-zA-Z]+$`)
	if err != nil {
		return err
	}
	Rex.NameRegex = nameRegex

	// Regular expression for OTP
	otpRegex, err := regexp.Compile(`^[0-9]+$`)
	if err != nil {
		return err
	}
	Rex.OtpRegex = otpRegex

	// Regular express for Passwords
	pwdRegex, err := regexp.Compile(`^[a-zA-Z0-9_.]+$`)
	if err != nil {
		return err
	}
	Rex.PasswordRegex = pwdRegex

	return nil
}
