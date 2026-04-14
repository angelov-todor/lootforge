"use client";
import { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Box,
  type SelectChangeEvent,
} from "@mui/material";
import { api } from "@/lib/api";
import { useGroup } from "@/components/GroupContext";
import type { StrategyType, StrategyConfig } from "@/types";

interface CreateGroupDialogProps {
  open: boolean;
  onClose: () => void;
}

export function CreateGroupDialog({ open, onClose }: CreateGroupDialogProps) {
  const { refreshGroups, setSelectedGroup } = useGroup();
  const [name, setName] = useState("");
  const [strategyType, setStrategyType] = useState<StrategyType>("pure_random");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // Weighted Luck config
  const [increment, setIncrement] = useState(10);
  const [decrement, setDecrement] = useState(5);
  const [defaultLuck, setDefaultLuck] = useState(50);

  // DKP config
  const [winCost, setWinCost] = useState(50);

  const handleSubmit = async () => {
    const trimmed = name.trim();
    if (!trimmed) {
      setError("Name is required");
      return;
    }

    setLoading(true);
    setError("");

    const strategy: StrategyConfig = { type: strategyType };
    if (strategyType === "weighted_luck") {
      strategy.weightedLuck = { increment, decrement, defaultLuck };
    } else if (strategyType === "dkp") {
      strategy.dkp = { winCost };
    }

    try {
      const group = await api.createGroup(trimmed, strategy);
      await refreshGroups();
      setSelectedGroup(group);
      handleClose();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to create group");
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setName("");
    setStrategyType("pure_random");
    setError("");
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>Create New Group</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          fullWidth
          label="Group Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          sx={{ mt: 1, mb: 2 }}
          error={!!error}
          helperText={error}
        />
        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel>Loot Strategy</InputLabel>
          <Select
            value={strategyType}
            label="Loot Strategy"
            onChange={(e: SelectChangeEvent<StrategyType>) =>
              setStrategyType(e.target.value as StrategyType)
            }
          >
            <MenuItem value="pure_random">Pure Random</MenuItem>
            <MenuItem value="weighted_luck">Weighted Luck</MenuItem>
            <MenuItem value="round_robin">Round Robin</MenuItem>
            <MenuItem value="dkp">DKP</MenuItem>
          </Select>
        </FormControl>

        {strategyType === "weighted_luck" && (
          <Box sx={{ display: "flex", gap: 2 }}>
            <TextField
              label="Increment"
              type="number"
              size="small"
              value={increment}
              onChange={(e) => setIncrement(Number(e.target.value))}
            />
            <TextField
              label="Decrement"
              type="number"
              size="small"
              value={decrement}
              onChange={(e) => setDecrement(Number(e.target.value))}
            />
            <TextField
              label="Default Luck"
              type="number"
              size="small"
              value={defaultLuck}
              onChange={(e) => setDefaultLuck(Number(e.target.value))}
            />
          </Box>
        )}

        {strategyType === "dkp" && (
          <TextField
            label="Win Cost"
            type="number"
            size="small"
            value={winCost}
            onChange={(e) => setWinCost(Number(e.target.value))}
          />
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button onClick={handleSubmit} variant="contained" disabled={loading}>
          {loading ? "Creating…" : "Create"}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
