package model

type (
	EmailConfig struct {
		ID            string
		OutboundEmail string `json:"outboundEmail" validate:"required,email"`
		SMTPHost      string `json:"smtpHost" validate:"required"`
		SMTPPort      int    `json:"smtpPort" validate:"required"`
		SMTPUser      string `json:"smtpUser" validate:"required"`
		SMTPPass      string `json:"smtpPass" validate:"required"`
		HasError      bool
		IsActive      bool `json:"isActive"`
	}

	UpdateEmailConfig struct {
		ID            string
		OutboundEmail string `json:"outboundEmail"`
		SMTPHost      string `json:"smtpHost"`
		SMTPPort      int    `json:"smtpPort"`
		SMTPUser      string `json:"smtpUser"`
		SMTPPass      string `json:"smtpPass"`
		HasError      bool
		IsActive      bool `json:"isActive"`
	}
)
