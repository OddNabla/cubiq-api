package service

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	// "encoding/json"
	"io"
	// "os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// FirebaseUploader handles uploading files to Firebase Storage
type FirebaseUploader struct {
}

func (fu *FirebaseUploader) UploadFile(ctx context.Context, data []byte, objectName, contentType string) (string, error) {
	bucketName := "cubiqapi.firebasestorage.app"
	firebaseServiceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(firebaseServiceAccountPath))
	if err != nil {
		return "", err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)
	writer.ContentType = contentType
	writer.Metadata = map[string]string{
		"uploadedBy": "GolangUploader",
	}

	reader := bytes.NewReader(data)
	if _, err := io.Copy(writer, reader); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	url, err := getUrl(bucketName, objectName)
	return url, err
}

func getUrl(bucketName, objectName string) (string, error) {
	firebaseServiceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(firebaseServiceAccountPath))
	if err != nil {
		return "", err
	}
	defer client.Close()

	privateKeyBytes := []byte(os.Getenv("FIREBASE_PRIVATE_KEY"))

	url, err := storage.SignedURL(bucketName, objectName, &storage.SignedURLOptions{
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
		GoogleAccessID: os.Getenv("FIREBASE_CLIENT_EMAIL"),
		SignBytes: func(b []byte) ([]byte, error) {
			privateKey := privateKeyBytes
			block, _ := pem.Decode(privateKey)
			if block == nil || block.Type != "PRIVATE KEY" {
				return nil, fmt.Errorf("failed to decode PEM block containing private key")
			}

			parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key: %v", err)
			}

			rsaKey, ok := parsedKey.(*rsa.PrivateKey)
			if !ok {
				return nil, fmt.Errorf("not an RSA private key")
			}

			hashed := sha256.Sum256(b)
			return rsa.SignPKCS1v15(nil, rsaKey, crypto.SHA256, hashed[:])
		},
	})
	if err != nil {
		return "", err
	}
	return url, nil
}
