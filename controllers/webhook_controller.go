package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kytapay/webhook-v2/helpers"
	"github.com/kytapay/webhook-v2/models"
	"github.com/kytapay/webhook-v2/repositories"
	"github.com/kytapay/webhook-v2/services"
)

type WebhookController struct {
	db                  *sql.DB
	transactionRepo     *repositories.TransactionRepository
	merchantRepo        *repositories.MerchantRepository
	walletRepo          *repositories.WalletRepository
	feesRepo            *repositories.FeesRepository
	callbackRepo        *repositories.CallbackRepository
	userRepo            *repositories.UserRepository
	telegramService     *services.TelegramService
	callbackService     *services.CallbackService
}

func NewWebhookController(db *sql.DB) *WebhookController {
	return &WebhookController{
		db:              db,
		transactionRepo: repositories.NewTransactionRepository(db),
		merchantRepo:    repositories.NewMerchantRepository(db),
		walletRepo:      repositories.NewWalletRepository(db),
		feesRepo:        repositories.NewFeesRepository(db),
		callbackRepo:    repositories.NewCallbackRepository(db),
		userRepo:        repositories.NewUserRepository(db),
		telegramService: services.NewTelegramService(),
		callbackService: services.NewCallbackService(),
	}
}

// HandleLinkQuQRIS handles webhook callback from LinkQu for QRIS
func (wc *WebhookController) HandleLinkQuQRIS(c *gin.Context) {
	// Validasi client-id dan client-secret dari header
	clientID := c.GetHeader("client-id")
	clientSecret := c.GetHeader("client-secret")

	if !helpers.VerifyLinkQuSignature(clientID, clientSecret) {
		wc.sendTelegramAlert("🚨 <b>Unauthorized Callback Attempt</b>\n\n⚠️ <b>Security Alert:</b>\n• Source: QRIS LinkQu\n• IP Address: <code>" + c.ClientIP() + "</code>\n• Client ID: <code>" + clientID + "</code>", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract data
	partnerRef, _ := data["partner_reff"].(string)
	status, _ := data["status"].(string)
	amount, _ := data["amount"].(float64)
	transactionTime, _ := data["transaction_time"].(string)
	callbackType, _ := data["type"].(string)

	if partnerRef == "" {
		wc.sendTelegramAlert("⚠️ <b>Callback Error</b>\n\n• Source: QRIS LinkQu\n• Issue: Missing payment ID", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Check callback type
	if strings.ToUpper(callbackType) == "SETTLE" {
		// Settlement notification - only send Telegram
		wc.sendSettlementNotification(partnerRef, amount, transactionTime, "QRIS", "LinkQu")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Process transaction (type = "pay")
	wc.processTransaction(partnerRef, status, amount, transactionTime, "QRIS", "LinkQu")

	// Always return 200 OK with success response
	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "2002800",
		"responseMessage": "Successful",
	})
}

// HandleLinkQuEWallet handles webhook callback from LinkQu for E-Wallet
func (wc *WebhookController) HandleLinkQuEWallet(c *gin.Context) {
	// Validasi client-id dan client-secret dari header
	clientID := c.GetHeader("client-id")
	clientSecret := c.GetHeader("client-secret")

	if !helpers.VerifyLinkQuSignature(clientID, clientSecret) {
		wc.sendTelegramAlert("🚨 <b>Unauthorized Callback Attempt</b>\n\n⚠️ <b>Security Alert:</b>\n• Source: E-Wallet LinkQu\n• IP Address: <code>" + c.ClientIP() + "</code>\n• Client ID: <code>" + clientID + "</code>", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract data
	partnerRef, _ := data["partner_reff"].(string)
	status, _ := data["status"].(string)
	amount, _ := data["amount"].(float64)
	transactionTime, _ := data["transaction_time"].(string)
	callbackType, _ := data["type"].(string)

	if partnerRef == "" {
		wc.sendTelegramAlert("⚠️ <b>Callback Error</b>\n\n• Source: E-Wallet LinkQu\n• Issue: Missing payment ID", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Check callback type
	if strings.ToUpper(callbackType) == "SETTLE" {
		// Settlement notification - only send Telegram
		wc.sendSettlementNotification(partnerRef, amount, transactionTime, "E-Wallet", "LinkQu")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Process transaction (type = "pay")
	wc.processTransaction(partnerRef, status, amount, transactionTime, "EWALLET", "LinkQu")

	// Always return 200 OK with success response
	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "2002800",
		"responseMessage": "Successful",
	})
}

// HandlePakaiLinkVA handles webhook callback from PakaiLink for Virtual Account
// Note: PakaiLink VA does not send X-SIGNATURE or X-TIMESTAMP headers
func (wc *WebhookController) HandlePakaiLinkVA(c *gin.Context) {
	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	var requestData map[string]interface{}
	if err := json.Unmarshal(body, &requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract transactionData
	transactionData, ok := requestData["transactionData"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract data
	partnerRef, _ := transactionData["partnerReferenceNo"].(string)
	paymentFlagStatus, _ := transactionData["paymentFlagStatus"].(string)
	callbackType, _ := transactionData["callbackType"].(string)
	paidAmount, _ := transactionData["paidAmount"].(map[string]interface{})
	
	var amount float64
	if paidAmount != nil {
		if valueStr, ok := paidAmount["value"].(string); ok {
			if valueFloat, err := strconv.ParseFloat(valueStr, 64); err == nil {
				amount = valueFloat
			}
		}
	}

	// Map paymentFlagStatus to status
	status := "PENDING"
	if paymentFlagStatus == "00" {
		status = "SUCCESS"
	}

	// Use current time as date (PakaiLink VA doesn't send X-TIMESTAMP)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
	}
	date := time.Now().In(loc).Format("2006-01-02T15:04:05Z07:00")

	if partnerRef == "" {
		wc.sendTelegramAlert("⚠️ <b>Callback Error</b>\n\n• Source: VA PakaiLink\n• Issue: Missing payment ID", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Check callback type
	if strings.ToLower(callbackType) == "settlement" {
		// Settlement notification - only send Telegram
		wc.sendSettlementNotification(partnerRef, amount, date, "Virtual Account", "PakaiLink")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Process transaction (callbackType = "payment")
	wc.processTransaction(partnerRef, status, amount, date, "VA", "PakaiLink")

	// Always return 200 OK with success response
	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "2002800",
		"responseMessage": "Successful",
	})
}

// processTransaction processes the transaction update
func (wc *WebhookController) processTransaction(paymentID, status string, amount float64, date, source, provider string) {
	// Get transaction
	transaction, err := wc.transactionRepo.GetTransactionByGrantID(paymentID)
	if err != nil {
		formattedAmount := helpers.FormatNumber(amount, 0)
		wc.sendTelegramAlert(fmt.Sprintf("ℹ️ <b>Transaction Not Found</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Status: %s\n• Amount: Rp %s\n• Date: %s", source, paymentID, status, formattedAmount, date), "HTML")
		return
	}

	// Get merchant payment
	merchantPayment, err := wc.merchantRepo.GetMerchantPaymentByGatewayRef(paymentID)
	if err != nil {
		formattedAmount := helpers.FormatNumber(amount, 0)
		wc.sendTelegramAlert(fmt.Sprintf("ℹ️ <b>Merchant Transaction Not Found</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Status: %s\n• Amount: Rp %s\n• Date: %s", source, paymentID, status, formattedAmount, date), "HTML")
		return
	}

	// Check if already processed
	if merchantPayment.Status != "Pending" {
		wc.sendTelegramAlert(fmt.Sprintf("⚠️ <b>Duplicate Callback Prevented</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Current Status: <b>%s</b>\n• Attempted Status: %s", source, paymentID, merchantPayment.Status, status), "HTML")
		return
	}

	// Normalize status
	normalizedStatus := helpers.NormalizeStatus(status)
	merchantNormalizedStatus := helpers.MerchantNormalizeStatus(status)

	// Update transaction
	err = wc.transactionRepo.UpdateTransaction(paymentID, normalizedStatus, int64(amount))
	if err != nil {
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Updating Transaction</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", source, paymentID, err.Error()), "HTML")
		return
	}

	// Get transactions record
	transactions, _ := wc.transactionRepo.GetTransactionsByGrantID(paymentID)

	if transactions == nil {
		// Update merchant payment only
		err = wc.merchantRepo.UpdateMerchantPayment(paymentID, merchantNormalizedStatus, amount)
		if err != nil {
			wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Updating Merchant Payment</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", source, paymentID, err.Error()), "HTML")
			return
		}
	} else {
		// Update transactions
		err = wc.transactionRepo.UpdateTransactions(paymentID, merchantNormalizedStatus, merchantNormalizedStatus)
		if err != nil {
			wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Updating Transactions</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", source, paymentID, err.Error()), "HTML")
			return
		}
	}

	// Get merchant and fees
	merchant, err := wc.merchantRepo.GetMerchantByID(*merchantPayment.MerchantID)
	if err != nil {
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Getting Merchant</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", source, paymentID, err.Error()), "HTML")
		return
	}

	feeReguler, _ := wc.feesRepo.GetFeesLimit(10, *merchantPayment.PaymentMethodID)
	// feeExpress, _ := wc.feesRepo.GetFeesExpress(10) // Not used in payment processing

	// Determine user_id
	var userID int
	if merchant != nil {
		userID = merchant.UserID
	} else if transactions != nil && transactions.UserID != nil {
		userID = *transactions.UserID
	}

	// Get wallet
	wallet, err := wc.walletRepo.GetUserWallet(userID)
	if err != nil {
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Getting Wallet</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", source, paymentID, err.Error()), "HTML")
		return
	}

	// Calculate charges
	var chargePercentage, chargeFixed float64
	if feeReguler != nil {
		chargePercentage = feeReguler.ChargePercentage
		chargeFixed = feeReguler.ChargeFixed
	} else {
		chargePercentage = 5.0
		chargeFixed = 5000.0
	}

	chargePercentageAmount := (amount * chargePercentage) / 100
	totalFee := chargePercentageAmount + chargeFixed

	// Determine payment status
	paymentStatus := merchantNormalizedStatus
	vaRealtimePaymentMethods := []int{2, 4, 6, 8} // Assuming these are realtime VA methods
	isRealtimeVA := false
	for _, pmID := range vaRealtimePaymentMethods {
		if merchantPayment.PaymentMethodID != nil && *merchantPayment.PaymentMethodID == pmID {
			isRealtimeVA = true
			break
		}
	}

	if merchantNormalizedStatus == "Success" && !isRealtimeVA {
		// Check if payment method requires settlement
		settlementMethods := []int{1, 2, 4, 6, 8, 11, 12, 13, 14, 15, 19}
		for _, pmID := range settlementMethods {
			if merchantPayment.PaymentMethodID != nil && *merchantPayment.PaymentMethodID == pmID {
				paymentStatus = "Pending_Settlement"
				break
			}
		}
	}

	// Update wallet balance for VA Success (non-realtime)
	if source == "VA" && merchantNormalizedStatus == "Success" && !isRealtimeVA {
		if transactions == nil {
			newBalance := wallet.Balance + (amount - totalFee)
			err = wc.walletRepo.UpdateWalletBalance(userID, newBalance)
			if err != nil {
				wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Updating Wallet</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", source, paymentID, err.Error()), "HTML")
				return
			}
		} else {
			newBalance := wallet.Balance + transactions.Total
			err = wc.walletRepo.UpdateWalletBalance(userID, newBalance)
			if err != nil {
				wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Updating Wallet</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", source, paymentID, err.Error()), "HTML")
				return
			}
		}
	}

	// Create transaction record if not exists
	if transactions == nil {
		currencyID := 1
		transactionTypeID := 10
		orderID := transaction.OrderID
		transactionsData := models.TransactionsData{
			UserID:                &userID,
			CurrencyID:            &currencyID,
			PaymentMethodID:       merchantPayment.PaymentMethodID,
			MerchantID:            merchantPayment.MerchantID,
			UUID:                  &orderID,
			GrantID:               &paymentID,
			TransactionReferenceID: 1,
			TransactionTypeID:     &transactionTypeID,
			UserType:              "registered",
			Subtotal:              amount,
			Percentage:            chargePercentage,
			ChargePercentage:      chargePercentageAmount,
			ChargeFixed:           chargeFixed,
			Total:                 amount - totalFee,
			PaymentStatus:         &merchantNormalizedStatus,
			Status:                paymentStatus,
		}

		err = wc.transactionRepo.CreateTransaction(transactionsData)
		if err != nil {
			wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Creating Transaction</b>\n\n• Source: %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", source, paymentID, err.Error()), "HTML")
			return
		}
	}

	// Send callback to merchant (only for payment callbacks, not settlement)
	if transactions == nil {
		// Use merchantNormalizedStatus (Success/Pending/Failed) instead of normalizedStatus (success/pending/expires)
		payloads := services.BuildPayloadV2(transaction, paymentID, merchant.BusinessName, merchantNormalizedStatus, date)
		payload := payloads[transaction.PaymentMethod]
		wc.sendCallbackToMerchant(transaction, payload)
	}

	// Send Telegram notification
	formattedAmount := helpers.FormatNumber(amount, 0)
	message := fmt.Sprintf("✅ <b>Pembayaran Berhasil</b>\n\n📋 <b>Detail Transaksi:</b>\n• ID Transaksi: <code>%s</code>\n• Order ID: <code>%s</code>\n• Metode: %s\n• Provider: %s\n• Jumlah: <b>Rp %s</b>\n• Status: <b>%s</b>\n• Waktu: %s",
		paymentID, transaction.OrderID, getPaymentMethodName(source), provider, formattedAmount, normalizedStatus, date)
	wc.sendTelegramAlert(message, "HTML")
}

// sendSettlementNotification sends settlement notification to Telegram
func (wc *WebhookController) sendSettlementNotification(paymentID string, amount float64, date string, paymentMethod, provider string) {
	// Get transaction for additional info
	transaction, err := wc.transactionRepo.GetTransactionByGrantID(paymentID)
	if err != nil {
		// If transaction not found, send basic notification
		formattedAmount := helpers.FormatNumber(amount, 0)
		message := fmt.Sprintf("💰 <b>Settlement Notification</b>\n\n📋 <b>Detail:</b>\n• ID Transaksi: <code>%s</code>\n• Metode: %s\n• Provider: %s\n• Jumlah: <b>Rp %s</b>\n• Waktu: %s",
			paymentID, paymentMethod, provider, formattedAmount, date)
		wc.sendTelegramAlert(message, "HTML")
		return
	}

	formattedAmount := helpers.FormatNumber(amount, 0)
	message := fmt.Sprintf("💰 <b>Settlement Berhasil</b>\n\n📋 <b>Detail Transaksi:</b>\n• ID Transaksi: <code>%s</code>\n• Order ID: <code>%s</code>\n• Metode: %s\n• Provider: %s\n• Jumlah: <b>Rp %s</b>\n• Waktu Settlement: %s",
		paymentID, transaction.OrderID, paymentMethod, provider, formattedAmount, date)
	wc.sendTelegramAlert(message, "HTML")
}

// getPaymentMethodName returns formatted payment method name
func getPaymentMethodName(source string) string {
	switch source {
	case "QRIS":
		return "QRIS"
	case "EWALLET":
		return "E-Wallet"
	case "VA":
		return "Virtual Account"
	default:
		return source
	}
}


// sendCallbackToMerchant sends callback to merchant
func (wc *WebhookController) sendCallbackToMerchant(transaction *models.TransactionInfo, payload interface{}) {
	_, err := wc.callbackRepo.GetCallbackByTransactionInfoID(transaction.ID)
	if err != nil {
		return
	}

	statusCode, responseBody, err := wc.callbackService.SendCallback(transaction.NotifyURL, payload, transaction.Token)
	if err != nil {
		errorMessage := "408 - Request Timeout"
		if err.Error() != "" {
			errorMessage = err.Error()
		}
		wc.callbackRepo.UpdateCallback(transaction.ID, "Failed", errorMessage, responseBody, payload)
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Request Timeout (Callback)</b>\n\n• Error: <code>%s</code>\n• URL: <code>%s</code>", err.Error(), transaction.NotifyURL), "HTML")
		return
	}

	if statusCode >= 200 && statusCode < 300 {
		wc.callbackRepo.UpdateCallback(transaction.ID, "Success", fmt.Sprintf("%d - Success", statusCode), responseBody, payload)
	} else {
		errorMessage := fmt.Sprintf("HTTP %d", statusCode)
		switch statusCode {
		case 400:
			errorMessage = "400 - Bad Request"
		case 401:
			errorMessage = "401 - Unauthorized"
		case 403:
			errorMessage = "403 - Forbidden"
		case 404:
			errorMessage = "404 - Not Found"
		case 422:
			errorMessage = "422 - Unprocessable Entity"
		case 500:
			errorMessage = "500 - Internal Server Error"
		case 502:
			errorMessage = "502 - Bad Gateway"
		case 503:
			errorMessage = "503 - Service Unavailable"
		}
		wc.callbackRepo.UpdateCallback(transaction.ID, "Failed", errorMessage, responseBody, payload)
	}
}

// sendTelegramAlert sends alert to Telegram
func (wc *WebhookController) sendTelegramAlert(message, parseMode string) {
	_ = wc.telegramService.SendMessage(message, parseMode)
}

// HandleLinkQuPayoutBank handles webhook callback from LinkQu for Bank Payout
func (wc *WebhookController) HandleLinkQuPayoutBank(c *gin.Context) {
	// Validasi client-id dan client-secret dari header
	clientID := c.GetHeader("client-id")
	clientSecret := c.GetHeader("client-secret")

	if !helpers.VerifyLinkQuSignature(clientID, clientSecret) {
		wc.sendTelegramAlert("🚨 <b>Unauthorized Callback Attempt</b>\n\n⚠️ <b>Security Alert:</b>\n• Source: Bank Payout LinkQu\n• IP Address: <code>"+c.ClientIP()+"</code>\n• Client ID: <code>"+clientID+"</code>", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract data
	partnerRef, _ := data["partner_reff"].(string)
	status, _ := data["status"].(string)
	amount, _ := data["amount"].(float64)
	transactionTime, _ := data["transaction_time"].(string)

	if partnerRef == "" {
		wc.sendTelegramAlert("⚠️ <b>Callback Error</b>\n\n• Source: Bank Payout LinkQu\n• Issue: Missing payment ID", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Process payout transaction
	wc.processPayoutTransaction(partnerRef, status, amount, transactionTime, "Bank", "LinkQu")

	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "2002800",
		"responseMessage": "Successful",
	})
}

// HandleLinkQuPayoutEWallet handles webhook callback from LinkQu for E-Wallet Payout
func (wc *WebhookController) HandleLinkQuPayoutEWallet(c *gin.Context) {
	// Validasi client-id dan client-secret dari header
	clientID := c.GetHeader("client-id")
	clientSecret := c.GetHeader("client-secret")

	if !helpers.VerifyLinkQuSignature(clientID, clientSecret) {
		wc.sendTelegramAlert("🚨 <b>Unauthorized Callback Attempt</b>\n\n⚠️ <b>Security Alert:</b>\n• Source: E-Wallet Payout LinkQu\n• IP Address: <code>"+c.ClientIP()+"</code>\n• Client ID: <code>"+clientID+"</code>", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract data
	partnerRef, _ := data["partner_reff"].(string)
	status, _ := data["status"].(string)
	amount, _ := data["amount"].(float64)
	transactionTime, _ := data["transaction_time"].(string)

	if partnerRef == "" {
		wc.sendTelegramAlert("⚠️ <b>Callback Error</b>\n\n• Source: E-Wallet Payout LinkQu\n• Issue: Missing payment ID", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Process payout transaction
	wc.processPayoutTransaction(partnerRef, status, amount, transactionTime, "E-Wallet", "LinkQu")

	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "2002800",
		"responseMessage": "Successful",
	})
}

// HandlePakaiLinkPayoutBank handles webhook callback from PakaiLink for Bank Payout
// Note: PakaiLink Payout does not send X-SIGNATURE or X-TIMESTAMP headers
func (wc *WebhookController) HandlePakaiLinkPayoutBank(c *gin.Context) {
	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	var requestData map[string]interface{}
	if err := json.Unmarshal(body, &requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract transactionData
	transactionData, ok := requestData["transactionData"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract data
	partnerRef, _ := transactionData["partnerReferenceNo"].(string)
	paymentFlagStatus, _ := transactionData["paymentFlagStatus"].(string)
	paidAmount, _ := transactionData["paidAmount"].(map[string]interface{})

	var amount float64
	if paidAmount != nil {
		if valueStr, ok := paidAmount["value"].(string); ok {
			if valueFloat, err := strconv.ParseFloat(valueStr, 64); err == nil {
				amount = valueFloat
			}
		}
	}

	// Map paymentFlagStatus to status
	status := "PENDING"
	if paymentFlagStatus == "00" {
		status = "SUCCESS"
	} else if paymentFlagStatus != "" {
		status = "FAILED"
	}

	// Use current time as date (PakaiLink Payout doesn't send X-TIMESTAMP)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
	}
	date := time.Now().In(loc).Format("2006-01-02T15:04:05Z07:00")

	if partnerRef == "" {
		wc.sendTelegramAlert("⚠️ <b>Callback Error</b>\n\n• Source: Bank Payout PakaiLink\n• Issue: Missing payment ID", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Process payout transaction
	wc.processPayoutTransaction(partnerRef, status, amount, date, "Bank", "PakaiLink")

	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "2002800",
		"responseMessage": "Successful",
	})
}

// HandlePakaiLinkPayoutEWallet handles webhook callback from PakaiLink for E-Wallet Payout
// Note: PakaiLink Payout does not send X-SIGNATURE or X-TIMESTAMP headers
func (wc *WebhookController) HandlePakaiLinkPayoutEWallet(c *gin.Context) {
	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	var requestData map[string]interface{}
	if err := json.Unmarshal(body, &requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract transactionData
	transactionData, ok := requestData["transactionData"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Extract data
	partnerRef, _ := transactionData["partnerReferenceNo"].(string)
	paymentFlagStatus, _ := transactionData["paymentFlagStatus"].(string)
	paidAmount, _ := transactionData["paidAmount"].(map[string]interface{})

	var amount float64
	if paidAmount != nil {
		if valueStr, ok := paidAmount["value"].(string); ok {
			if valueFloat, err := strconv.ParseFloat(valueStr, 64); err == nil {
				amount = valueFloat
			}
		}
	}

	// Map paymentFlagStatus to status
	status := "PENDING"
	if paymentFlagStatus == "00" {
		status = "SUCCESS"
	} else if paymentFlagStatus != "" {
		status = "FAILED"
	}

	// Use current time as date (PakaiLink Payout doesn't send X-TIMESTAMP)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
	}
	date := time.Now().In(loc).Format("2006-01-02T15:04:05Z07:00")

	if partnerRef == "" {
		wc.sendTelegramAlert("⚠️ <b>Callback Error</b>\n\n• Source: E-Wallet Payout PakaiLink\n• Issue: Missing payment ID", "HTML")
		c.JSON(http.StatusOK, gin.H{
			"responseCode":    "2002800",
			"responseMessage": "Successful",
		})
		return
	}

	// Process payout transaction
	wc.processPayoutTransaction(partnerRef, status, amount, date, "E-Wallet", "PakaiLink")

	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "2002800",
		"responseMessage": "Successful",
	})
}

// processPayoutTransaction processes the payout transaction update
func (wc *WebhookController) processPayoutTransaction(paymentID, status string, amount float64, date, paymentMethod, provider string) {
	// Get transaction
	transaction, err := wc.transactionRepo.GetTransactionByGrantID(paymentID)
	if err != nil {
		formattedAmount := helpers.FormatNumber(amount, 0)
		wc.sendTelegramAlert(fmt.Sprintf("ℹ️ <b>Transaction Not Found</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Status: %s\n• Amount: Rp %s\n• Date: %s", paymentMethod, provider, paymentID, status, formattedAmount, date), "HTML")
		return
	}

	// Get transactions record
	transactions, err := wc.transactionRepo.GetTransactionsByGrantID(paymentID)
	if err != nil || transactions == nil {
		formattedAmount := helpers.FormatNumber(amount, 0)
		wc.sendTelegramAlert(fmt.Sprintf("ℹ️ <b>Transactions Record Not Found</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Status: %s\n• Amount: Rp %s\n• Date: %s", paymentMethod, provider, paymentID, status, formattedAmount, date), "HTML")
		return
	}

	// Get merchant payout
	merchantPayout, err := wc.merchantRepo.GetMerchantPayoutByGatewayRef(paymentID)
	if err != nil {
		formattedAmount := helpers.FormatNumber(amount, 0)
		wc.sendTelegramAlert(fmt.Sprintf("ℹ️ <b>Merchant Payout Not Found</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Status: %s\n• Amount: Rp %s\n• Date: %s", paymentMethod, provider, paymentID, status, formattedAmount, date), "HTML")
		return
	}

	// Check if already processed
	if merchantPayout.Status != "Pending" {
		wc.sendTelegramAlert(fmt.Sprintf("⚠️ <b>Duplicate Callback Prevented</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Current Status: <b>%s</b>\n• Attempted Status: %s", paymentMethod, provider, paymentID, merchantPayout.Status, status), "HTML")
		return
	}

	// Normalize status
	normalizedStatus := helpers.NormalizeStatus(status)
	normalizedStatus2 := helpers.MerchantNormalizeStatus(status)

	// Update transaction
	err = wc.transactionRepo.UpdateTransaction(paymentID, normalizedStatus, int64(amount))
	if err != nil {
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Updating Transaction</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", paymentMethod, provider, paymentID, err.Error()), "HTML")
		return
	}

	// Update transactions
	err = wc.transactionRepo.UpdateTransactions(paymentID, normalizedStatus2, normalizedStatus2)
	if err != nil {
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Updating Transactions</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", paymentMethod, provider, paymentID, err.Error()), "HTML")
		return
	}

	// Update merchant payout
	err = wc.merchantRepo.UpdateMerchantPayout(paymentID, normalizedStatus2, amount)
	if err != nil {
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Updating Merchant Payout</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", paymentMethod, provider, paymentID, err.Error()), "HTML")
		return
	}

	// Get merchant and user
	merchant, err := wc.merchantRepo.GetMerchantByID(*merchantPayout.MerchantID)
	if err != nil {
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Getting Merchant</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", paymentMethod, provider, paymentID, err.Error()), "HTML")
		return
	}

	// Get user (need to add GetUserByID to repository)
	userID := merchant.UserID
	wallet, err := wc.walletRepo.GetUserWallet(userID)
	if err != nil {
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Getting Wallet</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", paymentMethod, provider, paymentID, err.Error()), "HTML")
		return
	}

	// Get fees
	feeReguler, _ := wc.feesRepo.GetFeesLimit(9, *merchantPayout.PaymentMethodID)
	feeExpress, _ := wc.feesRepo.GetFeesExpress(9)

	// Calculate charges
	var chargePercentage, chargeFixed float64
	if feeReguler != nil && feeExpress != nil {
		chargePercentage = feeReguler.ChargePercentage
		chargeFixed = feeReguler.ChargeFixed
	} else {
		chargePercentage = 1.5
		chargeFixed = 5000.0
	}

	chargePercentageAmount := (amount * chargePercentage) / 100

	// Get user to check role_id
	user, err := wc.userRepo.GetUserByID(userID)
	if err != nil {
		wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Getting User</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", paymentMethod, provider, paymentID, err.Error()), "HTML")
		return
	}

	// Calculate total fee based on role_id
	var finalTotalFee float64
	if user.RoleID != nil && *user.RoleID == 3 {
		// Role ID 3 = reguler fee
		finalTotalFee = chargePercentageAmount + chargeFixed
	} else {
		// Other roles = express fee
		var chargePercentageExpress, chargeFixedExpress float64
		if feeExpress != nil {
			chargePercentageExpress = feeExpress.ChargePercentage
			chargeFixedExpress = feeExpress.ChargeFixed
		} else {
			chargePercentageExpress = 1.5
			chargeFixedExpress = 7000.0
		}
		finalTotalFee = (amount*chargePercentageExpress)/100 + chargeFixedExpress
	}

	// Update wallet balance if status is Success (deduct amount + fee)
	if normalizedStatus2 == "Success" {
		newBalance := wallet.Balance - (amount + finalTotalFee)
		err = wc.walletRepo.UpdateWalletBalance(userID, newBalance)
		if err != nil {
			wc.sendTelegramAlert(fmt.Sprintf("❌ <b>Error Updating Wallet</b>\n\n• Source: %s Payout %s\n• Payment ID: <code>%s</code>\n• Error: <code>%s</code>", paymentMethod, provider, paymentID, err.Error()), "HTML")
			return
		}
	}

	// Send callback to merchant (V2 format only)
	payloads := services.BuildPayloadV2Payout(transaction, paymentID, normalizedStatus2, date)
	payload := payloads["PAYOUTS"]
	wc.sendCallbackToMerchant(transaction, payload)

	// Send Telegram notification
	formattedAmount := helpers.FormatNumber(amount, 0)
	message := fmt.Sprintf("✅ <b>Payout Status Updated</b>\n\n📋 <b>Detail Transaksi:</b>\n• ID Transaksi: <code>%s</code>\n• Order ID: <code>%s</code>\n• Metode: %s\n• Provider: %s\n• Jumlah: <b>Rp %s</b>\n• Status: <b>%s</b>\n• Waktu: %s",
		paymentID, transaction.OrderID, paymentMethod, provider, formattedAmount, normalizedStatus2, date)
	wc.sendTelegramAlert(message, "HTML")
}

