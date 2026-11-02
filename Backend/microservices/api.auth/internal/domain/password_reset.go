package domain

import "time"

// PasswordResetToken é um pedido de redefinição de senha. Só o hash (SHA-256) do token é gravado: o valor
// puro existe apenas no link enviado por e-mail. Vale uma vez (UsedAt) e por pouco tempo (ExpiresAt). Os
// tokens usados ou substituídos não são apagados: a contagem de pedidos recentes (limite por usuário) olha
// para eles.
type PasswordResetToken struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;index"`
	TokenHash   string    `gorm:"size:64;not null;uniqueIndex"`
	ExpiresAt   time.Time `gorm:"not null"`
	UsedAt      *time.Time
	RequestedIP string `gorm:"size:45"`
	CreatedAt   time.Time
}

// TableName especifica o nome da tabela
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// ForgotPasswordRequest é o pedido de envio do link de redefinição, pelo e-mail da conta.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email,max=100"`
}

// ResetPasswordRequest troca a senha usando o token recebido por e-mail. A senha segue a mesma regra do cadastro de usuário.
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required,max=200"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=72"`
}
