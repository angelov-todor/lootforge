"use client";
import { createTheme } from "@mui/material/styles";

const theme = createTheme({
  palette: {
    mode: "dark",
    primary: { main: "#7c3aed" },
    secondary: { main: "#f59e0b" },
    background: { default: "#0f172a", paper: "#1e293b" },
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
  },
});

export default theme;
