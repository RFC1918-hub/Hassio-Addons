# Profilarr Home Assistant Add-on

![Supports aarch64 Architecture][aarch64-shield]
![Supports amd64 Architecture][amd64-shield]
![Supports armhf Architecture][armhf-shield]
![Supports armv7 Architecture][armv7-shield]
![Supports i386 Architecture][i386-shield]

A Sonarr/Radarr companion app for automating the setup of Custom Formats, Quality Profiles, and Naming Schemes.

## About

Profilarr is a powerful tool that helps you manage and automate your Sonarr and Radarr instances. It simplifies the configuration of:

- Custom Formats
- Quality Profiles
- Naming Schemes

This add-on packages Profilarr for easy installation on Home Assistant.

## Installation

1. Add this repository to your Home Assistant add-on store
2. Install the "Profilarr" add-on
3. Configure the add-on options (optional)
4. Start the add-on
5. Access the web interface on port 6868

## Configuration

```yaml
TZ: "UTC"
```

### Option: `TZ`

Your timezone (e.g., `America/New_York`, `Europe/London`). Default is `UTC`.

## Usage

1. After starting the add-on, access the web interface at `http://homeassistant.local:6868`
2. Configure your Sonarr/Radarr instances
3. Set up your custom formats and quality profiles
4. Let Profilarr automate the rest!

## Data Storage

Configuration data is stored in `/config/profilarr/` on your Home Assistant instance.

## Support

For issues with the add-on, please open an issue on the [GitHub repository][github].

For issues with Profilarr itself, please visit the [official Profilarr repository][profilarr].

## License

MIT License

[aarch64-shield]: https://img.shields.io/badge/aarch64-yes-green.svg
[amd64-shield]: https://img.shields.io/badge/amd64-yes-green.svg
[armhf-shield]: https://img.shields.io/badge/armhf-yes-green.svg
[armv7-shield]: https://img.shields.io/badge/armv7-yes-green.svg
[i386-shield]: https://img.shields.io/badge/i386-yes-green.svg
[github]: https://github.com/RFC1918-hub/Hassio-Addons
[profilarr]: https://github.com/santiagosayshey/profilarr
