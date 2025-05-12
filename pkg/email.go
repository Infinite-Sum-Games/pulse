package pkg

import (
	"fmt"

	
	cmd "github.com/IAmRiteshKoushik/pulse/cmd"
	"gopkg.in/gomail.v2"
)

func SendOTPEmail(toEmail, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "gmharish285@gmail.com") // Your email
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your OTP Code for ASoC registration")
	m.SetBody("text/plain", fmt.Sprintf("Your OTP code is: %s\nIt is valid for 5 minutes.", otp))
	m.AddAlternative("text/html", fmt.Sprintf("<strong>Your OTP code is: %s</strong><br>It is valid for 5 minutes.", otp))

	d := gomail.NewDialer("smtp.gmail.com", 587, "gmharish285@gmail.com", cmd.EnvVars.EmailAppPassword)

	return d.DialAndSend(m)
}
