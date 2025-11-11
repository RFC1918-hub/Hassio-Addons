# Jellyseerr - Media Request Manager

Jellyseerr is a free and open source software application for managing requests for your media library. It is a fork of Overseerr built to bring support for Jellyfin & Emby media servers!

## Features

- Full Jellyfin/Emby/Plex integration
- Supports Movies, TV Shows, and Music
- User management and request system
- Mobile friendly interface
- Granular permission system
- Integration with Sonarr and Radarr
- Support for multiple media servers
- Email & webhook notifications
- Local and external user support

## Configuration

### Persistence
All configuration data is stored in `/data` which maps to the add-on's configuration folder, ensuring persistence across restarts and updates.

### First Time Setup
1. After installation, access the web interface at `http://homeassistant.local:5055`
2. Follow the setup wizard to configure your Jellyfin/Emby/Plex server
3. Configure Sonarr and Radarr integration if desired
4. Set up users and permissions

### Environment Variables
The following environment variables are pre-configured:
- `TZ`: Timezone (set to Africa/Johannesburg by default, modify in config.yaml if needed)
- `CONFIG_DIR`: Configuration directory (automatically set to /data)

## Support

For more information, visit the [official documentation](https://docs.seerr.dev/).

## Image Source
This add-on uses the official Jellyseerr Docker image from GitHub Container Registry:
`ghcr.io/fallenbagel/jellyseerr:latest`
