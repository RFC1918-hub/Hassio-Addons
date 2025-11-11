# Changelog

All notable changes to this project will be documented in this file.

## [1.0.3] - 2025-11-11

### Fixed
- Exclude transaction history table rows from meter parsing
- Focus on actual meter display elements (panels, cards, divs)
- Filter out "Product Type" and "Download Invoice" elements

### Improved
- Enhanced page structure analysis when meters not found
- Show headings, panels, forms, and currency elements
- Better debugging information for troubleshooting HTML structure

## [1.0.2] - 2025-11-11

### Added
- Configurable log level (debug, info, warning, error)
- Enhanced HTML parsing with 5 different strategies
- Debug logging to save HTML to /tmp/meters_page.html
- More robust meter number detection (8-12 digit patterns)
- Flexible balance extraction (handles various currency formats)
- Smart meter type detection (electricity, water, gas)

### Improved
- Better error messages showing what HTML structure was found
- Parse meter data from tables, divs, cards, and custom elements
- Log detailed parsing information for troubleshooting
- Default log level set to DEBUG for initial setup

## [1.0.1] - 2025-11-11

### Fixed
- Fixed Docker image registry error by removing external image reference
- Fixed bashio unbound variable error in run.sh script
- Simplified run.sh to use standard bash instead of bashio functions
- Added explicit bash package to Dockerfile
- Add-on now builds locally instead of pulling from GitHub Container Registry

## [1.0.0] - 2025-11-11

### Added
- Initial release of MidCity Utilities Sensor add-on
- Native Home Assistant sensor integration using Supervisor API
- Automatic meter discovery from MidCity Utilities portal
- Support for multiple meters (electricity and water)
- Configurable scan interval
- Multi-architecture support (aarch64, amd64, armhf, armv7, i386)
- Comprehensive logging for troubleshooting
- Balance displayed in South African Rand (ZAR)

### Features
- Scrapes meter data from buyprepaid.midcityutilities.co.za
- Creates native HA sensors (no MQTT required)
- Automatic state updates based on scan interval
- Proper error handling and retry logic
