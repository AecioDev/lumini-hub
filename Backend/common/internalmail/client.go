// Package internalmail é o cliente das rotas internas de e-mail do api.core: quem precisa mandar e-mail
// (hoje, o api.auth na recuperação de senha) pede ao api.core, que tem a configuração SMTP e os templates.
// A chamada vai direto ao api.core (CORE_SERVICE_URL), sem passar pelo gateway, com um token de serviço.
package internalmail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"lumini-hub/common/config"
	"lumini-hub/common/utils"
)

// Sender pede o envio de e-mails ao api.core. Existe como interface para os testes trocarem a chamada de rede.
type Sender interface {
	SendPasswordReset(ctx context.Context, req PasswordResetRequest) error
}

// PasswordResetRequest é o pedido de envio do e-mail de redefinição de senha. Token é o valor puro (o
// link é montado pelo api.core, que conhece o endereço do sistema).
type PasswordResetRequest struct {
	To             string `json:"to"`
	Name           string `json:"name"`
	Token          string `json:"token"`
	CompanyID      *uint  `json:"company_id"`
	ExpiresMinutes int    `json:"expires_minutes"`
}

// Client é o Sender de verdade, por HTTP.
type Client struct {
	baseURL string
	caller  string
	cfg     *config.Config
	http    *http.Client
}

// NewClient cria o cliente. caller é o nome do serviço que chama (ex.: "api.auth").
func NewClient(caller string, cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.Services.CoreURL,
		caller:  caller,
		cfg:     cfg,
		// Maior que o limite de conexão SMTP do api.core (20 s): o erro de lá chega aqui em vez de estourar antes.
		http: &http.Client{Timeout: 30 * time.Second},
	}
}

// SendPasswordReset pede ao api.core o envio do e-mail de redefinição de senha.
func (c *Client) SendPasswordReset(ctx context.Context, req PasswordResetRequest) error {
	return c.post(ctx, "/api/internal/mail/password-reset", req)
}

func (c *Client) post(ctx context.Context, path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token, err := utils.GenerateServiceToken(c.caller, c.cfg)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return fmt.Errorf("api.core indisponível: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	var apiResp utils.Response
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	_ = json.Unmarshal(raw, &apiResp)
	detail := apiResp.Message
	if apiResp.Error != "" {
		detail += ": " + apiResp.Error
	}
	return fmt.Errorf("api.core recusou o envio (HTTP %d): %s", resp.StatusCode, detail)
}
