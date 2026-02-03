package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"blog-suiseiseki/services"
)

type WebhookHandler struct {
	syncService *services.SyncService
	secret      string
}

func NewWebhookHandler(syncService *services.SyncService, secret string) *WebhookHandler {
	return &WebhookHandler{
		syncService: syncService,
		secret:      secret,
	}
}

// HandleWebhook handles GitHub Webhook requests.
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	// Verify signature if secret is configured
	if h.secret != "" {
		signature := c.GetHeader("X-Hub-Signature-256")
		if signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing signature"})
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
			return
		}

		if !h.verifySignature(body, signature) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "signature verification failed"})
			return
		}

		c.Request.Body = io.NopCloser(io.Reader(bytes.NewReader(body)))
	}

	if err := h.syncService.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sync ok"})
}

func (h *WebhookHandler) verifySignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
