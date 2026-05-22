package main

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

// TeamConfig holds per-team notification settings.
type TeamConfig struct {
	ChannelName        string `json:"channel_name"`
	Enabled            bool   `json:"enabled"`
	ShowPurpose        bool   `json:"show_purpose"`
	NotifyOnConversion bool   `json:"notify_on_conversion"`
}

// configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration. Any public fields will be deserialized from the Mattermost server configuration
// in OnConfigurationChange.
type configuration struct {
	// TeamConfigs is a JSON-encoded map[string]TeamConfig keyed by team ID.
	TeamConfigs string `json:"TeamConfigs"`
}

// getTeamConfigs parses TeamConfigs into a usable map.
func (c *configuration) getTeamConfigs() (map[string]TeamConfig, error) {
	if c.TeamConfigs == "" {
		return make(map[string]TeamConfig), nil
	}

	var configs map[string]TeamConfig
	if err := json.Unmarshal([]byte(c.TeamConfigs), &configs); err != nil {
		return nil, errors.Wrap(err, "invalid TeamConfigs JSON")
	}

	return configs, nil
}

// Clone shallow copies the configuration.
func (c *configuration) Clone() *configuration {
	clone := *c
	return &clone
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	configuration := new(configuration)

	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	if _, err := configuration.getTeamConfigs(); err != nil {
		return err
	}

	p.setConfiguration(configuration)

	return nil
}
