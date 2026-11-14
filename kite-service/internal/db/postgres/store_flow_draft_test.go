package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/kitecloud/kite/kite-service/internal/model"
	"github.com/kitecloud/kite/kite-service/internal/store"
	"github.com/kitecloud/kite/kite-service/pkg/flow"
	"github.com/stretchr/testify/require"
)

func TestFlowDrafts(t *testing.T) {
	c := testDB(t)
	ctx := context.Background()
	_, err := c.DB.Exec(ctx, `
INSERT INTO users (id, email, display_name, discord_id, discord_username, created_at, updated_at)
VALUES ('draft_user', 'draft@example.com', 'Test', 'draft_user', 'test', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO apps (id, name, owner_user_id, creator_user_id, discord_token, discord_id, created_at, updated_at)
VALUES ('draft_app', 'Test', 'draft_user', 'draft_user', 't', 'draft_app', now(), now()) ON CONFLICT DO NOTHING;`)
	require.NoError(t, err)

	now := time.Now().UTC().Truncate(time.Second)
	live := flow.FlowData{Nodes: []flow.FlowNode{{ID: "live"}}}
	draft := flow.FlowData{Nodes: []flow.FlowNode{{ID: "draft"}}}
	save := store.EditorChange("draft_user", "")
	save.SaveDraft = true
	publish := store.EditorChange("draft_user", "")
	publish.PublishDraft = true

	t.Run("command", func(t *testing.T) {
		_, err := c.DB.Exec(ctx, "DELETE FROM commands WHERE id = 'draft_command'")
		require.NoError(t, err)
		cmd, err := c.CreateCommand(ctx, &model.Command{
			ID: "draft_command", Name: "test", AppID: "draft_app", CreatorUserID: "draft_user",
			Enabled: true, FlowSource: live, CreatedAt: now, UpdatedAt: now,
		})
		require.NoError(t, err)
		require.Nil(t, cmd.DraftFlowSource)
		save.BaseVersion = cmd.ContentVersion
		cmd.DraftFlowSource = &draft
		saved, err := c.UpdateCommand(ctx, cmd, save)
		require.NoError(t, err)
		require.Equal(t, live, saved.FlowSource)
		require.Equal(t, draft, *saved.DraftFlowSource)
		require.True(t, saved.UpdatedAt.Equal(now))
		require.True(t, saved.Enabled)
		_, err = c.UpdateCommand(ctx, cmd, save)
		require.ErrorIs(t, err, store.ErrOutdated)

		publish.BaseVersion = saved.ContentVersion
		saved.FlowSource = draft
		saved.UpdatedAt = now.Add(time.Second)
		published, err := c.UpdateCommand(ctx, saved, publish)
		require.NoError(t, err)
		require.Equal(t, draft, published.FlowSource)
		require.Nil(t, published.DraftFlowSource)
		require.True(t, published.Enabled)
	})

	for _, source := range []model.EventSource{model.EventSourceDiscord, model.EventSourceSchedule, model.EventSourceWebhook} {
		t.Run(string(source), func(t *testing.T) {
			id := "draft_" + string(source)
			_, err := c.DB.Exec(ctx, "DELETE FROM event_listeners WHERE id = $1", id)
			require.NoError(t, err)
			listener, err := c.CreateEventListener(ctx, &model.EventListener{
				ID: id, Source: source, AppID: "draft_app", CreatorUserID: "draft_user",
				Enabled: true, FlowSource: live, CreatedAt: now, UpdatedAt: now,
			})
			require.NoError(t, err)
			require.Nil(t, listener.DraftFlowSource)
			save.BaseVersion = listener.ContentVersion
			listener.DraftFlowSource = &draft
			saved, err := c.UpdateEventListener(ctx, listener, save)
			require.NoError(t, err)
			require.Equal(t, live, saved.FlowSource)
			require.Equal(t, draft, *saved.DraftFlowSource)
			require.True(t, saved.UpdatedAt.Equal(now))
			require.True(t, saved.Enabled)
			_, err = c.UpdateEventListener(ctx, listener, save)
			require.ErrorIs(t, err, store.ErrOutdated)

			publish.BaseVersion = saved.ContentVersion
			saved.FlowSource = draft
			saved.UpdatedAt = now.Add(time.Second)
			published, err := c.UpdateEventListener(ctx, saved, publish)
			require.NoError(t, err)
			require.Equal(t, draft, published.FlowSource)
			require.Nil(t, published.DraftFlowSource)
			require.True(t, published.Enabled)
		})
	}
}
