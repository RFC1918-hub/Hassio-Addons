# Changelog

All notable changes to this project will be documented in this file.

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
