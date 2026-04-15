"use client";
import { Box, Typography } from "@mui/material";

interface LogoProps {
  size?: "small" | "large";
}

export function Logo({ size = "small" }: LogoProps) {
  if (size === "large") {
    return (
      <Box sx={{ textAlign: "center" }}>
        <Box sx={{ fontSize: 48, position: "relative", display: "inline-block" }}>
          <Box
            component="span"
            sx={{ filter: "drop-shadow(0 0 8px rgba(124,58,237,0.5))" }}
          >
            &#9879;
          </Box>
          <Box
            component="span"
            sx={{
              position: "absolute",
              top: -8,
              right: -12,
              fontSize: 24,
              filter: "drop-shadow(0 0 6px rgba(245,158,11,0.5))",
            }}
          >
            &#127922;
          </Box>
        </Box>
        <Typography
          variant="h4"
          sx={{ fontWeight: 800, mt: 1 }}
        >
          <Box
            component="span"
            sx={{
              background: "linear-gradient(135deg, #a855f7, #f59e0b)",
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
            }}
          >
            LOOT
          </Box>
          <Box component="span" sx={{ color: "white" }}>
            FORGE
          </Box>
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
      <Box sx={{ fontSize: 20, position: "relative", display: "inline-block", lineHeight: 1 }}>
        <span>&#9879;</span>
        <Box
          component="span"
          sx={{ position: "absolute", top: -4, right: -6, fontSize: 10 }}
        >
          &#127922;
        </Box>
      </Box>
      <Typography variant="subtitle1" sx={{ fontWeight: 800, lineHeight: 1 }}>
        <Box
          component="span"
          sx={{ color: "primary.main" }}
        >
          LOOT
        </Box>
        <Box component="span" sx={{ color: "white" }}>
          FORGE
        </Box>
      </Typography>
    </Box>
  );
}
