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
log_level = os.environ.get('LOG_LEVEL', 'INFO').upper()
logging.basicConfig(
    level=getattr(logging, log_level, logging.INFO),
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Home Assistant Supervisor API configuration
SUPERVISOR_TOKEN = os.environ.get('SUPERVISOR_TOKEN')
HA_URL = 'http://supervisor/core/api'

# Check if supervisor token is available
if SUPERVISOR_TOKEN:
    logger.info("SUPERVISOR_TOKEN found")
else:
    logger.warning("SUPERVISOR_TOKEN not found - API calls may fail")

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

            # Save HTML for debugging
            try:
                with open('/tmp/meters_page.html', 'w', encoding='utf-8') as f:
                    f.write(response.text)
                logger.info("Saved HTML to /tmp/meters_page.html for debugging")
            except Exception as e:
                logger.warning(f"Could not save HTML for debugging: {e}")

            # Parse the HTML content
            soup = BeautifulSoup(response.text, 'html.parser')

            meters = []

            # Log some info about the page structure
            logger.info(f"Page title: {soup.title.string if soup.title else 'No title'}")

            # Try to find common container elements
            logger.info("Searching for meter data in HTML...")

            # Try multiple strategies to find meters (NOT transaction history)
            meter_cards = []

            # Strategy 1: Look for panel or well sections (common in Bootstrap)
            meter_cards = soup.find_all('div', class_=['panel', 'well', 'panel-body'])
            if meter_cards:
                logger.info(f"Found {len(meter_cards)} panel/well elements")
                # Filter out elements that look like transaction history
                meter_cards = [card for card in meter_cards if 'Product Type' not in card.get_text()]

            # Strategy 2: Look for cards (but not in tables)
            if not meter_cards:
                all_cards = soup.find_all('div', class_='card')
                # Exclude cards that are inside tables or contain transaction data
                meter_cards = []
                for card in all_cards:
                    text = card.get_text()
                    # Skip if it's transaction history
                    if 'Product Type' not in text and 'Download Invoice' not in text:
                        meter_cards.append(card)
                if meter_cards:
                    logger.info(f"Found {len(meter_cards)} card elements (excluding transaction history)")

            # Strategy 3: Look for container divs with specific classes
            if not meter_cards:
                meter_cards = soup.find_all('div', class_=['meter-info', 'meter-display', 'account-info'])
                if meter_cards:
                    logger.info(f"Found {len(meter_cards)} meter info containers")

            # Strategy 4: Look for divs that contain balance/credit info but aren't tables
            if not meter_cards:
                # Find all divs, but exclude those in tables
                all_divs = soup.find_all('div')
                meter_cards = []
                for div in all_divs:
                    # Skip if parent is a table
                    if div.find_parent('table'):
                        continue
                    text = div.get_text()
                    # Look for divs with meter numbers or balance info
                    if any(keyword in text.lower() for keyword in ['meter', 'balance', 'credit', 'account']):
                        # But not transaction history
                        if 'Product Type' not in text and 'Download Invoice' not in text:
                            meter_cards.append(div)

                # Limit to reasonable sized elements (not too small, not too large)
                meter_cards = [c for c in meter_cards if 50 < len(c.get_text()) < 500]
                if meter_cards:
                    logger.info(f"Found {len(meter_cards)} divs with meter-related content")

            # Strategy 5: Look for the main content area
            if not meter_cards:
                main_content = soup.find('div', class_=['container', 'content', 'main-content'])
                if main_content:
                    # Get direct children that might be meter cards
                    meter_cards = main_content.find_all('div', recursive=False)
                    if meter_cards:
                        logger.info(f"Found {len(meter_cards)} direct children of main content")

            if not meter_cards:
                logger.warning("No meter cards found with any strategy")
                # Log useful page structure info for debugging
                logger.info("=== Analyzing page structure for meter data ===")

                # Look for all headings
                headings = soup.find_all(['h1', 'h2', 'h3', 'h4'])
                if headings:
                    logger.info(f"Found {len(headings)} headings:")
                    for h in headings[:10]:
                        logger.info(f"  - {h.name}: {h.get_text(strip=True)[:100]}")

                # Look for panels or sections
                panels = soup.find_all('div', class_=['panel', 'well', 'section'])
                if panels:
                    logger.info(f"Found {len(panels)} panels/sections")
                    for i, panel in enumerate(panels[:5]):
                        text = panel.get_text(separator=' ', strip=True)[:200]
                        logger.info(f"  Panel {i+1}: {text}")

                # Look for forms (meters might be in a form)
                forms = soup.find_all('form')
                if forms:
                    logger.info(f"Found {len(forms)} forms")
                    for i, form in enumerate(forms):
                        logger.info(f"  Form {i+1} action: {form.get('action', 'N/A')}")

                # Look for spans/divs with 'R' (currency) outside tables
                currency_elements = []
                for elem in soup.find_all(['span', 'div', 'p']):
                    if elem.find_parent('table'):
                        continue
                    text = elem.get_text(strip=True)
                    if text.startswith('R ') or text.startswith('R\xa0'):
                        currency_elements.append(elem)

                if currency_elements:
                    logger.info(f"Found {len(currency_elements)} currency elements outside tables:")
                    for elem in currency_elements[:10]:
                        parent_class = elem.parent.get('class', ['no-class'])
                        logger.info(f"  - {elem.get_text(strip=True)} (parent: {parent_class})")

                logger.info("=== End page structure analysis ===")
                return None

            for card in meter_cards:
                try:
                    meter_data = self.parse_meter_card(card)
                    if meter_data:
                        meters.append(meter_data)
                        logger.info(f"Successfully parsed meter: {meter_data}")
                except Exception as e:
                    logger.error(f"Error parsing meter card: {e}")
                    continue

            logger.info(f"Retrieved data for {len(meters)} meter(s)")
            return meters if meters else None

        except Exception as e:
            logger.error(f"Error retrieving meter data: {e}")
            return None

    def parse_meter_card(self, card):
        """Parse individual meter card to extract data."""
        import re
        meter_data = {}

        # Get all text from the card
        card_text = card.get_text(separator=' ', strip=True)
        logger.debug(f"Parsing card text: {card_text[:200]}")

        # Try to extract meter number - look for patterns like digits
        meter_number = None

        # Strategy 1: Look for data attributes
        meter_number = card.get('data-meter-id') or card.get('data-meter')

        # Strategy 2: Look for specific elements
        if not meter_number:
            meter_elem = (
                card.find('span', class_='meter-number') or
                card.find('div', class_='meter-number') or
                card.find(class_=re.compile(r'meter.*number', re.I))
            )
            if meter_elem:
                meter_number = meter_elem.text.strip()

        # Strategy 3: Look for patterns in text (meter numbers are typically 8-12 digits)
        if not meter_number:
            meter_match = re.search(r'\b(\d{8,12})\b', card_text)
            if meter_match:
                meter_number = meter_match.group(1)

        # Strategy 4: Look for text containing "Meter" followed by a number
        if not meter_number:
            meter_match = re.search(r'Meter\s*:?\s*(\d+)', card_text, re.I)
            if meter_match:
                meter_number = meter_match.group(1)

        # Try to extract balance/credit - look for specific patterns with units
        balance = None
        unit = None

        # Strategy 1: Look for balance with kWh (electricity)
        balance_match = re.search(r'balance[:\s]+([\d,]+\.?\d*)\s*kWh', card_text, re.I)
        if balance_match:
            balance_str = balance_match.group(1).replace(',', '')
            unit = 'kWh'
            try:
                balance = float(balance_str)
                meter_data['balance'] = balance
                logger.debug(f"Found kWh balance: {balance}")
            except ValueError:
                pass

        # Strategy 2: Look for balance with L or m³ (water)
        if not balance:
            balance_match = re.search(r'balance[:\s]+([\d,]+\.?\d*)\s*([Ll]|m³|m3)', card_text, re.I)
            if balance_match:
                balance_str = balance_match.group(1).replace(',', '')
                unit = 'm³' if 'm' in balance_match.group(2).lower() else 'L'
                try:
                    balance = float(balance_str)
                    meter_data['balance'] = balance
                    logger.debug(f"Found water balance: {balance} {unit}")
                except ValueError:
                    pass

        # Strategy 3: Look for balance with R (South African Rand)
        if not balance:
            balance_match = re.search(r'R\s*([\d,]+\.?\d*)', card_text)
            if balance_match:
                balance_str = balance_match.group(1).replace(',', '')
                # Make sure it's not the meter number (meter numbers don't have decimals typically)
                if '.' in balance_str or float(balance_str) < 100000:
                    unit = 'ZAR'
                    try:
                        balance = float(balance_str)
                        meter_data['balance'] = balance
                        logger.debug(f"Found ZAR balance: {balance}")
                    except ValueError:
                        pass

        # Strategy 4: Look for any decimal number with "balance" nearby (not meter ID)
        if not balance:
            balance_match = re.search(r'balance[:\s]+([\d,]+\.\d+)', card_text, re.I)
            if balance_match:
                balance_str = balance_match.group(1).replace(',', '')
                try:
                    balance = float(balance_str)
                    meter_data['balance'] = balance
                    logger.debug(f"Found generic balance: {balance}")
                except ValueError:
                    pass

        # Try to extract meter type (electricity/water)
        meter_type = 'unknown'

        # Strategy 1: Determine from unit
        if unit == 'kWh':
            meter_type = 'electricity'
        elif unit in ['L', 'm³']:
            meter_type = 'water'
        elif unit == 'ZAR':
            meter_type = 'prepaid'

        # Strategy 2: Look for data attributes
        if meter_type == 'unknown':
            meter_type = card.get('data-meter-type') or card.get('data-type') or 'unknown'

        # Strategy 3: Look for specific elements
        if meter_type == 'unknown':
            type_elem = (
                card.find('span', class_='meter-type') or
                card.find('div', class_='meter-type') or
                card.find(class_=re.compile(r'type', re.I))
            )
            if type_elem:
                meter_type = type_elem.text.strip()

        # Strategy 4: Look for keywords in text
        if meter_type == 'unknown':
            card_text_lower = card_text.lower()
            if 'electricity' in card_text_lower or 'electric' in card_text_lower:
                meter_type = 'electricity'
            elif 'water' in card_text_lower:
                meter_type = 'water'
            elif 'gas' in card_text_lower:
                meter_type = 'gas'

        if meter_number:
            meter_data['meter_number'] = str(meter_number).strip()
            meter_data['meter_type'] = str(meter_type).lower()
            meter_data['unit'] = unit if unit else 'kWh'  # Default to kWh
            meter_data['last_updated'] = datetime.now().isoformat()
            logger.debug(f"Parsed meter data: {meter_data}")
        else:
            logger.debug("Could not find meter number in card")

        return meter_data if meter_data.get('meter_number') and meter_data.get('balance') is not None else None

    def update_ha_sensor(self, meter_data):
        """Update Home Assistant sensor using the REST API."""
        try:
            if not SUPERVISOR_TOKEN:
                logger.error("Cannot update sensors - SUPERVISOR_TOKEN not available")
                logger.info("Sensor data that would be created:")
                for meter in meter_data:
                    logger.info(f"  - Meter {meter.get('meter_number')}: {meter.get('balance')} {meter.get('unit')}")
                return

            for meter in meter_data:
                meter_number = meter.get('meter_number', 'unknown')
                balance = meter.get('balance', 0)
                meter_type = meter.get('meter_type', 'unknown')
                unit = meter.get('unit', 'kWh')

                # Create unique entity_id
                entity_id = f"sensor.midcity_{meter_type}_{meter_number.replace(' ', '_').lower()}"

                # Determine device class and icon based on meter type
                if meter_type == 'electricity':
                    device_class = 'energy'
                    icon = 'mdi:lightning-bolt'
                    unit_of_measurement = unit
                elif meter_type == 'water':
                    device_class = 'water'
                    icon = 'mdi:water'
                    unit_of_measurement = unit
                else:
                    device_class = 'monetary'
                    icon = 'mdi:cash'
                    unit_of_measurement = 'ZAR' if unit == 'ZAR' else unit

                # Prepare state data
                state_data = {
                    'state': balance,
                    'attributes': {
                        'meter_number': meter_number,
                        'meter_type': meter_type,
                        'unit_of_measurement': unit_of_measurement,
                        'device_class': device_class,
                        'friendly_name': f'MidCity {meter_type.title()} {meter_number}',
                        'last_updated': meter.get('last_updated'),
                        'icon': icon
                    }
                }

                logger.debug(f"Attempting to create/update sensor: {entity_id}")
                logger.debug(f"API URL: {HA_URL}/states/{entity_id}")
                logger.debug(f"State data: {state_data}")

                # Update state via Home Assistant API
                response = requests.post(
                    f'{HA_URL}/states/{entity_id}',
                    headers=self.headers,
                    json=state_data,
                    timeout=10
                )

                if response.status_code in [200, 201]:
                    logger.info(f"Successfully updated sensor: {entity_id} = {balance} {unit_of_measurement}")
                else:
                    logger.error(f"Failed to update sensor {entity_id}: {response.status_code} - {response.text}")
                    logger.debug(f"Request headers: {self.headers}")
                    logger.debug(f"SUPERVISOR_TOKEN available: {bool(SUPERVISOR_TOKEN)}")

        except Exception as e:
            logger.error(f"Error updating Home Assistant sensor: {e}", exc_info=True)

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
    log_level = config.get('log_level', 'info').upper()

    # Update log level if specified in config
    logger.setLevel(getattr(logging, log_level, logging.INFO))
    logging.getLogger().setLevel(getattr(logging, log_level, logging.INFO))

    if not username or not password:
        logger.error("Username and password are required in configuration")
        sys.exit(1)

    # Create and run sensor
    sensor = MidCityUtilitiesSensor(username, password, scan_interval)
    sensor.run()


if __name__ == '__main__':
    main()
