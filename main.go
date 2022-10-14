package main

import (
	_ "embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/wneessen/go-mail"
	"log"
	"os"
)

var (
	//go:embed config.toml
	config string
)

type message struct {
	Username     string
	Password     string
	Title        string
	Host         string
	Port         uint16
	FromName     string
	FromEmail    string
	ToName       string
	ToEmail      string
	ReplyToName  string
	ReplyToEmail string
	Subject      string
	Message      string
}

func main() {
	o := &message{}
	//o := &order{}
	fmt.Println(config)

	err := toml.Unmarshal([]byte(config), o)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", o)

	// Create a new mail message
	m := mail.NewMsg()

	// To set address header fields like "From", "To", "Cc" or "Bcc" you have different methods
	// at your hands. Some perform input validation, some ignore invalid addresses. Some perform
	// the formatting for you.

	if err := m.ReplyToFormat(o.ReplyToName, o.ReplyToEmail); err != nil {
		fmt.Printf("failed to set REPLY-TO address: %s\n", err)
		os.Exit(1)
	}
	if err := m.FromFormat(o.FromName, o.FromEmail); err != nil {
		fmt.Printf("failed to set FROM address: %s\n", err)
		os.Exit(1)
	}
	if err := m.AddToFormat(o.ToName, o.ToEmail); err != nil {
		fmt.Printf("failed to set TO address: %s\n", err)
		os.Exit(1)
	}
	//m.CcIgnoreInvalid("cc@example.com", "invalidaddress+example.com")

	// Set a subject line
	m.Subject(fmt.Sprint(o.Title, " - ", o.Subject))

	// And some other common headers...
	//
	// Sets a valid "Date" header field with the current time
	m.SetDate() // Current time

	// Generates a valid and unique "Message-ID"
	m.SetMessageID()

	// Sets the "Precedence"-Header to "bulk" to indicate a "bulk mail"
	//m.SetBulk()

	// Set a "high" importance to the mail (this sets several Header fields to
	// satisfy the different common mail clients like Mail.app and Outlook)
	m.SetImportance(mail.ImportanceHigh)

	// Add your mail message to body
	m.SetBodyString(mail.TypeTextPlain, o.Message)

	// Attach a file from your local FS
	// We override the attachment fromName using the WithFileName() Option
	//m.AttachFile("/home/tester/test.txt", mail.WithFileName("attachment.txt"))

	// Next let's create a Client
	// We have lots of With* options at our disposal to stear the Client. It will set sane
	// options by default, though

	// Let's assume we need to perform SMTP AUTH with the sending server, though. Since we
	// use SMTP PLAIN AUTH, let's also make sure to enforce strong TLS
	c, err := mail.NewClient(o.Host, mail.WithPort(int(o.Port)),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(o.Username),
		mail.WithPassword(o.Password), mail.WithTLSPolicy(mail.TLSMandatory))
	if err != nil {
		fmt.Printf("failed to create new mail client: %s\n", err)
		os.Exit(1)
	}

	// Now that we have our client, we can connect to the server and send our mail message
	// via the convenient DialAndSend() method. You have the option to Dial() and Send()
	// separately as well
	if err := c.DialAndSend(m); err != nil {
		fmt.Printf("failed to send mail: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Mail successfully sent.")
}
