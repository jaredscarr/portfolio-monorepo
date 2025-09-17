"use client";

import { createTheme, ThemeOptions } from "@mui/material/styles";

// Define your color palette - customize these colors using the tools mentioned above
const colors = {
  primary: {
    light: "#bdc7bf",
    dark: "#546357",
  },
  secondary: {
    light: "#ae675b",
    dark: "#bc8176",
  },
  background: {
    light: "#ffffff",
    dark: "#121212",
  },
  surface: {
    light: "#f5f5f5",
    dark: "#1e1e1e",
  },
  text: {
    primary: {
      light: "#000000",
      dark: "#ffffff",
    },
    secondary: {
      light: "#666666",
      dark: "#b0b0b0",
    },
  },
};

// Base theme options shared across all color schemes
const baseTheme: ThemeOptions = {
  typography: {
    fontFamily: "var(--font-geist-sans), sans-serif",
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    // Add consistent component overrides here
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: "none", // Prevent uppercase transformation
        },
      },
    },
  },
};

// Create theme with colorSchemes API for manual theme toggling
export const theme = createTheme({
  ...baseTheme,
  colorSchemes: {
    light: {
      palette: {
        primary: {
          main: colors.primary.light,
        },
        secondary: {
          main: colors.secondary.light,
        },
        background: {
          default: colors.background.light,
          paper: colors.surface.light,
        },
        text: {
          primary: colors.text.primary.light,
          secondary: colors.text.secondary.light,
        },
      },
    },
    dark: {
      palette: {
        primary: {
          main: colors.primary.dark,
        },
        secondary: {
          main: colors.secondary.dark,
        },
        background: {
          default: colors.background.dark,
          paper: colors.surface.dark,
        },
        text: {
          primary: colors.text.primary.dark,
          secondary: colors.text.secondary.dark,
        },
      },
    },
  },
  // Set default color scheme
  defaultColorScheme: "light",
  // Enable CSS variables and configure selector for manual toggling
  cssVariables: {
    colorSchemeSelector: "data-mui-color-scheme",
  },
});
