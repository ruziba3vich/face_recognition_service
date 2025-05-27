package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/face_recognition_service/genprotos/face_recognition_service"
	lgg "github.com/ruziba3vich/prodonik_lgger"
)

type FaceRecognitionHandler struct {
	faceEmbedderClient face_recognition_service.FaceEmbedderClient
	logger             *lgg.Logger
}

func NewFaceRecognitionHandler(faceEmbedderClient face_recognition_service.FaceEmbedderClient, logger *lgg.Logger) *FaceRecognitionHandler {
	return &FaceRecognitionHandler{
		faceEmbedderClient: faceEmbedderClient,
		logger:             logger,
	}
}

func (h *FaceRecognitionHandler) HandleImageEmbedding(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.faceEmbedderClient.GetEmbedding(ctx, &face_recognition_service.ImageRequest{Image: imageData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("gRPC error: %v", err)})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"embedding_length": len(resp.Embedding),
		"embedding":        resp.Embedding,
	})
}

/*

	conn, err := grpc.Dial("localhost:7178", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC connection failed"})
		return
	}
	defer conn.Close()

	client := face_recognition_service.NewFaceEmbedderClient(conn)
*/
