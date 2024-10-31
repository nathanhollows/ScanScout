package services

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/a-h/templ"
	emails "github.com/nathanhollows/Rapua/internal/templates/emails"
	"github.com/nathanhollows/Rapua/models"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	// SendPasswordReset(ctx context.Context, user models.User) (*rest.Response, error)
	SendVerificationEmail(ctx context.Context, user models.User) (*rest.Response, error)
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}

func (s emailService) SendVerificationEmail(ctx context.Context, user models.User) (*rest.Response, error) {
	from := mail.NewEmail(os.Getenv("SENDGRID_FROM_NAME"), os.Getenv("SENDGRID_FROM_EMAIL"))
	to := mail.NewEmail(user.Name, user.Email)
	subject := "Please verify your email"

	url := templ.URL(os.Getenv("SITE_URL") + "/verify-email/" + user.EmailToken)

	plainTextContent := `Tap the link below to finish verifying your account with Rapua. If you didn't register, you can safely ignore this email and your email will be automatically deleted from our system.
Verify your email

	%v

Cheers,
Nathan`
	plainTextContent = fmt.Sprintf(plainTextContent, url)

	// Render the html email template
	w := new(bytes.Buffer)
	c := emails.VerifyEmail(url)
	c.Render(ctx, w)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, w.String())
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	return response, err
}
