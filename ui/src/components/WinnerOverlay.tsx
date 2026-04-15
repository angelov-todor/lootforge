"use client";
import { Box, Button, Typography, Backdrop } from "@mui/material";
import type { RollResponse } from "@/types";

const overlayKeyframes = `
@keyframes winnerScale {
  0% { transform: scale(0.3); opacity: 0; }
  60% { transform: scale(1.05); }
  100% { transform: scale(1); opacity: 1; }
}
@keyframes sparkle {
  0%, 100% { opacity: 0.2; transform: scale(0.8); }
  50% { opacity: 0.8; transform: scale(1.2); }
}
`;

interface WinnerOverlayProps {
  result: RollResponse;
  onPass: () => void;
  onDone: () => void;
  loading: boolean;
}

const sparkles = [
  { top: "10%", left: "15%", size: 20, delay: "0s" },
  { top: "8%", right: "18%", size: 16, delay: "0.3s" },
  { bottom: "25%", left: "12%", size: 14, delay: "0.6s" },
  { bottom: "15%", right: "14%", size: 18, delay: "0.2s" },
  { top: "35%", left: "8%", size: 12, delay: "0.8s" },
  { top: "20%", right: "8%", size: 15, delay: "0.5s" },
];

export function WinnerOverlay({ result, onPass, onDone, loading }: WinnerOverlayProps) {
  return (
    <Backdrop
      open
      sx={{
        zIndex: 1300,
        background: "radial-gradient(ellipse at center, rgba(99,102,241,0.35) 0%, rgba(15,23,42,0.98) 60%)",
        flexDirection: "column",
        gap: 1,
      }}
    >
      <style>{overlayKeyframes}</style>

      {/* Sparkles */}
      {sparkles.map((s, i) => (
        <Box
          key={i}
          sx={{
            position: "absolute",
            top: s.top,
            left: s.left,
            right: s.right,
            bottom: s.bottom,
            fontSize: s.size,
            animation: `sparkle 2s ease-in-out ${s.delay} infinite`,
            opacity: 0.3,
          }}
        >
          &#10022;
        </Box>
      ))}

      {/* Content */}
      <Box
        sx={{
          textAlign: "center",
          animation: "winnerScale 0.5s ease-out",
        }}
      >
        <Typography sx={{ fontSize: 48 }}>&#127881;</Typography>
        <Typography
          variant="overline"
          sx={{ letterSpacing: 3, color: "primary.main", fontSize: "0.75rem" }}
        >
          Winner
        </Typography>
        <Typography
          variant="h3"
          sx={{
            fontWeight: "bold",
            color: "white",
            textShadow: "0 0 30px rgba(99,102,241,0.5)",
            my: 1,
          }}
        >
          {result.winner.name}
        </Typography>
        {result.item && (
          <Typography sx={{ color: "secondary.main", fontSize: "1.1rem" }}>
            &#9876; {result.item}
          </Typography>
        )}
        <Box sx={{ display: "flex", gap: 1.5, justifyContent: "center", mt: 3 }}>
          <Button
            variant="outlined"
            color="secondary"
            onClick={onPass}
            disabled={loading}
            sx={{ px: 4, py: 1.2, borderWidth: 2 }}
          >
            &#8634; Pass
          </Button>
          <Button
            variant="contained"
            onClick={onDone}
            sx={{ px: 4, py: 1.2 }}
          >
            &#10003; Done
          </Button>
        </Box>
      </Box>
    </Backdrop>
  );
}
