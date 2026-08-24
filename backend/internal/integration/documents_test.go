package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

// sampleDocumentPath points at the shared fixture already used by the
// PDF/text reader unit tests, so upload tests exercise a real file on disk
// rather than an in-memory stand-in.
const sampleDocumentPath = "../utils/testdata/sample.txt"

func authedRequest(t *testing.T, client *http.Client, token, method, url string, body []byte) *http.Response {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, reader)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { res.Body.Close() })

	return res
}

// TestDocuments_FullCRUDLifecycle drives a document created via the JSON
// endpoint through list, read and delete, mirroring the notes lifecycle
// test but for the documents resource.
func TestDocuments_FullCRUDLifecycle(t *testing.T) {
	server := testutil.NewServer(t)
	token := testutil.TokenFor(t, "user_a")

	res := authedRequest(t, server.Client(), token, http.MethodPost, server.URL+"/documents",
		[]byte(`{"name":"notes.txt","type":".txt","content":"some document content"}`))
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var created models.Document
	require.NoError(t, json.NewDecoder(res.Body).Decode(&created))
	require.Equal(t, "notes.txt", created.Name)
	require.NotZero(t, created.ID)

	docURL := fmt.Sprintf("%s/documents/%d", server.URL, created.ID)

	res = authedRequest(t, server.Client(), token, http.MethodGet, server.URL+"/documents", nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var documents []models.Document
	require.NoError(t, json.NewDecoder(res.Body).Decode(&documents))
	require.Len(t, documents, 1)

	res = authedRequest(t, server.Client(), token, http.MethodGet, docURL, nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res = authedRequest(t, server.Client(), token, http.MethodDelete, docURL, nil)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	res = authedRequest(t, server.Client(), token, http.MethodGet, docURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

// buildMultipartUpload writes a multipart/form-data body with the given file
// under the "file" field name, matching what DocumentHandler.UploadDocument
// expects from r.FormFile("file").
func buildMultipartUpload(t *testing.T, fieldName string, content []byte) (*bytes.Buffer, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, "sample.txt")
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	return body, writer.FormDataContentType()
}

// TestDocuments_UploadPersistsAndIsListable proves a real multipart file
// upload flows through DocumentHandler -> DocumentService -> SQLite and is
// visible afterwards via GET /documents.
func TestDocuments_UploadPersistsAndIsListable(t *testing.T) {
	server := testutil.NewServer(t)
	token := testutil.TokenFor(t, "user_a")

	content, err := os.ReadFile(sampleDocumentPath)
	require.NoError(t, err)

	body, contentType := buildMultipartUpload(t, "file", content)

	req, err := http.NewRequest(http.MethodPost, server.URL+"/documents/upload", body)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", contentType)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var uploaded models.Document
	require.NoError(t, json.NewDecoder(res.Body).Decode(&uploaded))
	require.Equal(t, filepath.Base(sampleDocumentPath), uploaded.Name)
	require.Equal(t, string(content), uploaded.Content)

	res = authedRequest(t, server.Client(), token, http.MethodGet, server.URL+"/documents", nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var documents []models.Document
	require.NoError(t, json.NewDecoder(res.Body).Decode(&documents))
	require.Len(t, documents, 1)
	require.Equal(t, uploaded.ID, documents[0].ID)
}

// TestDocuments_UploadOverLimitRejected proves a file larger than the 10 MiB
// upload ceiling is rejected with 413 by the real HTTP stack, not just by a
// unit-tested limit that might be wired incorrectly in production.
func TestDocuments_UploadOverLimitRejected(t *testing.T) {
	server := testutil.NewServer(t)
	token := testutil.TokenFor(t, "user_a")

	oversized := bytes.Repeat([]byte("a"), (10<<20)+1) // 10 MiB + 1 byte
	body, contentType := buildMultipartUpload(t, "file", oversized)

	req, err := http.NewRequest(http.MethodPost, server.URL+"/documents/upload", body)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", contentType)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusRequestEntityTooLarge, res.StatusCode)
}
