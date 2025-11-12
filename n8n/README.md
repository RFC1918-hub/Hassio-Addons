# Home Assistant Add-on: n8n

## About

n8n is a fair-source-licensed workflow automation tool. It helps you to connect different services and build automated workflows using a visual editor.

## Features

- Visual workflow editor
- 200+ nodes to automate tasks
- Self-hosted
- Real-time workflow execution
- Always runs the latest n8n version from official Docker image

## Installation

1. Add this repository to your Home Assistant instance
2. Install the n8n add-on
3. Start the add-on
4. Click "OPEN WEB UI" to open the n8n interface (port 5678)

## Configuration

### Options

- `n8n_host` (optional): The hostname or domain where n8n will be accessible (e.g., "n8n.example.com" or "192.168.1.100:5678")
- `n8n_protocol` (optional): Protocol to use - either "http" or "https" (default: "http")
- `timezone` (optional): Timezone for n8n (default: "UTC")

### Example configuration:

```yaml
n8n_host: "n8n.example.com"
n8n_protocol: "https"
timezone: "America/New_York"
```

If `n8n_host` is not set, n8n will use the local access URL. Set this to your external domain/hostname if you want webhooks and external integrations to work properly.

## Support

Got questions? Feel free to [open an issue on our GitHub repository](https://github.com/RFC1918-hub/Hassio-Add-ons/issues).