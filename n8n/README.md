# Home Assistant Add-on: n8n

## About

n8n is a fair-source-licensed workflow automation tool. It helps you to connect different services and build automated workflows using a visual editor.

## Features

- Visual workflow editor
- 200+ nodes to automate tasks
- Self-hosted
- Real-time workflow execution
- Version 3.0.0

## Installation

1. Add this repository to your Home Assistant instance
2. Install the n8n add-on
3. Start the add-on
4. Click "OPEN WEB UI" to open the n8n interface (port 5678)

## Configuration

The add-on comes with a default configuration:

```yaml
environment:
  GENERIC_TIMEZONE: "Africa/Johannesburg"
  TZ: "Africa/Johannesburg"
  N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS: "true"
  N8N_RUNNERS_ENABLED: "true"
```

## Support

Got questions? Feel free to [open an issue on our GitHub repository](https://github.com/RFC1918-hub/Hassio-Add-ons/issues).