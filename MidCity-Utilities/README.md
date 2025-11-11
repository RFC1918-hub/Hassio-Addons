# MidCity Utilities Sensor Add-on for Home Assistant

Monitor your MidCity Utilities prepaid meters directly in Home Assistant with native sensor entities.

## About

This add-on allows you to monitor your MidCity Utilities prepaid meters (electricity and water) by scraping data from the MidCity Utilities website and creating native Home Assistant sensor entities.

## Features

- Native Home Assistant sensor entities via MQTT Discovery
- Automatic meter discovery
- Entities with unique_id for full UI management
- Automatic device grouping - sensors appear under "MidCity Utilities Sensor" device
- Configurable scan interval
- Support for multiple meters
- Displays balance in kWh for electricity meters
- Extracts predicted zero balance date
- Full UI customization support (rename, change icon, assign to area)

## Installation

### Method 1: Add Repository (Recommended)

1. In Home Assistant, navigate to **Settings** → **Add-ons** → **Add-on Store**
2. Click the **⋮** (three dots) in the top right corner
3. Select **Repositories**
4. Add this repository URL: `https://github.com/Hassio-Addons/MidCity-Utilities`
5. Click **Add**
6. Find "MidCity Utilities Sensor" in the add-on store and click **Install**

### Method 2: Manual Installation

1. Copy this folder to `/addons/midcity_utilities/` on your Home Assistant host
2. Restart Home Assistant
3. Navigate to **Settings** → **Add-ons** → **Add-on Store**
4. Refresh the page
5. Find "MidCity Utilities Sensor" and click **Install**

## Configuration

After installation, configure the add-on with your MidCity Utilities credentials:

```yaml
username: your_email@example.com
password: your_password
scan_interval: 300
log_level: info
```

### Configuration Options

| Option | Required | Default | Description |
|--------|----------|---------|-------------|
| `username` | Yes | - | Your MidCity Utilities email/username |
| `password` | Yes | - | Your MidCity Utilities password |
| `scan_interval` | No | 300 | Update interval in seconds (minimum 60) |
| `log_level` | No | info | Log level: debug, info, warning, or error |

### Prerequisites

This add-on requires the **Mosquitto broker** add-on to be installed and running for MQTT Discovery. If you don't have it:

1. Go to **Settings** → **Add-ons** → **Add-on Store**
2. Search for "Mosquitto broker"
3. Click **Install**
4. Start the Mosquitto broker add-on
5. Then configure and start this add-on

## Usage

1. Install and start the Mosquitto broker add-on (if not already installed)
2. Configure this add-on with your MidCity Utilities credentials
3. Start the add-on
4. Check the logs to ensure successful connection
5. Sensors will be automatically created via MQTT Discovery with entity IDs like:
   - `sensor.midcity_electricity_<meter_number>`
   - `sensor.midcity_water_<meter_number>`
6. Sensors automatically appear under the "MidCity Utilities Sensor" device
7. All sensors are fully manageable from the UI (Settings → Devices & Services → Entities)

## Sensor Entities

The add-on creates sensor entities with the following properties:

- **State**: Current balance (kWh for electricity, m³ for water)
- **Unique ID**: Automatically assigned for UI management
- **Device**: Grouped under "MidCity Utilities Sensor"
- **Attributes**:
  - `meter_number`: Your meter number
  - `meter_type`: Type of meter (electricity/water)
  - `last_updated`: Timestamp of last update
  - `predicted_zero_date`: Date when balance expected to reach zero
  - `attribution`: "Data from MidCity Utilities"
- **Properties**:
  - `unit_of_measurement`: kWh (electricity), m³ (water), or ZAR
  - `device_class`: energy (electricity), water (water), or monetary
  - `icon`: Lightning bolt (electricity), water drop (water)

## Example Automation

```yaml
automation:
  - alias: "Low Electricity Alert"
    trigger:
      - platform: numeric_state
        entity_id: sensor.midcity_electricity_12345678
        below: 50
    action:
      - service: notify.mobile_app
        data:
          message: "Electricity balance is low: R{{ states('sensor.midcity_electricity_12345678') }}"
```

## Troubleshooting

### Check Logs

Always check the add-on logs for error messages:
1. Go to **Settings** → **Add-ons** → **MidCity Utilities Sensor**
2. Click on the **Log** tab

### Common Issues

#### Login Failed
- Verify your username and password are correct
- Check if you can log in at https://buyprepaid.midcityutilities.co.za/

#### No Sensors Created
- Check the add-on logs for errors
- Ensure the add-on has proper permissions
- Verify that meters are visible when logged into the website

#### Sensors Not Updating
- Check your scan_interval setting
- Verify internet connection
- Review logs for any error messages

## Support

If you encounter issues:
1. Check the logs for error messages
2. Verify your credentials work on the MidCity Utilities website
3. Open an issue on GitHub with relevant log entries

## Changelog

### Version 1.2.0
- **BREAKING:** Switched to MQTT Discovery for proper entity creation
- Entities now have unique_id for full UI management
- Automatic device grouping under "MidCity Utilities Sensor"
- No manual token configuration required
- Sensors fully customizable from UI

### Version 1.0.0
- Initial release
- Native Home Assistant sensor integration
- Support for electricity and water meters
- Configurable scan interval

## Credits

This add-on scrapes data from MidCity Utilities prepaid portal at https://buyprepaid.midcityutilities.co.za/

## License

MIT License
