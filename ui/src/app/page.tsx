"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Box, Button, CircularProgress } from "@mui/material";
import GoogleIcon from "@mui/icons-material/Google";
import { Logo } from "@/components/Logo";

export default function LoginPage() {
  const router = useRouter();
  const [checking, setChecking] = useState(true);
  const [signingIn, setSigningIn] = useState(false);

  useEffect(() => {
    // Lazy check auth status — don't block first paint
    import("@/lib/firebase").then(({ onAuthStateChanged }) => {
      const unsub = onAuthStateChanged((user) => {
        if (user) {
          router.replace("/dashboard");
        } else {
          setChecking(false);
        }
        unsub();
      });
    });
  }, [router]);

  const handleSignIn = async () => {
    setSigningIn(true);
    try {
      const { signInWithGoogle } = await import("@/lib/firebase");
      await signInWithGoogle();
      router.replace("/dashboard");
    } catch {
      setSigningIn(false);
    }
  };

  if (checking) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", minHeight: "100vh", background: "linear-gradient(135deg, #0f1117 0%, #1a1d2e 100%)" }}>
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
        background: "linear-gradient(135deg, #0f1117 0%, #1a1d2e 100%)",
        gap: 4,
        px: 3,
      }}
    >
      <Logo size="large" />
      <Button
        variant="contained"
        size="large"
        startIcon={<GoogleIcon />}
        onClick={handleSignIn}
        disabled={signingIn}
        sx={{ py: 1.5, px: 4, maxWidth: 320, width: "100%" }}
      >
        {signingIn ? "Signing in..." : "Sign in with Google"}
      </Button>
    </Box>
  );
}
