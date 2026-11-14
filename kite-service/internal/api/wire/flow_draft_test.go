package wire

import (
	"testing"

	"github.com/kitecloud/kite/kite-service/pkg/flow"
	"github.com/stretchr/testify/require"
)

func TestDraftUpdateAllowsUnfinishedBlocks(t *testing.T) {
	unfinished := flow.FlowData{
		Nodes: []flow.FlowNode{{ID: "unfinished"}},
		Edges: []flow.FlowEdge{},
	}

	for _, action := range []string{"save", "publish", "", "invalid"} {
		t.Run("action="+action, func(t *testing.T) {
			command := CommandUpdateRequest{FlowSource: unfinished, DraftAction: action}
			listener := EventListenerUpdateRequest{FlowSource: unfinished, DraftAction: action}
			if action == "save" {
				require.NoError(t, command.Validate())
				require.NoError(t, listener.Validate())
			} else {
				require.Error(t, command.Validate())
				require.Error(t, listener.Validate())
			}
		})
	}
}

func TestDraftUpdateRequiresFlowArrays(t *testing.T) {
	command := CommandUpdateRequest{DraftAction: "save"}
	listener := EventListenerUpdateRequest{DraftAction: "save"}
	require.Error(t, command.Validate())
	require.Error(t, listener.Validate())
}
