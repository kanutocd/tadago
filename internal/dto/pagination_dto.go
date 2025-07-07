package dto

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type PaginationQuery struct {
	Cursor string `form:"cursor" json:"cursor,omitempty"`
	Limit  int    `form:"limit" json:"limit,omitempty" binding:"min=1,max=100"`
}

type PaginationResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Limit      int    `json:"limit"`
	Count      int    `json:"count"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type Cursor struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func EncodeCursor(cursor Cursor) string {
	if cursor.ID == "" {
		return ""
	}

	jsonData, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(jsonData)
}

func DecodeCursor(encodedCursor string) (Cursor, error) {
	var cursor Cursor
	if encodedCursor == "" {
		return cursor, nil
	}

	decodedData, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return cursor, err
	}

	err = json.Unmarshal(decodedData, &cursor)
	return cursor, err
}
