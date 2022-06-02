package main

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type qrInput struct {
	Content string `json:"content" binding:"required"`
}

func qrHandler(c *gin.Context) {
	var input qrInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	png, err := qrcode.Encode(input.Content, qrcode.Medium, 512)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pngReader := bytes.NewReader(png)
	if err := printImage8bitDouble(pngReader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	nocut := c.Query("nocut")
	if nocut == "1" || nocut == "true" {
		return
	}
	cutPaper()
}
