"use client";
export const dynamic = "force-dynamic";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Box, Button, CircularProgress } from "@mui/material";
import GoogleIcon from "@mui/icons-material/Google";
import { useAuth } from "@/components/AuthProvider";
import { Logo } from "@/components/Logo";

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
      <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", minHeight: "100vh" }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "100vh",
        background: "linear-gradient(135deg, #0f172a 0%, #1e1b4b 100%)",
        gap: 4,
        px: 3,
      }}
    >
      <Logo size="large" />
      <Button
        variant="contained"
        size="large"
        startIcon={<GoogleIcon />}
        onClick={signIn}
        sx={{ py: 1.5, px: 4, maxWidth: 320, width: "100%" }}
      >
        Sign in with Google
      </Button>
    </Box>
  );
}
