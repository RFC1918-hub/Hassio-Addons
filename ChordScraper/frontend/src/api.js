// API client for Chord Scraper backend
const API_BASE_URL = process.env.REACT_APP_API_URL || '';

/**
 * Search for guitar tabs on Ultimate Guitar
 * @param {string} query - Search query (song title)
 * @returns {Promise<Array>} Array of search results
 */
export async function searchTabs(query) {
  const response = await fetch(`${API_BASE_URL}/search?title=${encodeURIComponent(query)}`);

  if (!response.ok) {
    throw new Error(`Search failed: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Get tab in OnSong format by ID
 * @param {number} tabId - Tab ID from Ultimate Guitar
 * @returns {Promise<string>} Tab content in OnSong format
 */
export async function getOnSongFormat(tabId) {
  const response = await fetch(`${API_BASE_URL}/onsong`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ id: tabId }),
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch tab: ${response.statusText}`);
  }

  return response.text();
}

/**
 * Get worship chords from WorshipChords.com
 * @param {string} url - WorshipChords.com URL
 * @returns {Promise<string>} Chord content
 */
export async function getWorshipChords(url) {
  const response = await fetch(`${API_BASE_URL}/worshipchords`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ url }),
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch worship chords: ${response.statusText}`);
  }

  return response.text();
}

/**
 * Send chord content to Google Drive via n8n webhook
 * @param {Object} data - Data to send
 * @param {string} data.content - Chord content
 * @param {string} data.song - Song name
 * @param {string} data.artist - Artist name
 * @param {number} data.id - Tab ID
 * @param {boolean} data.isManualSubmission - Manual submission flag
 * @param {boolean} data.requiresAutomation - Automation required flag
 * @returns {Promise<Object>} Response from webhook
 */
export async function sendToGoogleDrive(data) {
  const response = await fetch(`${API_BASE_URL}/send-to-drive`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `Failed to send to Google Drive: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Health check
 * @returns {Promise<boolean>} True if server is healthy
 */
export async function healthCheck() {
  try {
    const response = await fetch(`${API_BASE_URL}/health`);
    return response.ok;
  } catch (error) {
    return false;
  }
}
