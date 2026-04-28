package context

import (
	"context"
	"encoding/json"

	"github.com/spf13/cobra"
)

type SettingsContext struct {
	SiteId      string
	RawSettings map[string]json.RawMessage
}

const settingsContextKey = "settings"

func SetSettingsFromContext(cmd *cobra.Command, settings *SettingsContext) context.Context {
	ctx := context.WithValue(cmd.Context(), settingsContextKey, settings)
	return ctx
}

func GetSettingsFromContext(cmd *cobra.Command) (*SettingsContext, bool) {
	settings := cmd.Context().Value(settingsContextKey)
	if settings == nil {
		return nil, false
	}
	s, ok := settings.(*SettingsContext)
	return s, ok
}
