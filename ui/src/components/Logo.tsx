"use client";
import { Box, Typography } from "@mui/material";
import { DieIcon } from "@/components/DieIcon";

interface LogoProps {
  size?: "small" | "large";
}

export function Logo({ size = "small" }: LogoProps) {
  if (size === "large") {
    return (
      <Box sx={{ textAlign: "center" }}>
        <Box
          sx={{
            width: 72,
            height: 72,
            borderRadius: 3,
            background: "linear-gradient(135deg, #6366f1, #818cf8)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            mx: "auto",
            mb: 2,
            boxShadow: "0 4px 20px rgba(99,102,241,0.3)",
          }}
        >
          <DieIcon size={40} color="white" />
        </Box>
        <Typography variant="h4" sx={{ fontWeight: 800 }}>
          <Box
            component="span"
            sx={{
              background: "linear-gradient(135deg, #818cf8, #f59e0b)",
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
            }}
          >
            LOOT
          </Box>
          <Box component="span" sx={{ color: "text.primary" }}>
            FORGE
          </Box>
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 0.75 }}>
      <Box
        sx={{
          width: 28,
          height: 28,
          borderRadius: 1,
          background: "linear-gradient(135deg, #6366f1, #818cf8)",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <DieIcon size={18} color="white" />
      </Box>
      <Typography variant="subtitle1" sx={{ fontWeight: 800, lineHeight: 1 }}>
        <Box component="span" sx={{ color: "primary.main" }}>
          LOOT
        </Box>
        <Box component="span" sx={{ color: "text.primary" }}>
          FORGE
        </Box>
      </Typography>
    </Box>
  );
}
