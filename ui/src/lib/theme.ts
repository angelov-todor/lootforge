"use client";
import { createTheme } from "@mui/material/styles";

export function buildTheme(mode: "light" | "dark") {
  return createTheme({
    palette: {
      mode,
      primary: { main: "#6366f1" },
      secondary: { main: "#f59e0b" },
      ...(mode === "dark"
        ? {
            background: { default: "#0f1117", paper: "#1a1d2e" },
          }
        : {
            background: { default: "#f8fafc", paper: "#ffffff" },
          }),
    },
    typography: {
      fontFamily: "inherit",
    },
    components: {
      MuiButton: {
        styleOverrides: {
          root: { textTransform: "none", borderRadius: 8 },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: { borderRadius: 12 },
        },
      },
    },
  });
}

// Default export for backward compatibility
const theme = buildTheme("dark");
export default theme;
