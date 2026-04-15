"use client";
import { useState, useEffect } from "react";
import { Box, Chip, CircularProgress, TextField, Typography } from "@mui/material";
import { useMembers } from "@/hooks/useMembers";
import { useRolls } from "@/hooks/useRolls";
import { WinnerOverlay } from "@/components/WinnerOverlay";

export function RollStation() {
  const { members } = useMembers();
  const { executeRoll, rollResult, loading, clearResult } = useRolls();

  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [item, setItem] = useState("");

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
    setSelected(new Set(remaining));
    try {
      await executeRoll(remaining, item);
    } catch (e) {
      console.error("Pass roll failed:", e);
    }
  };

  const handleDone = () => {
    clearResult();
    setItem("");
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        px: 2,
        pt: 3,
        gap: 2,
        minHeight: "calc(100vh - 200px)",
      }}
    >
      {rollResult && (
        <WinnerOverlay
          result={rollResult}
          onPass={handlePass}
          onDone={handleDone}
          loading={loading}
        />
      )}

      <TextField
        placeholder="Item name (optional)"
        value={item}
        onChange={(e) => setItem(e.target.value)}
        size="small"
        disabled={!!rollResult}
        sx={{
          width: "80%",
          maxWidth: 320,
          "& .MuiInputBase-input": { textAlign: "center" },
        }}
      />

      <Typography
        variant="overline"
        sx={{ color: "text.secondary", letterSpacing: 1 }}
      >
        Participants
      </Typography>

      <Box sx={{ display: "flex", flexWrap: "wrap", justifyContent: "center", gap: 0.75 }}>
        {members.map((m) => (
          <Chip
            key={m.id}
            label={m.name}
            onClick={() => toggleMember(m.id)}
            disabled={!!rollResult}
            color={selected.has(m.id) ? "primary" : "default"}
            variant={selected.has(m.id) ? "filled" : "outlined"}
            sx={{ fontWeight: selected.has(m.id) ? 600 : 400 }}
          />
        ))}
      </Box>

      <Box sx={{ display: "flex", gap: 1.5 }}>
        <Typography
          variant="caption"
          sx={{ color: "primary.main", cursor: "pointer" }}
          onClick={selectAll}
        >
          Select All
        </Typography>
        <Typography variant="caption" sx={{ color: "text.secondary" }}>
          |
        </Typography>
        <Typography
          variant="caption"
          sx={{ color: "text.secondary", cursor: "pointer" }}
          onClick={clearAll}
        >
          Clear
        </Typography>
      </Box>

      <Box sx={{ flex: 1, display: "flex", alignItems: "center", justifyContent: "center" }}>
        <Box
          component="button"
          onClick={handleRoll}
          disabled={selected.size === 0 || loading || !!rollResult}
          sx={{
            width: 140,
            height: 140,
            borderRadius: "50%",
            background: "linear-gradient(135deg, #7c3aed, #a855f7)",
            border: "none",
            cursor: selected.size === 0 || loading ? "not-allowed" : "pointer",
            opacity: selected.size === 0 || loading ? 0.5 : 1,
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            boxShadow: "0 0 40px rgba(124,58,237,0.3)",
            transition: "transform 0.15s, box-shadow 0.15s",
            "&:hover:not(:disabled)": {
              transform: "scale(1.05)",
              boxShadow: "0 0 50px rgba(124,58,237,0.5)",
            },
            "&:active:not(:disabled)": {
              transform: "scale(0.95)",
            },
          }}
        >
          {loading ? (
            <CircularProgress size={36} sx={{ color: "white" }} />
          ) : (
            <>
              <Box sx={{ fontSize: 36, lineHeight: 1 }}>&#127922;</Box>
              <Typography
                sx={{ color: "white", fontWeight: "bold", fontSize: 18, mt: 0.5 }}
              >
                ROLL
              </Typography>
            </>
          )}
        </Box>
      </Box>

      <Typography variant="caption" sx={{ color: "text.secondary", pb: 2 }}>
        {selected.size} of {members.length} members selected
      </Typography>
    </Box>
  );
}
