# Profilarr - Profile Manager

Profilarr is a profile management tool for your *arr applications, allowing you to manage and sync quality profiles across Radarr, Sonarr, and other arr services.

## Features

- Centralized profile management for *arr applications
- Quality profile synchronization
- Custom format management
- Profile templates and presets
- Web-based interface
- Support for Radarr, Sonarr, and other arr services

## Configuration

### Persistence
All configuration data is stored in `/data` which maps to the add-on's configuration folder, ensuring persistence across restarts and updates.

### First Time Setup
1. After installation, access the web interface at `http://homeassistant.local:6868`
2. Configure your arr service connections (Radarr, Sonarr, etc.)
3. Set up your quality profiles and custom formats
4. Sync profiles across your services

### Environment Variables
The following environment variables can be configured in the add-on configuration:
- `PUID`: User ID (defaults to 1000)
- `PGID`: Group ID (defaults to 1000)
- `UMASK`: File creation mask (defaults to 022)
- `TZ`: Timezone (defaults to UTC)

## Options

### PUID
The user ID that Profilarr should run as.

**Default**: `1000`

### PGID
The group ID that Profilarr should run as.

**Default**: `1000`

### UMASK
The umask for file creation.

**Default**: `"022"`

### TZ
The timezone for the application.

**Default**: `UTC`

## Support

For more information, visit the [Profilarr GitHub repository](https://github.com/Dictionarry-Hub/profilarr).

## Image Source
This add-on uses the official Profilarr Docker image:
`santiagosayshey/profilarr:latest`
