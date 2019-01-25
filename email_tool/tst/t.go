package main

import (
	"gopkg.in/gomail.v2"
)

func main() {

	m := gomail.NewMessage()
	m.SetAddressHeader("From", "devops@example.com", "AAAA")
	//m.SetHeader("To", "fff <user@example.com>", "user@example.com")
	m.SetHeader("To", m.FormatAddress("user@example.com", "ffffff"), m.FormatAddress("user@example.com", "bbbb"), m.FormatAddress("user@example.com", ""))
	m.SetAddressHeader("Cc", "user@example.com", "DDD")

	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	m.Attach("/data/home/fenggx/qrcode_for_gh_eb9e7425031c_258.jpg")

	d := gomail.NewDialer("smtp.exmail.qq.com", 465, "devops@example.com", "*****")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}
