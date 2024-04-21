package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func UploadFile(c *fiber.Ctx) error {
	ctx := context.Background()
	bucketName := os.Getenv("MINIO_BUCKET")
	file, err := c.FormFile("fileUpload")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get Buffer from file
	buffer, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	defer buffer.Close()

	// Create minio connection.
	minioClient, err := MinioConnection()
	if err != nil {
		// Return status 500 and minio connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	objectName := file.Filename
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	// Upload the zip file with PutObject
	info, err := minioClient.PutObject(ctx, bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"info":  info,
	})
}
