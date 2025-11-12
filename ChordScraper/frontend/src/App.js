import React from 'react';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';
import SearchPage from './components/SearchPage';

// Material Design 3 (Material You) Dark Theme
const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#D0BCFF',
      light: '#EADDFF',
      dark: '#9A82DB',
      contrastText: '#381E72',
    },
    secondary: {
      main: '#CCC2DC',
      light: '#E8DEF8',
      dark: '#9A8C98',
      contrastText: '#332D41',
    },
    tertiary: {
      main: '#EFB8C8',
      light: '#FFD8E4',
      dark: '#B58392',
      contrastText: '#492532',
    },
    error: {
      main: '#F2B8B5',
      light: '#FFDAD6',
      dark: '#B3261E',
    },
    warning: {
      main: '#F9DEDC',
      light: '#FFEDE9',
      dark: '#8C1D18',
    },
    info: {
      main: '#AAC7E7',
      light: '#D1E4FF',
      dark: '#00639B',
    },
    success: {
      main: '#A6D4A8',
      light: '#D5F0D6',
      dark: '#316B31',
    },
    background: {
      default: '#1C1B1F',
      paper: '#2B2930',
    },
    surface: {
      main: '#2B2930',
      variant: '#49454F',
    },
    text: {
      primary: '#E6E1E5',
      secondary: '#CAC4D0',
      disabled: 'rgba(230, 225, 229, 0.38)',
    },
    divider: 'rgba(202, 196, 208, 0.12)',
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h3: {
      fontWeight: 600,
      letterSpacing: '0px',
    },
    h6: {
      fontWeight: 600,
      letterSpacing: '0.15px',
    },
    button: {
      fontWeight: 500,
      letterSpacing: '0.1px',
    },
  },
  shape: {
    borderRadius: 16,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          fontWeight: 500,
          borderRadius: 20,
          paddingLeft: 24,
          paddingRight: 24,
          paddingTop: 10,
          paddingBottom: 10,
        },
        contained: {
          boxShadow: 'none',
          '&:hover': {
            boxShadow: '0px 1px 2px rgba(0, 0, 0, 0.3), 0px 1px 3px 1px rgba(0, 0, 0, 0.15)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
          borderRadius: 20,
          boxShadow: 'none',
          border: '1px solid rgba(202, 196, 208, 0.12)',
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
          borderRadius: 20,
        },
        elevation3: {
          boxShadow: '0px 1px 3px rgba(0, 0, 0, 0.3), 0px 4px 8px 3px rgba(0, 0, 0, 0.15)',
        },
      },
    },
    MuiTab: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          fontWeight: 500,
          fontSize: '0.875rem',
          minHeight: 48,
          letterSpacing: '0.1px',
        },
      },
    },
    MuiTabs: {
      styleOverrides: {
        indicator: {
          height: 3,
          borderRadius: '3px 3px 0 0',
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            borderRadius: 12,
          },
        },
      },
    },
    MuiDialog: {
      styleOverrides: {
        paper: {
          borderRadius: 28,
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          fontWeight: 500,
        },
      },
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <SearchPage />
    </ThemeProvider>
  );
}

export default App;
