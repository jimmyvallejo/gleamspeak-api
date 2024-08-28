package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/common"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
)

type signedURLRequest struct {
	Filename string `json:"filename"`
	Filetype string `json:"filetype"`
}

func (h *Handlers) GetSignedURL(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var request signedURLRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	timestamp := time.Now().UnixNano()
	uniqueID := generateUniqueID()
	ext := path.Ext(request.Filename)
	filename := fmt.Sprintf("%d-%s-%s%s", timestamp, user.ID, uniqueID, ext)

	key := path.Join("public", filename)

	presignClient := s3.NewPresignClient(h.S3)

	presignResult, err := presignClient.PresignPutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket:      aws.String("gleamspeak-bucket"),
			Key:         aws.String(key),
			ContentType: aws.String(request.Filetype),
		},
		s3.WithPresignExpires(time.Minute*15),
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate pre-signed URL")
		return
	}

	publicURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "gleamspeak-bucket", key)

	response := SignedURLResponse{
		URL:       presignResult.URL,
		PublicURL: publicURL,
	}

	respondWithJSON(w, http.StatusOK, response)
}
