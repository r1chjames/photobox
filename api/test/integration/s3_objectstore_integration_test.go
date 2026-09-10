//go:build integration

package integration

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	objS3 "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/object/s3"
	wsstore "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/workspace"
)

// TestS3ObjectStore_WorkspaceIsolation_Integration exercises the real S3
// adapter + workspace-bound accessor against the configured bucket
// (Wasabi/S3/MinIO). Skipped unless S3_* env vars are set, so the default
// integration run needs no credentials.
//
// It proves the two properties that matter for tenancy:
//   - writes are readable back through the same workspace-bound store, and
//   - a store bound to workspace A cannot read workspace B's object even
//     though both hold the same photo ID.
func TestS3ObjectStore_WorkspaceIsolation_Integration(t *testing.T) {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")
	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		t.Skip("S3_ENDPOINT/S3_ACCESS_KEY/S3_SECRET_KEY/S3_BUCKET not set; skipping live S3 test")
	}

	store, err := objS3.New(objS3.Config{
		Endpoint:  endpoint,
		Region:    os.Getenv("S3_REGION"),
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		UseSSL:    os.Getenv("S3_USE_SSL") != "false",
	})
	require.NoError(t, err, "connecting to S3 endpoint %s", endpoint)

	const (
		wsA = "aaaaaaaa-1111-1111-1111-111111111111"
		wsB = "bbbbbbbb-2222-2222-2222-222222222222"
	)
	sA, err := wsstore.New(store, wsA)
	require.NoError(t, err)
	sB, err := wsstore.New(store, wsB)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Unique photo ID per run so repeated runs cannot collide.
	photoID := "itest-" + time.Now().UTC().Format("20060102150405.000000000")
	// photoID must be path-safe for the key builder.
	photoID = strings.ReplaceAll(photoID, ".", "-")

	payload := []byte("integration-test-object-" + photoID)
	require.NoError(t, sA.PutOriginal(ctx, photoID, bytes.NewReader(payload), int64(len(payload)), "application/octet-stream"))

	// Cleanup on exit.
	defer func() {
		_ = sA.DeleteOriginal(context.Background(), photoID)
	}()

	// Same workspace reads it back.
	rc, err := sA.GetOriginal(ctx, photoID)
	require.NoError(t, err, "workspace A should read its own object")
	got, err := io.ReadAll(rc)
	rc.Close()
	require.NoError(t, err)
	assert.Equal(t, payload, got, "round-tripped bytes must match")

	// Cross-workspace: store B builds a different key, so the same photo ID
	// must resolve to nothing for workspace B. This is the structural
	// isolation guarantee.
	rcB, err := sB.GetOriginal(ctx, photoID)
	if err == nil {
		_, readErr := io.ReadAll(rcB)
		rcB.Close()
		require.Error(t, readErr, "workspace B must not be able to read workspace A's object")
	} else {
		// Some S3 implementations surface the miss at Get time rather than read time.
		assert.Error(t, err)
	}

	exists, err := store.Exists(ctx, "originals/"+wsB+"/"+photoID)
	require.NoError(t, err)
	assert.False(t, exists, "workspace B must not see workspace A's object under its own prefix")
}
