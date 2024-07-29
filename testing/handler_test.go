package testing

import (
	"articles-system/lib/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreate_ValidPayload(t *testing.T) {
    customLogger, db, redis, app := SetupTestApp()
    defer func() {
        customLogger.Close()
        db.Close()
        redis.Close()
    }()

    payload := models.AddArticle{
        Author: "John Doe",
        Title:  "Sample Article 1",
        Body:   "This is the first sample article.",
    }

    testCreateArticle(t, app, payload, fiber.StatusOK)
}

func TestCreate_InvalidPayload(t *testing.T) {
    customLogger, db, redis, app := SetupTestApp()
    defer func() {
        customLogger.Close()
        db.Close()
        redis.Close()
    }()

    payload := models.AddArticle{}

    testCreateArticle(t, app, payload, fiber.StatusBadRequest)
}

func testCreateArticle(t *testing.T, app *fiber.App, payload models.AddArticle, expectedStatus int) {
    payloadBytes, err := json.Marshal(payload)
    assert.Nil(t, err, "Error marshalling payload")

    req := httptest.NewRequest(fiber.MethodPost, "/articles", bytes.NewReader(payloadBytes))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req, -1)
    if err != nil {
        assert.Nil(t, err, "Test create article failed")
        return
    }

    assert.Equal(t, expectedStatus, resp.StatusCode)
}

func TestGetArticles(t *testing.T) {
	customLogger, db, redis, app := SetupTestApp()
	defer func() {
		customLogger.Close()
		db.Close()
		redis.Close()
	}()

	queryParams := [2][2]string{
		{"query", "sample"},
		{"author", "John%20Doe"},
	}

	var queryParamsString string
	if len(queryParams) > 0 {
		for i, value := range queryParams {
			if i == 0 {
				queryParamsString += "?" + value[0] + "=" + value[1]
			} else {
				queryParamsString += "&" + value[0] + "=" + value[1]
			}
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/articles"+queryParamsString, nil)

	resp, err := app.Test(req, -1)
	if err != nil {
		assert.Nil(t, err, "Test get articles failed")
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		assert.Nil(t, err, "Error reading response body")
		return
	}

	var responseData models.Response
	if err = app.Config().JSONDecoder(respBody, &responseData); err != nil {
		assert.Nil(t, err, "Error decoding response")
		return
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.NotNil(t, responseData.Data)
}
