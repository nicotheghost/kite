package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/kitecloud/kite/kite-service/internal/model"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestTemplateVersionLabels(t *testing.T) {
	c := testDB(t)
	ctx := context.Background()
	_, err := c.DB.Exec(ctx, `
INSERT INTO users (id, email, display_name, discord_id, discord_username, created_at, updated_at)
VALUES ('label_user', 'label@example.com', 'Test', 'label_user', 'test', now(), now()) ON CONFLICT DO NOTHING;
DELETE FROM bundles WHERE id = 'label_bundle';`)
	require.NoError(t, err)

	now := time.Now().UTC()
	b := &model.Bundle{
		ID: "label_bundle", Name: "Test", OwnerUserID: null.StringFrom("label_user"),
		Version: 1, Content: []byte("{}"), CreatedAt: now, UpdatedAt: now,
		Page: model.BundleListing{Name: "Test"},
	}
	require.NoError(t, c.CreateBundle(ctx, b, 10))
	b.Version, b.VersionLabel = 2, "1.01"
	_, err = c.PublishBundleVersion(ctx, b, 1)
	require.NoError(t, err)
	loaded, err := c.Bundle(ctx, b.ID)
	require.NoError(t, err)
	require.Equal(t, "1.01", loaded.VersionLabel)
	require.Equal(t, "1.01", loaded.LatestVersionLabel)
	require.Len(t, loaded.VersionNotes, 1)
	require.Equal(t, "1.01", loaded.VersionNotes[0].VersionLabel)

	_, err = c.DB.Exec(ctx, "UPDATE bundles SET approved_version = 1 WHERE id = $1", b.ID)
	require.NoError(t, err)
	loaded, err = c.Bundle(ctx, b.ID)
	require.NoError(t, err)
	require.Equal(t, 1, loaded.Version)
	require.Empty(t, loaded.VersionLabel)
	require.Equal(t, "1.01", loaded.LatestVersionLabel)
	owned, err := c.BundlesByOwner(ctx, b.OwnerUserID.String)
	require.NoError(t, err)
	require.Len(t, owned, 1)
	require.Equal(t, "1.01", owned[0].VersionLabel)
}
