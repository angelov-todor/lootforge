"use client";
import { useState, useEffect } from "react";
import {
  Box,
  Button,
  Card,
  CardContent,
  Checkbox,
  FormControlLabel,
  TextField,
  Typography,
  CircularProgress,
} from "@mui/material";
import { useMembers } from "@/hooks/useMembers";
import { useRolls } from "@/hooks/useRolls";

const winnerKeyframes = `
@keyframes winnerAppear {
  0% { transform: scale(0.6); opacity: 0; box-shadow: 0 0 0px #7c3aed; }
  60% { transform: scale(1.08); box-shadow: 0 0 40px #7c3aed; }
  100% { transform: scale(1); opacity: 1; box-shadow: 0 0 20px #7c3aed88; }
}
`;

export function RollStation() {
  const { members } = useMembers();
  const { executeRoll, rollResult, loading, clearResult } = useRolls();

  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [item, setItem] = useState("");

  // When members change or result clears, reset selection
  useEffect(() => {
    if (!rollResult) {
      setSelected(new Set(members.map((m) => m.id)));
    }
  }, [members, rollResult]);

  const toggleMember = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const selectAll = () => setSelected(new Set(members.map((m) => m.id)));
  const clearAll = () => setSelected(new Set());

  const handleRoll = async () => {
    await executeRoll(Array.from(selected), item);
  };

  const handlePass = async () => {
    if (!rollResult) return;
    const remaining = Array.from(selected).filter((id) => id !== rollResult.winner.id);
    if (remaining.length === 0) return;
    await executeRoll(remaining, item);
  };

  const handleDone = () => {
    clearResult();
    setItem("");
  };

  return (
    <Box>
      <style>{winnerKeyframes}</style>

      {/* Winner display */}
      {rollResult && (
        <Card
          sx={{
            mb: 3,
            border: "2px solid",
            borderColor: "primary.main",
            animation: "winnerAppear 0.5s ease-out",
            textAlign: "center",
          }}
        >
          <CardContent>
            <Typography variant="overline" color="text.secondary">
              Winner
            </Typography>
            <Typography variant="h4" color="primary" sx={{ my: 1, fontWeight: "bold" }}>
              {rollResult.winner.name}
            </Typography>
            {rollResult.item && (
              <Typography variant="body2" color="text.secondary">
                Item: {rollResult.item}
              </Typography>
            )}
            <Box sx={{ mt: 2, display: "flex", gap: 2, justifyContent: "center" }}>
              <Button variant="outlined" onClick={handlePass} disabled={loading}>
                Pass
              </Button>
              <Button variant="contained" onClick={handleDone}>
                Done
              </Button>
            </Box>
          </CardContent>
        </Card>
      )}

      {/* Item field */}
      <TextField
        label="Item (optional)"
        value={item}
        onChange={(e) => setItem(e.target.value)}
        size="small"
        sx={{ mb: 2, width: 300 }}
        disabled={!!rollResult}
      />

      {/* Participant list */}
      <Box sx={{ mb: 2 }}>
        <Box sx={{ display: "flex", gap: 1, mb: 1 }}>
          <Button size="small" variant="outlined" onClick={selectAll} disabled={!!rollResult}>
            Select All
          </Button>
          <Button size="small" variant="outlined" onClick={clearAll} disabled={!!rollResult}>
            Clear
          </Button>
        </Box>
        <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
          {members.map((m) => (
            <FormControlLabel
              key={m.id}
              control={
                <Checkbox
                  checked={selected.has(m.id)}
                  onChange={() => toggleMember(m.id)}
                  size="small"
                  disabled={!!rollResult}
                />
              }
              label={m.name}
            />
          ))}
        </Box>
      </Box>

      {/* Roll button */}
      {!rollResult && (
        <Button
          variant="contained"
          size="large"
          onClick={handleRoll}
          disabled={selected.size === 0 || loading}
          sx={{ minWidth: 160, py: 1.5, fontSize: "1.1rem" }}
          startIcon={loading ? <CircularProgress size={18} color="inherit" /> : undefined}
        >
          {loading ? "Rolling…" : "Roll"}
        </Button>
      )}
    </Box>
  );
}
