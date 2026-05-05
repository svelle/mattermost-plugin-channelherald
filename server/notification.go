package main

import (
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// postTypeChannelConverted is the system post type Mattermost creates when a private channel is
// made public. We watch for this in MessageHasBeenPosted to detect conversions.
const postTypeChannelConverted = "system_channel_converted"

// ChannelHasBeenCreated is invoked after a channel is created. We notify the configured channel
// for the team when a new public channel is created.
func (p *Plugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	if channel.Type != model.ChannelTypeOpen {
		return
	}

	p.postChannelNotification(channel, channel.CreatorId, "created")
}

// MessageHasBeenPosted is invoked after a message is posted. We use it to detect when a private
// channel has been converted to public via the system post Mattermost creates for that event.
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	if post.Type != postTypeChannelConverted {
		return
	}

	channel, appErr := p.API.GetChannel(post.ChannelId)
	if appErr != nil {
		p.API.LogError("Failed to get channel for conversion notification", "channel_id", post.ChannelId, "err", appErr.Error())
		return
	}

	if channel.Type != model.ChannelTypeOpen {
		return
	}

	cfg := p.getConfiguration()
	teamConfig, ok := cfg.getTeamConfigs()[channel.TeamId]
	if !ok || !teamConfig.Enabled || !teamConfig.NotifyOnConversion {
		return
	}

	p.postChannelNotification(channel, post.UserId, "made public")
}

// postChannelNotification builds and posts a notification message to the configured channel for
// the team.
func (p *Plugin) postChannelNotification(channel *model.Channel, actorUserID, action string) {
	cfg := p.getConfiguration()
	teamConfig, ok := cfg.getTeamConfigs()[channel.TeamId]
	if !ok || !teamConfig.Enabled || teamConfig.ChannelName == "" {
		return
	}

	notifChannel, appErr := p.API.GetChannelByName(channel.TeamId, teamConfig.ChannelName, false)
	if appErr != nil {
		p.API.LogError("Failed to get notification channel",
			"channel_name", teamConfig.ChannelName,
			"team_id", channel.TeamId,
			"err", appErr.Error(),
		)
		return
	}

	siteURL := p.API.GetConfig().ServiceSettings.SiteURL
	if siteURL == nil || *siteURL == "" {
		p.API.LogError("SiteURL is not configured; cannot build channel link")
		return
	}

	team, appErr := p.API.GetTeam(channel.TeamId)
	if appErr != nil {
		p.API.LogError("Failed to get team", "team_id", channel.TeamId, "err", appErr.Error())
		return
	}

	channelURL := fmt.Sprintf("%s/%s/channels/%s", *siteURL, team.Name, channel.Name)

	var message string
	if actorUserID != "" {
		actor, appErr := p.API.GetUser(actorUserID)
		if appErr != nil {
			p.API.LogError("Failed to get actor user", "user_id", actorUserID, "err", appErr.Error())
			// Fall back to generic message rather than failing entirely.
			message = fmt.Sprintf("A new channel was %s: [%s](%s)", action, channel.DisplayName, channelURL)
		} else {
			message = fmt.Sprintf("@%s just %s a new channel: [%s](%s)", actor.Username, action, channel.DisplayName, channelURL)
		}
	} else {
		message = fmt.Sprintf("A new channel was %s: [%s](%s)", action, channel.DisplayName, channelURL)
	}

	if teamConfig.ShowPurpose && channel.Purpose != "" {
		message += fmt.Sprintf("\n> %s", channel.Purpose)
	}

	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: notifChannel.Id,
		Message:   message,
	}

	if _, err := p.API.CreatePost(post); err != nil {
		p.API.LogError("Failed to create notification post", "err", err.Error())
	}
}
