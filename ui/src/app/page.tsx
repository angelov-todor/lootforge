"use client";
export const dynamic = "force-dynamic";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Box, Button, Typography, Paper, CircularProgress } from "@mui/material";
import GoogleIcon from "@mui/icons-material/Google";
import { useAuth } from "@/components/AuthProvider";

export default function LoginPage() {
  const { user, loading, signIn } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && user) {
      router.replace("/dashboard");
    }
  }, [user, loading, router]);

  if (loading) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          minHeight: "100vh",
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "100vh",
        background: "linear-gradient(135deg, #0f172a 0%, #1e1b4b 100%)",
      }}
    >
      <Paper sx={{ p: 6, textAlign: "center", maxWidth: 400, width: "100%" }}>
        <Typography variant="h3" gutterBottom sx={{ fontWeight: "bold", color: "primary.main" }}>
          LootForge
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
          Fair loot distribution for gaming groups
        </Typography>
        <Button
          variant="contained"
          size="large"
          fullWidth
          startIcon={<GoogleIcon />}
          onClick={signIn}
          sx={{ py: 1.5 }}
        >
          Sign in with Google
        </Button>
      </Paper>
    </Box>
  );
}
