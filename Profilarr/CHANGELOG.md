# Changelog

## 1.1.0

- Changed: Use official Docker image directly (santiagosayshey/profilarr:latest)
- Simplified: No custom build, just use pre-built image
- Limited architectures to amd64 and aarch64 (supported by official image)

## 1.0.7

- Debug: Add directory listing to find correct entry point

## 1.0.6

- Fixed: Start Profilarr directly using node instead of entrypoint script
- Install Node.js and Python dependencies
- Simplified startup process

## 1.0.5

- Fixed: Use su-exec instead of gosu (more compatible with Alpine)
- Added symbolic link for gosu compatibility

## 1.0.4

- Fixed: Install gosu package to support official entrypoint

## 1.0.3

- Fixed: Use COPY --from to extract official Profilarr Docker image
- Simplified configuration (removed PUID/PGID/UMASK as not needed)
- Use official entrypoint from Profilarr image

## 1.0.2

- Fixed: Use Home Assistant base images for proper architecture support

## 1.0.1

- Fixed: Build from source instead of using non-existent Docker image
- Changed to use LinuxServer base image
- Clone and build Profilarr directly from GitHub

## 1.0.0

- Initial release
- Support for Profilarr latest version
- Configurable PUID, PGID, UMASK, and timezone
- Web interface on port 6868
