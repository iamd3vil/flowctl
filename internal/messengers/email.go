package messengers

import (
	"fmt"
	"log/slog"

	"github.com/cvhariharan/flowctl/internal/config"
	"github.com/knadh/smtppool/v2"
)

// EmailMessenger sends emails using an SMTP connection pool
type EmailMessenger struct {
	pool   *smtppool.Pool
	from   string
	logger *slog.Logger
}

// NewEmailMessenger creates a new EmailMessenger with the given SMTP configuration
func NewEmailMessenger(cfg config.SMTPConfig, logger *slog.Logger) (*EmailMessenger, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("email messenger is disabled")
	}

	var sslType smtppool.SSLType
	switch cfg.SSL {
	case "tls":
		sslType = smtppool.SSLTLS
	case "starttls":
		sslType = smtppool.SSLSTARTTLS
	default:
		sslType = smtppool.SSLNone
	}

	pool, err := smtppool.New(smtppool.Opt{
		Host:     cfg.Host,
		Port:     cfg.Port,
		MaxConns: cfg.MaxConns,
		Auth:     &smtppool.LoginAuth{Username: cfg.Username, Password: cfg.Password},
		SSL:      sslType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP pool: %w", err)
	}

	fromAddr := cfg.FromAddress
	if cfg.FromName != "" {
		fromAddr = fmt.Sprintf("%s <%s>", cfg.FromName, cfg.FromAddress)
	}

	return &EmailMessenger{
		pool:   pool,
		from:   fromAddr,
		logger: logger,
	}, nil
}

// Send sends an email message to all recipients
func (e *EmailMessenger) Send(msg Message) error {
	if len(msg.Recipients) == 0 {
		return nil
	}

	// Extract email addresses from recipients
	to := make([]string, 0, len(msg.Recipients))
	for _, r := range msg.Recipients {
		if r.Email != "" {
			to = append(to, r.Email)
		}
	}

	if len(to) == 0 {
		return nil
	}

	email := smtppool.Email{
		From:    e.from,
		To:      to,
		Subject: msg.Title,
		HTML:    []byte(msg.Body),
	}

	if err := e.pool.Send(email); err != nil {
		e.logger.Error("failed to send email",
			"to", to,
			"subject", msg.Title,
			"error", err,
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	e.logger.Debug("email sent",
		"to", to,
		"subject", msg.Title,
	)
	return nil
}

// Close closes the SMTP connection pool
func (e *EmailMessenger) Close() {
	if e.pool != nil {
		e.pool.Close()
	}
}
