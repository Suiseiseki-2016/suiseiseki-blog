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

// HandleWebhook 处理GitHub Webhook
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	// 验证签名（如果配置了secret）
	if h.secret != "" {
		signature := c.GetHeader("X-Hub-Signature-256")
		if signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少签名"})
			return
		}

		// 读取请求体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "读取请求体失败"})
			return
		}

		// 验证签名
		if !h.verifySignature(body, signature) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"})
			return
		}

		// 重新设置请求体（因为已经被读取了）
		c.Request.Body = io.NopCloser(io.Reader(bytes.NewReader(body)))
	}

	// 执行同步
	if err := h.syncService.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "同步成功"})
}

func (h *WebhookHandler) verifySignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
