# MidCity Utilities Sensor Add-on for Home Assistant

Monitor your MidCity Utilities prepaid meters directly in Home Assistant with native sensor entities.

## About

This add-on allows you to monitor your MidCity Utilities prepaid meters (electricity and water) by scraping data from the MidCity Utilities website and creating native Home Assistant sensor entities.

## Features

- Native Home Assistant sensor entities (no MQTT required)
- Automatic meter discovery
- Configurable scan interval
- Support for multiple meters
- Displays balance in South African Rand (ZAR)

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
```

### Configuration Options

| Option | Required | Default | Description |
|--------|----------|---------|-------------|
| `username` | Yes | - | Your MidCity Utilities email/username |
| `password` | Yes | - | Your MidCity Utilities password |
| `scan_interval` | No | 300 | Update interval in seconds (minimum 60) |

## Usage

1. Configure the add-on with your credentials
2. Start the add-on
3. Check the logs to ensure successful connection
4. Sensors will be automatically created with entity IDs like:
   - `sensor.midcity_electricity_<meter_number>`
   - `sensor.midcity_water_<meter_number>`

## Sensor Entities

The add-on creates sensor entities with the following attributes:

- **State**: Current balance in ZAR
- **Attributes**:
  - `meter_number`: Your meter number
  - `meter_type`: Type of meter (electricity/water)
  - `last_updated`: Timestamp of last update
  - `unit_of_measurement`: ZAR
  - `device_class`: monetary

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

### Version 1.0.0
- Initial release
- Native Home Assistant sensor integration
- Support for electricity and water meters
- Configurable scan interval

## Credits

This add-on scrapes data from MidCity Utilities prepaid portal at https://buyprepaid.midcityutilities.co.za/

## License

MIT License
