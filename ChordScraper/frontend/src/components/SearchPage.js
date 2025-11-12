import React, { useState, useCallback } from 'react';
import {
  Box,
  Container,
  TextField,
  Button,
  Typography,
  Card,
  CardContent,
  CardActions,
  Chip,
  CircularProgress,
  Tabs,
  Tab,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Snackbar,
  Alert,
  Stack,
  Paper,
  Divider,
  Textarea,
} from '@mui/material';
import { Search as SearchIcon, CloudUpload, MusicNote, Star } from '@mui/icons-material';
import axios from 'axios';
import debounce from 'lodash/debounce';
import { getOnSongFormat } from '../api';

const API_URL = window.location.origin;

function SearchPage() {
  const [tabIndex, setTabIndex] = useState(0);
  const [searchTerm, setSearchTerm] = useState('');
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedTab, setSelectedTab] = useState(null);
  const [onsongContent, setOnsongContent] = useState('');
  const [modalOpen, setModalOpen] = useState(false);
  const [sendingToDrive, setSendingToDrive] = useState(false);
  const [manualSong, setManualSong] = useState('');
  const [manualArtist, setManualArtist] = useState('');
  const [manualContent, setManualContent] = useState('');
  const [submittingManual, setSubmittingManual] = useState(false);
  const [ugTabId, setUgTabId] = useState('');
  const [loadingUgTab, setLoadingUgTab] = useState(false);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'info' });

  const showSnackbar = (message, severity = 'info') => {
    setSnackbar({ open: true, message, severity });
  };

  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
  };

  // Debounced search function
  const debouncedSearch = useCallback(
    debounce(async (query) => {
      if (!query || query.length < 3) {
        setResults([]);
        return;
      }

      try {
        setLoading(true);
        const response = await axios.get(`${API_URL}/search?title=${encodeURIComponent(query)}`);
        // Sort by rating (highest first)
        const sortedResults = response.data.sort((a, b) => (b.rating || 0) - (a.rating || 0));
        setResults(sortedResults);
      } catch (error) {
        showSnackbar(error.response?.data?.error || 'Failed to search', 'error');
      } finally {
        setLoading(false);
      }
    }, 500),
    []
  );

  const handleSearchChange = (e) => {
    const value = e.target.value;
    setSearchTerm(value);
    debouncedSearch(value);
  };

  const handleGetOnsong = async (tab) => {
    try {
      setSelectedTab(tab);
      const response = await axios.post(`${API_URL}/onsong`, { id: tab.id });
      setOnsongContent(response.data);
      setModalOpen(true);
    } catch (error) {
      showSnackbar(
        typeof error.response?.data === 'string'
          ? error.response.data
          : error.response?.data?.message || 'Failed to get OnSong format',
        'error'
      );
    }
  };

  const handleSendToDrive = async () => {
    try {
      setSendingToDrive(true);
      await axios.post(`${API_URL}/send-to-drive`, {
        content: onsongContent,
        song: selectedTab?.song,
        artist: selectedTab?.artist,
        id: String(selectedTab?.id), // Convert to string to support all ID types
        isManualSubmission: selectedTab?.isManual || false,
      });
      showSnackbar('Successfully sent to Google Drive!', 'success');
      setModalOpen(false);
    } catch (error) {
      showSnackbar(
        error.response?.data?.error || error.response?.data?.message || 'Failed to send to Google Drive',
        'error'
      );
    } finally {
      setSendingToDrive(false);
    }
  };

  const handleGetUgTabById = async () => {
    if (!ugTabId) {
      showSnackbar('Please enter an Ultimate Guitar tab ID', 'error');
      return;
    }

    // Convert to integer and validate
    const tabIdNumber = parseInt(ugTabId, 10);
    if (isNaN(tabIdNumber) || tabIdNumber <= 0) {
      showSnackbar('Please enter a valid numeric tab ID', 'error');
      return;
    }

    try {
      setLoadingUgTab(true);
      const onsongData = await getOnSongFormat(tabIdNumber);

      // Try to extract song/artist info from the formatted content
      const lines = onsongData.split('\n');
      let song = 'Unknown Song';
      let artist = 'Unknown Artist';

      if (lines.length > 0 && lines[0].trim()) {
        song = lines[0].trim();
      }

      if (lines.length > 1 && lines[1].trim()) {
        artist = lines[1].trim();
      }

      setSelectedTab({ song, artist, id: ugTabId });
      setOnsongContent(onsongData);
      setModalOpen(true);
      setUgTabId(''); // Clear the input after successful fetch
      showSnackbar('Tab loaded successfully', 'success');
    } catch (error) {
      showSnackbar(
        error.message || 'Failed to get tab by ID',
        'error'
      );
    } finally {
      setLoadingUgTab(false);
    }
  };

  const handleManualSubmission = async () => {
    if (!manualSong.trim()) {
      showSnackbar('Please enter a song title', 'error');
      return;
    }

    if (!manualContent.trim()) {
      showSnackbar('Please enter chord content', 'error');
      return;
    }

    try {
      setSubmittingManual(true);
      const song = manualSong.trim();
      const artist = manualArtist.trim() || 'Unknown Artist';
      const id = 'manual-' + Date.now();

      // Call the format-manual endpoint to get formatted content
      const response = await axios.post(`${API_URL}/format-manual`, {
        song: song,
        artist: artist,
        content: manualContent,
      });

      // Set up the preview modal with formatted content
      setSelectedTab({ song, artist, id, isManual: true });
      setOnsongContent(response.data);
      setModalOpen(true);

      // Clear the form after showing preview
      setManualSong('');
      setManualArtist('');
      setManualContent('');
    } catch (error) {
      showSnackbar(
        error.response?.data?.error || error.response?.data?.message || 'Failed to format content',
        'error'
      );
    } finally {
      setSubmittingManual(false);
    }
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Box sx={{ mb: 4 }}>
        <Typography
          variant="h3"
          component="h1"
          gutterBottom
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 2,
            fontWeight: 700,
          }}
        >
          <MusicNote sx={{ fontSize: 48 }} />
          Chord Scraper
        </Typography>
        <Typography variant="subtitle1" color="text.secondary">
          Search Ultimate Guitar tabs • Convert to OnSong format
        </Typography>
      </Box>

      <Paper elevation={3} sx={{ mb: 4 }}>
        <Tabs
          value={tabIndex}
          onChange={(e, newValue) => setTabIndex(newValue)}
          variant="fullWidth"
          sx={{ borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab label="Search Tabs" />
          <Tab label="Tab by ID" />
          <Tab label="Manual Submission" />
        </Tabs>

        {/* Ultimate Guitar Tab */}
        {tabIndex === 0 && (
          <Box sx={{ p: 3 }}>
            <TextField
              fullWidth
              variant="outlined"
              placeholder="Search for a song... (type at least 3 characters)"
              value={searchTerm}
              onChange={handleSearchChange}
              InputProps={{
                startAdornment: <SearchIcon sx={{ mr: 1, color: 'text.secondary' }} />,
              }}
              sx={{ mb: 1 }}
            />

            {results.length > 0 && !loading && (
              <Typography variant="caption" color="text.secondary" sx={{ display: 'flex', alignItems: 'center', gap: 0.5, mb: 2 }}>
                <Star sx={{ fontSize: 16, color: 'warning.main' }} />
                Sorted by highest rating first
              </Typography>
            )}

            {loading && (
              <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
                <CircularProgress />
              </Box>
            )}

            <Stack spacing={2}>
              {results.map((result) => (
                <Card key={result.id} variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {result.song || 'Unknown Song'}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      Artist: {result.artist || 'Unknown Artist'}
                    </Typography>
                    <Box sx={{ display: 'flex', gap: 1, alignItems: 'center', mt: 1 }}>
                      <Chip label={result.type || 'Unknown'} color="primary" size="small" />
                      <Chip
                        icon={<Star sx={{ fontSize: 18 }} />}
                        label={result.rating ? result.rating.toFixed(2) : 'N/A'}
                        color="warning"
                        size="small"
                      />
                    </Box>
                  </CardContent>
                  <CardActions>
                    <Button
                      variant="contained"
                      onClick={() => handleGetOnsong(result)}
                      fullWidth
                    >
                      Get OnSong Format
                    </Button>
                  </CardActions>
                </Card>
              ))}
            </Stack>
          </Box>
        )}

        {/* Ultimate Guitar Tab by ID */}
        {tabIndex === 1 && (
          <Box sx={{ p: 3 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Enter an Ultimate Guitar tab ID to fetch directly. You can find the ID in the Ultimate Guitar URL.
            </Typography>
            <Typography variant="caption" color="text.secondary" sx={{ mb: 2, display: 'block' }}>
              Example: For URL "https://tabs.ultimate-guitar.com/tab/1234567", the ID is "1234567"
            </Typography>
            <Stack direction="row" spacing={2}>
              <TextField
                fullWidth
                variant="outlined"
                placeholder="Enter Ultimate Guitar tab ID..."
                value={ugTabId}
                onChange={(e) => setUgTabId(e.target.value)}
              />
              <Button
                variant="contained"
                color="primary"
                onClick={handleGetUgTabById}
                disabled={loadingUgTab}
                sx={{ minWidth: 150 }}
              >
                {loadingUgTab ? <CircularProgress size={24} /> : 'Get Tab'}
              </Button>
            </Stack>
          </Box>
        )}

        {/* Manual Submission Tab */}
        {tabIndex === 2 && (
          <Box sx={{ p: 3 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Submit raw chord text directly. Artist is optional (will default to "Unknown").
            </Typography>
            <Stack spacing={2}>
              <TextField
                fullWidth
                variant="outlined"
                placeholder="Song Title"
                value={manualSong}
                onChange={(e) => setManualSong(e.target.value)}
              />
              <TextField
                fullWidth
                variant="outlined"
                placeholder="Artist Name (optional)"
                value={manualArtist}
                onChange={(e) => setManualArtist(e.target.value)}
              />
              <TextField
                fullWidth
                multiline
                rows={12}
                variant="outlined"
                placeholder="Paste or type the raw chord content here..."
                value={manualContent}
                onChange={(e) => setManualContent(e.target.value)}
                sx={{ fontFamily: 'monospace' }}
              />
              <Button
                variant="contained"
                color="primary"
                onClick={handleManualSubmission}
                disabled={submittingManual}
                startIcon={submittingManual ? <CircularProgress size={20} /> : null}
              >
                {submittingManual ? 'Formatting...' : 'Preview Format'}
              </Button>
            </Stack>
          </Box>
        )}
      </Paper>

      {/* OnSong Content Modal */}
      <Dialog
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          {selectedTab?.song || 'Unknown Song'} - {selectedTab?.artist || 'Unknown Artist'}
        </DialogTitle>
        <DialogContent dividers>
          <TextField
            fullWidth
            multiline
            rows={20}
            variant="outlined"
            value={typeof onsongContent === 'string'
              ? onsongContent
              : JSON.stringify(onsongContent, null, 2)
            }
            onChange={(e) => setOnsongContent(e.target.value)}
            sx={{
              fontFamily: 'monospace',
              '& .MuiInputBase-input': {
                fontFamily: 'monospace',
                fontSize: '0.875rem',
              },
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setModalOpen(false)}>Close</Button>
          <Button
            variant="contained"
            color="success"
            onClick={handleSendToDrive}
            disabled={sendingToDrive}
            startIcon={sendingToDrive ? <CircularProgress size={20} /> : <CloudUpload />}
          >
            {sendingToDrive ? 'Sending...' : 'Send to Google Drive'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Snackbar for notifications */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={5000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={handleCloseSnackbar} severity={snackbar.severity} variant="filled">
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Container>
  );
}

export default SearchPage;
