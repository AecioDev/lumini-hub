package middlewares

import (
	"net/http"
	"sync"
	"time"

	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
)

// RateLimitByIP limita quantas requisições cada IP pode fazer numa janela de tempo (janela deslizante, em
// memória). Serve de freio para rotas públicas sensíveis (recuperação de senha). É por instância do serviço:
// reiniciar zera a contagem, por isso as regras que importam de verdade (ex.: quantos pedidos por usuário)
// também ficam no banco.
func RateLimitByIP(max int, window time.Duration) gin.HandlerFunc {
	var (
		mu   sync.Mutex
		hits = map[string][]time.Time{}
	)

	return func(c *gin.Context) {
		now := time.Now()
		cutoff := now.Add(-window)
		ip := c.ClientIP()

		mu.Lock()
		recent := hits[ip][:0]
		for _, t := range hits[ip] {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}
		limited := len(recent) >= max
		if !limited {
			recent = append(recent, now)
		}
		if len(recent) == 0 {
			delete(hits, ip)
		} else {
			hits[ip] = recent
		}
		// Limpeza ocasional de IPs que não aparecem mais, para o mapa não crescer sem fim.
		if len(hits) > 10000 {
			for key, times := range hits {
				if len(times) == 0 || times[len(times)-1].Before(cutoff) {
					delete(hits, key)
				}
			}
		}
		mu.Unlock()

		if limited {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Muitas tentativas. Aguarde alguns minutos e tente novamente.", "")
			c.Abort()
			return
		}
		c.Next()
	}
}
