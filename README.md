# Channel Herald

[![Build Status](https://github.com/svelle/mattermost-plugin-channelherald/actions/workflows/ci.yml/badge.svg)](https://github.com/svelle/mattermost-plugin-channelherald/actions/workflows/ci.yml)

Channel Herald posts a notification whenever a new public channel is created in Mattermost. Notifications are sent by a dedicated bot and include a clickable link to the new channel, the creator's username, and optionally the channel's purpose. Each team can be configured independently.

## Features

- Notifies a chosen channel whenever a new public channel is created in a team
- Notification message format: `@username just created a new channel: [Channel Name](url)`
- Optionally includes the channel purpose beneath the message
- Optionally notifies when a private channel is made public
- Configured per-team — teams without a configured notification channel receive no notifications
- Uses absolute URLs so links work in email and mobile notifications

## Installation

1. Download the latest release from the [releases page](https://github.com/svelle/mattermost-plugin-channelherald/releases).
2. In Mattermost, go to **System Console → Plugins → Plugin Management** and upload the `.tar.gz` file.
3. Enable the plugin.

## Configuration

Go to **System Console → Plugins → Channel Herald**.

For each team you will see:

| Setting | Description |
|---|---|
| **Enable notifications for this team** | Turn notifications on or off for the team. |
| **Notification channel name** | The URL slug of the channel to post notifications in (e.g. `town-square`). The channel must already exist in that team. |
| **Include channel purpose in the notification** | Appends the new channel's purpose as a blockquote beneath the notification message. |
| **Also notify when a private channel is made public** | Posts a notification when an existing private channel is converted to public. |

Save the settings. The plugin will immediately start watching for new channels.

> **Tip:** Create a dedicated channel such as `new-channels` for notifications to keep them easy to find and mute if needed.

## Notification examples

New channel created:

> @alice just created a new channel: [Marketing 2026](https://your.server/team-name/channels/marketing-2026)

With purpose enabled:

> @alice just created a new channel: [Marketing 2026](https://your.server/team-name/channels/marketing-2026)
> Campaign planning and asset tracking for the 2026 fiscal year.

Private channel made public:

> @bob just made a new channel public: [Engineering Roadmap](https://your.server/team-name/channels/engineering-roadmap)

## Development

### Prerequisites

- Go 1.22+
- Node 20+ / npm 10+
- A running Mattermost server

### Building

```bash
make
```

This produces `dist/com.mattermost.plugin-channelherald.tar.gz` ready for upload.

### Deploying to a local server

Enable plugin uploads and local mode in your Mattermost server config:

```json
{
    "PluginSettings": {
        "EnableUploads": true
    },
    "ServiceSettings": {
        "EnableLocalMode": true,
        "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
    }
}
```

Then deploy:

```bash
make deploy
```

To watch for changes and redeploy automatically:

```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=<your-token>
make watch
```

### Running tests

```bash
make test
```

### Releasing

Tag the commit and push — CI will build and attach the plugin bundle to the GitHub release automatically.

```bash
make patch   # bump patch version and tag
make minor   # bump minor version and tag
make major   # bump major version and tag
```
