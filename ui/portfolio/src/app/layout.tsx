import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { AppRouterCacheProvider } from "@mui/material-nextjs/v15-appRouter";
import { ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { Box, Typography, IconButton } from "@mui/material";
import { GitHub, LinkedIn } from "@mui/icons-material";
import { theme } from "./theme";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
  display: "swap",
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
  display: "swap",
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <meta name="color-scheme" content="light dark" />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
        style={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}
      >
        <AppRouterCacheProvider>
          <ThemeProvider theme={theme}>
            <CssBaseline />
            <Box sx={{ flex: 1 }}>{children}</Box>
            <Box
              component="footer"
              sx={{
                py: 2,
                px: 2,
                textAlign: "center",
                borderTop: 1,
                borderColor: "divider",
              }}
            >
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "center",
                  gap: 1,
                  mb: 1,
                }}
              >
                <IconButton
                  component="a"
                  href="https://github.com/jaredscarr"
                  target="_blank"
                  rel="noopener noreferrer"
                  color="inherit"
                  aria-label="GitHub profile"
                  size="small"
                >
                  <GitHub />
                </IconButton>
                <IconButton
                  component="a"
                  href="https://linkedin.com/in/jaredscarr"
                  target="_blank"
                  rel="noopener noreferrer"
                  color="inherit"
                  aria-label="LinkedIn profile"
                  size="small"
                >
                  <LinkedIn />
                </IconButton>
              </Box>
              <Typography variant="body2" color="text.secondary">
                © {new Date().getFullYear()} Jared Scarr. All rights reserved.
              </Typography>
            </Box>
          </ThemeProvider>
        </AppRouterCacheProvider>
      </body>
    </html>
  );
}
