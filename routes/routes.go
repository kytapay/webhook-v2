package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kytapay/webhook-v2/controllers"
)

// SetupRoutes configures all routes for the webhook service
func SetupRoutes(r *gin.Engine, webhookController *controllers.WebhookController) {
	// Health check (support both GET and HEAD for Docker healthcheck)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "KytaPay Webhook v2",
		})
	})
	r.HEAD("/health", func(c *gin.Context) {
		c.Status(200)
	})

	// Webhook routes
	payments := r.Group("/payments")
	{
		linkqu := payments.Group("/linkqu")
		{
			linkqu.POST("/qris", webhookController.HandleLinkQuQRIS)
			linkqu.POST("/ewallet", webhookController.HandleLinkQuEWallet)
		}

		pakailink := payments.Group("/pakailink")
		{
			pakailink.POST("/va", webhookController.HandlePakaiLinkVA)
		}
	}

	// Payout webhook routes
	payouts := r.Group("/payouts")
	{
		linkqu := payouts.Group("/linkqu")
		{
			linkqu.POST("/bank", webhookController.HandleLinkQuPayoutBank)
			linkqu.POST("/ewallet", webhookController.HandleLinkQuPayoutEWallet)
		}

		pakailink := payouts.Group("/pakailink")
		{
			pakailink.POST("/bank", webhookController.HandlePakaiLinkPayoutBank)
			pakailink.POST("/ewallet", webhookController.HandlePakaiLinkPayoutEWallet)
		}
	}
}

