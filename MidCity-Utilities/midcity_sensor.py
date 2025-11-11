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
import paho.mqtt.client as mqtt

# Set up logging
log_level = os.environ.get('LOG_LEVEL', 'INFO').upper()
logging.basicConfig(
    level=getattr(logging, log_level, logging.INFO),
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Get MQTT configuration from Supervisor services
def get_mqtt_config():
    """Get MQTT configuration from Supervisor services API."""
    try:
        supervisor_token = os.environ.get('SUPERVISOR_TOKEN')
        if not supervisor_token:
            logger.warning("No SUPERVISOR_TOKEN found, using default MQTT config")
            return {
                'host': 'localhost',
                'port': 1883,
                'user': None,
                'password': None
            }

        headers = {
            'Authorization': f'Bearer {supervisor_token}',
            'Content-Type': 'application/json'
        }

        response = requests.get('http://supervisor/services/mqtt', headers=headers, timeout=10)
        if response.status_code == 200:
            mqtt_config = response.json()['data']
            logger.info(f"Retrieved MQTT config from Supervisor: {mqtt_config.get('host')}:{mqtt_config.get('port')}")
            return mqtt_config
        else:
            logger.warning(f"Could not get MQTT config from Supervisor: {response.status_code}")
            return {
                'host': 'localhost',
                'port': 1883,
                'user': None,
                'password': None
            }
    except Exception as e:
        logger.warning(f"Error getting MQTT config: {e}")
        return {
            'host': 'localhost',
            'port': 1883,
            'user': None,
            'password': None
        }

# Get MQTT configuration
MQTT_CONFIG = get_mqtt_config()
MQTT_HOST = MQTT_CONFIG.get('host', 'localhost')
MQTT_PORT = MQTT_CONFIG.get('port', 1883)
MQTT_USER = MQTT_CONFIG.get('username') or MQTT_CONFIG.get('user')
MQTT_PASSWORD = MQTT_CONFIG.get('password')

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

        # Initialize MQTT client
        self.mqtt_client = mqtt.Client()

        # Set MQTT credentials if available
        if MQTT_USER and MQTT_PASSWORD:
            logger.info(f"Using MQTT authentication with user: {MQTT_USER}")
            self.mqtt_client.username_pw_set(MQTT_USER, MQTT_PASSWORD)
        else:
            logger.info("Using MQTT without authentication (anonymous)")

        # Set MQTT callbacks
        self.mqtt_client.on_connect = self.on_mqtt_connect
        self.mqtt_client.on_disconnect = self.on_mqtt_disconnect

        self.mqtt_connected = False

        # Connect to MQTT broker
        try:
            logger.info(f"Connecting to MQTT broker at {MQTT_HOST}:{MQTT_PORT}")
            self.mqtt_client.connect(MQTT_HOST, MQTT_PORT, 60)
            self.mqtt_client.loop_start()
        except Exception as e:
            logger.error(f"Failed to connect to MQTT broker: {e}")
            logger.error(f"Make sure Mosquitto broker add-on is installed and running")

    def on_mqtt_connect(self, client, userdata, flags, rc):
        """MQTT connection callback."""
        if rc == 0:
            logger.info("Successfully connected to MQTT broker")
            self.mqtt_connected = True
        else:
            error_messages = {
                1: "Connection refused - incorrect protocol version",
                2: "Connection refused - invalid client identifier",
                3: "Connection refused - server unavailable",
                4: "Connection refused - bad username or password",
                5: "Connection refused - not authorized"
            }
            error_msg = error_messages.get(rc, f"Unknown error code: {rc}")
            logger.error(f"Failed to connect to MQTT broker: {error_msg} (code: {rc})")
            if rc == 5:
                logger.error("MQTT authentication failed. Check Mosquitto broker configuration.")
                logger.error("You may need to configure Mosquitto to allow anonymous connections or add user credentials.")
            self.mqtt_connected = False

    def on_mqtt_disconnect(self, client, userdata, rc):
        """MQTT disconnection callback."""
        logger.warning(f"Disconnected from MQTT broker with code: {rc}")
        self.mqtt_connected = False

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

            # Log some info about the page structure
            logger.info(f"Page title: {soup.title.string if soup.title else 'No title'}")

            # Use a different strategy: find meter number and balance separately then combine
            import re

            # Find meter number
            meter_number = None

            # Strategy 1: Look in "Select Meter" dropdown/text
            meter_select = soup.find(string=re.compile(r'Select Meter'))
            if meter_select:
                match = re.search(r'(\d{8,12})', meter_select)
                if match:
                    meter_number = match.group(1)
                    logger.info(f"Found meter number in 'Select Meter': {meter_number}")

            # Strategy 2: Look in select/option elements
            if not meter_number:
                select_elem = soup.find('select', {'name': 'meter_id'})
                if select_elem:
                    options = select_elem.find_all('option')
                    for option in options:
                        if option.get('value') and re.match(r'\d{8,12}', option.get('value')):
                            meter_number = option.get('value')
                            logger.info(f"Found meter number in select option: {meter_number}")
                            break

            # Strategy 3: Look anywhere in the page for meter number patterns
            if not meter_number:
                page_text = soup.get_text()
                match = re.search(r'(?:Meter|Account)[^\d]*(\d{10,12})', page_text, re.I)
                if match:
                    meter_number = match.group(1)
                    logger.info(f"Found meter number in page text: {meter_number}")

            # Find current balance - it's embedded in JavaScript JSON
            balance = None
            unit = None

            # Strategy 1: Look for balance in JavaScript chartObjects
            scripts = soup.find_all('script')
            for script in scripts:
                if script.string and 'chartObjects' in script.string:
                    script_text = script.string
                    # Look for the series data with balance
                    # Pattern: "series":[{"name":"Current balance","0":[145.65],"tooltip":{"valueSuffix":" kWh"}}]
                    match = re.search(r'"name":"Current balance"[^}]*?"0":\[([0-9.]+)\]', script_text)
                    if match:
                        balance = float(match.group(1))
                        unit = 'kWh'
                        logger.info(f"Found balance in JavaScript: {balance} {unit}")
                        break

                    # Alternative pattern
                    match = re.search(r'"series":\[{[^}]*?"0":\[([0-9.]+)\][^}]*?"valueSuffix":" kWh"', script_text)
                    if match:
                        balance = float(match.group(1))
                        unit = 'kWh'
                        logger.info(f"Found balance in JavaScript (alt): {balance} {unit}")
                        break

            # Strategy 2: Look for balance in HTML text (fallback)
            if not balance:
                balance_text_elem = soup.find(string=re.compile(r'Current meter balance[:\s]+[\d.]+\s*kWh', re.I))
                if balance_text_elem:
                    balance_text = balance_text_elem.strip()
                    logger.debug(f"Found balance text in HTML: {balance_text}")
                    match = re.search(r'([\d,]+\.?\d*)\s*kWh', balance_text, re.I)
                    if match:
                        balance_str = match.group(1).replace(',', '')
                        balance = float(balance_str)
                        unit = 'kWh'
                        logger.info(f"Found balance in HTML: {balance} {unit}")

            # Find predicted zero balance date
            predicted_zero_date = None

            # Strategy 1: Search in all text for the predicted date
            page_text = soup.get_text()
            # Look for "Predicted 0 balance date: 2025-12-13" pattern
            date_match = re.search(r'Predicted\s+0\s+balance\s+date[:\s]+(\d{4}-\d{2}-\d{2})', page_text, re.I)
            if date_match:
                predicted_zero_date = date_match.group(1)
                logger.info(f"Found predicted zero date: {predicted_zero_date}")
            else:
                # Try alternative format with slashes
                date_match = re.search(r'Predicted\s+0\s+balance\s+date[:\s]+(\d{1,2}[/-]\d{1,2}[/-]\d{4})', page_text, re.I)
                if date_match:
                    predicted_zero_date = date_match.group(1)
                    logger.info(f"Found predicted zero date: {predicted_zero_date}")
                else:
                    logger.debug("Predicted zero date not found in page text")

            # Create meter data if we found both
            meters = []
            if meter_number and balance is not None:
                meter_data = {
                    'meter_number': meter_number,
                    'balance': balance,
                    'meter_type': 'electricity',  # kWh indicates electricity
                    'unit': unit,
                    'last_updated': datetime.now().isoformat()
                }

                # Add predicted zero date if found
                if predicted_zero_date:
                    meter_data['predicted_zero_date'] = predicted_zero_date

                meters.append(meter_data)
                logger.info(f"Successfully combined meter data: {meter_data}")
            else:
                logger.warning(f"Could not combine meter data - meter_number: {meter_number}, balance: {balance}")

            logger.info(f"Retrieved data for {len(meters)} meter(s)")
            return meters if meters else None

        except Exception as e:
            logger.error(f"Error retrieving meter data: {e}", exc_info=True)
            return None

    def get_meter_data_old(self):
        """Old method - kept for reference."""
        try:
            response = self.session.get(METER_URL, timeout=30)

            if response.status_code != 200:
                logger.error(f"Failed to retrieve meter data. Status code: {response.status_code}")
                return None

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

    def publish_mqtt_discovery(self, meter_data):
        """Publish MQTT Discovery messages for Home Assistant."""
        try:
            if not self.mqtt_connected:
                logger.error("Cannot publish sensors - MQTT not connected")
                logger.info("Sensor data that would be created:")
                for meter in meter_data:
                    logger.info(f"  - Meter {meter.get('meter_number')}: {meter.get('balance')} {meter.get('unit')}")
                return

            for meter in meter_data:
                meter_number = meter.get('meter_number', 'unknown')
                balance = meter.get('balance', 0)
                meter_type = meter.get('meter_type', 'unknown')
                unit = meter.get('unit', 'kWh')
                predicted_zero_date = meter.get('predicted_zero_date', None)

                # Create unique_id for this sensor
                unique_id = f"midcity_{meter_type}_{meter_number}"
                object_id = f"midcity_{meter_type}_{meter_number}"

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

                # Friendly name for the entity
                friendly_name = f'MidCity {meter_type.title()}'

                # MQTT topics
                config_topic = f"homeassistant/sensor/midcity_utilities/{object_id}/config"
                state_topic = f"homeassistant/sensor/midcity_utilities/{object_id}/state"
                attributes_topic = f"homeassistant/sensor/midcity_utilities/{object_id}/attributes"

                # Discovery payload
                discovery_payload = {
                    "name": friendly_name,
                    "unique_id": unique_id,
                    "object_id": object_id,
                    "state_topic": state_topic,
                    "json_attributes_topic": attributes_topic,
                    "unit_of_measurement": unit_of_measurement,
                    "device_class": device_class,
                    "icon": icon,
                    "state_class": "measurement",
                    "device": {
                        "identifiers": ["midcity_utilities_sensor"],
                        "name": "MidCity Utilities Sensor",
                        "model": "MidCity Utilities Monitor",
                        "manufacturer": "MidCity Utilities",
                        "sw_version": "1.2.1"
                    }
                }

                # Publish discovery message
                logger.info(f"Publishing MQTT discovery for {unique_id}")
                self.mqtt_client.publish(
                    config_topic,
                    json.dumps(discovery_payload),
                    retain=True
                )

                # Publish state
                logger.debug(f"Publishing state: {balance}")
                self.mqtt_client.publish(
                    state_topic,
                    str(balance),
                    retain=True
                )

                # Publish attributes
                attributes = {
                    'meter_number': meter_number,
                    'meter_type': meter_type,
                    'last_updated': meter.get('last_updated'),
                    'attribution': 'Data from MidCity Utilities'
                }

                # Add predicted zero date if available
                if predicted_zero_date:
                    attributes['predicted_zero_date'] = predicted_zero_date
                    logger.debug(f"Adding predicted zero date to attributes: {predicted_zero_date}")

                logger.debug(f"Publishing attributes: {attributes}")
                self.mqtt_client.publish(
                    attributes_topic,
                    json.dumps(attributes),
                    retain=True
                )

                logger.info(f"Successfully published sensor: sensor.{object_id} = {balance} {unit_of_measurement}")

        except Exception as e:
            logger.error(f"Error publishing MQTT discovery: {e}", exc_info=True)

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
                    # Publish MQTT Discovery messages
                    self.publish_mqtt_discovery(meter_data)
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
    sensor = MidCityUtilitiesSensor(
        username,
        password,
        scan_interval
    )
    sensor.run()


if __name__ == '__main__':
    main()
