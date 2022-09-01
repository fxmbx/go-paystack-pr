package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Server struct {
}

type InitiatePaymentRequest struct {
	Email  string `json:"email" binding:"required"`
	Amount int64  `json:"amount" binding:"required,min=100"`
}

type PaystackResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

func (server *Server) InitiatePayment(ctx *gin.Context) {
	var req InitiatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message 31": err.Error(), "success": false, "data": nil})
		return
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message 36": err.Error(), "success": false, "data": nil})
		return
	}

	request, err := http.NewRequest("POST", "https://api.paystack.co/transaction/initialize", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PAYSTACK_SECRET_KEY")))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message 51": err.Error(), "success": false, "data": nil})
		return
	}

	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized {
		ctx.JSON(http.StatusBadRequest, gin.H{"message 56": fmt.Sprintf("unauthorized, %s", err.Error()), "success": false, "data": nil})
		return
	} else if response.StatusCode != http.StatusOK {
		log.Println(response)
		ctx.JSON(http.StatusBadRequest, gin.H{"message 61": response.StatusCode, "success": false, "data": nil})
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message 67": err.Error(), "success": false, "data": nil})
		return
	}

	var jsonFromPayStack PaystackResponse
	err = json.Unmarshal(body, &jsonFromPayStack)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message 73": err.Error(), "success": false, "data": nil})
		return
	}
	// err = json.NewDecoder(request.Body).Decode(&jsonFromPayStack)

	ctx.JSON(200, jsonFromPayStack)
}

func (server *Server) VerifyPayment(ctx *gin.Context) {
	reference := ctx.Query("reference")
	fmt.Println("💨" + reference)
	if len(reference) <= 0 || reference == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Reference not found, add reference to url param i.e <reference=the_reference_code>", "success": false, "data": nil})
		return
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.paystack.co/transaction/verify/%s", reference), nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PAYSTACK_SECRET_KEY")))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	var jsonResponse PaystackResponse
	err = json.Unmarshal(body, &jsonResponse)
	if jsonResponse.Data["status"] != "success" {
		fmt.Println("yippy")
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message 119": err.Error(), "success": false, "data": nil})
		return
	}
	ctx.JSON(200, gin.H{"success": jsonResponse.Status, "message": jsonResponse.Message, "data": jsonResponse.Data})

}
