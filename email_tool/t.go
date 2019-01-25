package main

import (
	"flag"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os"
	"strings"
)

func getFileEmails(message *gomail.Message, addrs string) ([]string, error) {
	addrs = strings.Trim(addrs, " \t")

	data, err := ioutil.ReadFile(addrs)
	if err != nil {
		return nil, err
	}

	return parseEmails(message, string(data), "\n"), nil
}

func parseEmails(message *gomail.Message, addrs string, split string) []string {
	addrs = strings.Trim(addrs, " \t")
	var format []string
	ens := strings.Split(addrs, split)
	for _, it := range ens {
		en := strings.Split(it, ":")
		var email, name string
		if len(en) >= 1 {
			email = en[0]
		}

		if len(en) >= 2 {
			name = en[1]
		}

		email = strings.Trim(email, " \t")
		name = strings.Trim(name, " \t")
		if len(email) > 0 {
			if email[0] == '#' {
				continue
			}
			format = append(format, message.FormatAddress(email, name))
		}
	}

	return format

}

func getEmails(message *gomail.Message, addrs string) ([]string, error) {
	addrs = strings.Trim(addrs, " \t")

	if len(addrs) > 0 && addrs[0] == '@' {
		return getFileEmails(message, addrs[1:])
	}

	return parseEmails(message, addrs, "-"), nil
}

func main() {
	// smtp
	var smtp, from, passwd string
	var port int
	var subject string
	// 发送者
	var sender string
	// 附件
	var attach string
	// body mime, default
	// text/html
	var mime string

	var to string
	var cc string

	flag.StringVar(&smtp, "smtp", "", "smtp address")
	flag.IntVar(&port, "port", 0, "smtp port")
	flag.StringVar(&from, "from", "", "from email")
	flag.StringVar(&passwd, "passwd", "", "from email password")

	flag.StringVar(&subject, "subject", "", "email subject")
	flag.StringVar(&sender, "sender", "", "sender name")
	flag.StringVar(&attach, "attach", "", "attach if have")

	flag.StringVar(&to, "to", "", "to email format:address:name-address:name split by '-'(example@example.com:name-example@example.com:name...). If first chat is '@' email will read the file")
	flag.StringVar(&cc, "cc", "", "cc email as same as 'to'")

	flag.StringVar(&mime, "mime", "text/html", "body mime")
	flag.Parse()

	if len(smtp) == 0 {
		fmt.Println("smtp flag can not find")
		return
	}

	if port == 0 {
		fmt.Println("port flag can not find")
		return
	}

	if len(from) == 0 {
		fmt.Println("from flag can not find")
		return
	}

	if len(passwd) == 0 {
		fmt.Println("passwd flag can not find")
		return
	}

	if len(subject) == 0 {
		fmt.Println("subject flag can not find")
		return
	}

	if len(sender) == 0 {
		fmt.Println("sender flag can not find")
		return
	}

	if len(to) == 0 {
		fmt.Println("to flag can not find")
		return
	}

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("read email body err:%s", err)
		return
	}

	m := gomail.NewMessage()
	// from 必须和NewDialer 一致，否则就报错
	// panic: gomail: could not send email 1: 501 mail from address must be same as authorization user
	m.SetAddressHeader("From", from, sender)
	//m.SetHeader("To", "fff <fengguangxiang@example.com>", "fengguangxiang@example.com")
	toemails, err := getEmails(m, to)
	if err != nil {
		fmt.Println("get to emails err:", err)
		return
	}

	m.SetHeader("To", toemails...)
	if len(cc) > 0 {
		ccemails, err := getEmails(m, cc)
		if err != nil {
			fmt.Println("get cc emails err:", err)
			return
		}

		m.SetHeader("Cc", ccemails...)
	}
	//m.SetAddressHeader("Cc", "fengguangxiang@example.com", "DDD")

	m.SetHeader("Subject", subject)
	m.SetBody(mime, string(body))

	if len(attach) > 0 {
		m.Attach(attach)
	}

	d := gomail.NewDialer(smtp, port, from, passwd)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("send email err:", err)
	}

}
