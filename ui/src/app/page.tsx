"use client";
import { Box, Button, Typography, Paper } from "@mui/material";

export default function LoginPage() {
  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "100vh",
      }}
    >
      <Paper sx={{ p: 6, textAlign: "center", maxWidth: 400 }}>
        <Typography variant="h3" gutterBottom sx={{ fontWeight: "bold" }}>
          LootForge
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
          Fair loot distribution for gaming groups
        </Typography>
        <Button variant="contained" size="large" fullWidth disabled>
          Sign in with Google (coming soon)
        </Button>
      </Paper>
    </Box>
  );
}
