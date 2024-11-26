package mail

import (
	"fmt"
)

func emailVerificationTemplate(link string) string {
	return fmt.Sprintf(`
    <h1> Verify your Email </h1>
    <h2> Please click on the link to verify your email </h2>
    <p> <a href="%s" target="_blank">%s</a> </p>
    <p> This link will expire in 30 minutes </p>
    `, link, link)
}
