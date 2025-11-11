#!/usr/bin/env python3
"""MidCity Utilities sensor for Home Assistant."""
import requests
import sys
import os
import json
import time
import logging
from bs4 import BeautifulSoup
from datetime import datetime

# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Home Assistant Supervisor API configuration
SUPERVISOR_TOKEN = os.environ.get('SUPERVISOR_TOKEN')
HA_URL = 'http://supervisor/core/api'

# MidCity Utilities URLs
LOGIN_URL = "https://buyprepaid.midcityutilities.co.za/ajax/login"
METER_URL = "https://buyprepaid.midcityutilities.co.za/meters"


class MidCityUtilitiesSensor:
    """MidCity Utilities Sensor class."""

    def __init__(self, username, password, scan_interval=300):
        """Initialize the sensor."""
        self.username = username
        self.password = password
        self.scan_interval = scan_interval
        self.session = requests.Session()
        self.headers = {
            'Authorization': f'Bearer {SUPERVISOR_TOKEN}',
            'Content-Type': 'application/json',
        }

    def login(self):
        """Login to MidCity Utilities website."""
        try:
            payload = {
                'email': self.username,
                'password': self.password
            }

            response = self.session.post(LOGIN_URL, data=payload, timeout=30)

            if response.text != '{"ok":true,"success":true}':
                logger.error("Login failed. Please check your username and password.")
                return False

            logger.info("Successfully logged in to MidCity Utilities")
            return True

        except Exception as e:
            logger.error(f"Login error: {e}")
            return False

    def get_meter_data(self):
        """Retrieve meter data from MidCity Utilities."""
        try:
            response = self.session.get(METER_URL, timeout=30)

            if response.status_code != 200:
                logger.error(f"Failed to retrieve meter data. Status code: {response.status_code}")
                return None

            # Parse the HTML content
            soup = BeautifulSoup(response.text, 'html.parser')

            meters = []

            # Find all meter cards - adjust selectors based on actual HTML structure
            meter_cards = soup.find_all('div', class_='meter-card') or soup.find_all('div', class_='card')

            if not meter_cards:
                logger.warning("No meter cards found. Attempting alternative parsing...")
                # Try alternative parsing methods
                meter_cards = soup.find_all('div', {'data-meter-id': True}) or soup.select('.meter, .prepaid-meter')

            for card in meter_cards:
                try:
                    meter_data = self.parse_meter_card(card)
                    if meter_data:
                        meters.append(meter_data)
                except Exception as e:
                    logger.error(f"Error parsing meter card: {e}")
                    continue

            logger.info(f"Retrieved data for {len(meters)} meter(s)")
            return meters

        except Exception as e:
            logger.error(f"Error retrieving meter data: {e}")
            return None

    def parse_meter_card(self, card):
        """Parse individual meter card to extract data."""
        meter_data = {}

        # Try to extract meter number
        meter_number = (
            card.get('data-meter-id') or
            (card.find('span', class_='meter-number') and card.find('span', class_='meter-number').text.strip()) or
            (card.find(string=lambda text: text and 'Meter' in text))
        )

        # Try to extract balance/credit
        balance_elem = (
            card.find('span', class_='balance') or
            card.find('div', class_='credit') or
            card.find(string=lambda text: text and 'R' in str(text))
        )

        if balance_elem:
            balance_text = balance_elem.text.strip() if hasattr(balance_elem, 'text') else str(balance_elem).strip()
            # Extract numeric value from balance text
            import re
            balance_match = re.search(r'R?\s*([\d,.]+)', balance_text)
            if balance_match:
                meter_data['balance'] = float(balance_match.group(1).replace(',', ''))

        # Try to extract meter type (electricity/water)
        meter_type = (
            card.get('data-meter-type') or
            (card.find('span', class_='meter-type') and card.find('span', class_='meter-type').text.strip()) or
            'unknown'
        )

        if meter_number:
            meter_data['meter_number'] = str(meter_number).strip()
            meter_data['meter_type'] = str(meter_type).lower()
            meter_data['last_updated'] = datetime.now().isoformat()

        return meter_data if meter_data else None

    def update_ha_sensor(self, meter_data):
        """Update Home Assistant sensor using the REST API."""
        try:
            for meter in meter_data:
                meter_number = meter.get('meter_number', 'unknown')
                balance = meter.get('balance', 0)
                meter_type = meter.get('meter_type', 'unknown')

                # Create unique entity_id
                entity_id = f"sensor.midcity_{meter_type}_{meter_number.replace(' ', '_').lower()}"

                # Prepare state data
                state_data = {
                    'state': balance,
                    'attributes': {
                        'meter_number': meter_number,
                        'meter_type': meter_type,
                        'unit_of_measurement': 'ZAR',
                        'device_class': 'monetary',
                        'friendly_name': f'MidCity {meter_type.title()} {meter_number}',
                        'last_updated': meter.get('last_updated'),
                        'icon': 'mdi:cash' if meter_type == 'electricity' else 'mdi:water'
                    }
                }

                # Update state via Home Assistant API
                response = requests.post(
                    f'{HA_URL}/states/{entity_id}',
                    headers=self.headers,
                    json=state_data,
                    timeout=10
                )

                if response.status_code in [200, 201]:
                    logger.info(f"Successfully updated sensor: {entity_id}")
                else:
                    logger.error(f"Failed to update sensor {entity_id}: {response.status_code} - {response.text}")

        except Exception as e:
            logger.error(f"Error updating Home Assistant sensor: {e}")

    def run(self):
        """Main run loop."""
        logger.info("Starting MidCity Utilities sensor...")
        logger.info(f"Scan interval: {self.scan_interval} seconds")

        while True:
            try:
                logger.info("Fetching meter data...")

                # Login
                if not self.login():
                    logger.error("Failed to login. Retrying in 60 seconds...")
                    time.sleep(60)
                    continue

                # Get meter data
                meter_data = self.get_meter_data()

                if meter_data:
                    # Update Home Assistant sensors
                    self.update_ha_sensor(meter_data)
                else:
                    logger.warning("No meter data retrieved")

                logger.info(f"Waiting {self.scan_interval} seconds until next update...")
                time.sleep(self.scan_interval)

            except KeyboardInterrupt:
                logger.info("Shutting down...")
                break
            except Exception as e:
                logger.error(f"Unexpected error: {e}")
                time.sleep(60)


def main():
    """Main function."""
    # Read configuration from options.json
    try:
        with open('/data/options.json', 'r') as f:
            config = json.load(f)
    except Exception as e:
        logger.error(f"Failed to read configuration: {e}")
        sys.exit(1)

    username = config.get('username')
    password = config.get('password')
    scan_interval = config.get('scan_interval', 300)

    if not username or not password:
        logger.error("Username and password are required in configuration")
        sys.exit(1)

    # Create and run sensor
    sensor = MidCityUtilitiesSensor(username, password, scan_interval)
    sensor.run()


if __name__ == '__main__':
    main()
