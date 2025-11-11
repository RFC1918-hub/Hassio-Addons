# FlareSolverr - Cloudflare Proxy

FlareSolverr is a proxy server to bypass Cloudflare and DDoS-GUARD protection. It is commonly used with applications like Prowlarr, Jackett, and other *arr applications to bypass Cloudflare protection on indexers.

## Features

- Bypass Cloudflare protection
- Bypass DDoS-GUARD protection
- RESTful API interface
- Session support for persistent cookies
- Automatic challenge solving
- Proxy support
- Prometheus metrics (optional)
- Headless browser automation

## How it works

FlareSolverr starts a proxy server and waits for user requests. When a request arrives, it uses Selenium with undetected-chromedriver to create a web browser (Chrome). It opens the URL and waits until the Cloudflare challenge is solved (or timeout). The HTML code and cookies are sent back to the user, which can then be used to bypass Cloudflare using other HTTP clients.

## Configuration

### First Time Setup
1. After installation, the API will be available at `http://homeassistant.local:8191`
2. Configure your applications (Prowlarr, Jackett, etc.) to use FlareSolverr
3. In Prowlarr: Settings → Indexers → Add FlareSolverr with URL `http://homeassistant.local:8191`

### API Usage

Example request to solve a Cloudflare challenge:

```bash
curl -L -X POST 'http://homeassistant.local:8191/v1' \
-H 'Content-Type: application/json' \
--data-raw '{
  "cmd": "request.get",
  "url": "http://example.com/",
  "maxTimeout": 60000
}'
```

## Options

### LOG_LEVEL
Verbosity of the logging.

**Default**: `info`  
**Options**: `debug`, `info`, `warning`, `error`

### LOG_FILE
Path to capture log to file.

**Default**: None  
**Example**: `/config/flaresolverr.log`

### LOG_HTML
Only for debugging. If true, all HTML that passes through the proxy will be logged to the console in debug level.

**Default**: `false`

### PROXY_URL
URL for proxy. Can be overwritten per request.

**Default**: None  
**Example**: `http://127.0.0.1:8080`

### PROXY_USERNAME
Username for proxy authentication.

**Default**: None

### PROXY_PASSWORD
Password for proxy authentication.

**Default**: None

### CAPTCHA_SOLVER
Captcha solving method (currently experimental).

**Default**: None

### TZ
Timezone used in the logs and the web browser.

**Default**: `UTC`  
**Example**: `Europe/London`

### LANG
Language used in the web browser.

**Default**: None  
**Example**: `en_GB`

### HEADLESS
Run the web browser in headless mode (recommended) or visible mode (for debugging).

**Default**: `true`

### TEST_URL
FlareSolverr makes a request on start to verify the browser is working. Change this if it's blocked in your country.

**Default**: `https://www.google.com`

### PROMETHEUS_ENABLED
Enable Prometheus metrics exporter.

**Default**: `false`

### PROMETHEUS_PORT
Listening port for Prometheus exporter.

**Default**: `8192`

## API Commands

### sessions.create
Create a persistent browser session that retains cookies.

### sessions.list
List all active sessions.

### sessions.destroy
Destroy a browser session to free resources.

### request.get
Make a GET request and solve any Cloudflare challenges.

### request.post
Make a POST request and solve any Cloudflare challenges.

## Integration with *arr Applications

### Prowlarr
1. Go to Settings → Indexers
2. Click "Add FlareSolverr"
3. Enter the URL: `http://homeassistant.local:8191`
4. Click "Test" and then "Save"
5. Configure indexers to use FlareSolverr tags

### Jackett
1. Go to Jackett Settings
2. Find "FlareSolverr API URL"
3. Enter: `http://homeassistant.local:8191`
4. Click "Apply server settings"

## Performance Notes

Web browsers consume significant memory. Avoid making many simultaneous requests on low-RAM systems. Each request launches a new browser instance unless using persistent sessions.

When using sessions, always close them when done to free resources.

## Support

For more information and API documentation, visit:
- [FlareSolverr GitHub](https://github.com/FlareSolverr/FlareSolverr)
- [API Documentation](https://github.com/FlareSolverr/FlareSolverr#usage)

## Image Source
This add-on uses the official FlareSolverr Docker image:
`ghcr.io/flaresolverr/flaresolverr:latest`
