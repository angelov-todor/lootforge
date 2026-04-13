"use client";
import { useMemo } from "react";
import { Box, CircularProgress, Divider, Paper, Typography } from "@mui/material";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { useStats } from "@/hooks/useStats";

function stdDev(values: number[]): number {
  if (values.length === 0) return 0;
  const mean = values.reduce((a, b) => a + b, 0) / values.length;
  const variance = values.reduce((sum, v) => sum + (v - mean) ** 2, 0) / values.length;
  return Math.sqrt(variance);
}

export function Statistics() {
  const { stats, loading } = useStats();

  const entries = useMemo(() => Object.entries(stats), [stats]);

  const totalWins = useMemo(() => entries.reduce((sum, [, v]) => sum + v, 0), [entries]);
  const uniqueWinners = entries.length;

  const fairness = useMemo(() => {
    if (totalWins === 0) return 0;
    const percentages = entries.map(([, v]) => (v / totalWins) * 100);
    return stdDev(percentages);
  }, [entries, totalWins]);

  const chartData = useMemo(
    () => entries.map(([id, wins]) => ({ id, wins })),
    [entries]
  );

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", py: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (entries.length === 0) {
    return (
      <Typography color="text.secondary" sx={{ py: 2 }}>
        No statistics available yet.
      </Typography>
    );
  }

  return (
    <Box>
      {/* Summary */}
      <Paper variant="outlined" sx={{ p: 2, mb: 3, display: "flex", gap: 4 }}>
        <Box>
          <Typography variant="overline" color="text.secondary">
            Total Wins
          </Typography>
          <Typography variant="h4" sx={{ fontWeight: "bold" }}>
            {totalWins}
          </Typography>
        </Box>
        <Divider orientation="vertical" flexItem />
        <Box>
          <Typography variant="overline" color="text.secondary">
            Unique Winners
          </Typography>
          <Typography variant="h4" sx={{ fontWeight: "bold" }}>
            {uniqueWinners}
          </Typography>
        </Box>
        <Divider orientation="vertical" flexItem />
        <Box>
          <Typography variant="overline" color="text.secondary">
            Fairness (std dev %)
          </Typography>
          <Typography variant="h4" sx={{ fontWeight: "bold" }}>
            {fairness.toFixed(1)}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            Lower = more equal distribution
          </Typography>
        </Box>
      </Paper>

      {/* Bar chart */}
      <Typography variant="subtitle1" gutterBottom>
        Wins per Member
      </Typography>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={chartData} margin={{ top: 8, right: 16, left: 0, bottom: 8 }}>
          <XAxis dataKey="id" tick={{ fontSize: 12 }} />
          <YAxis allowDecimals={false} tick={{ fontSize: 12 }} />
          <Tooltip />
          <Bar dataKey="wins" fill="#7c3aed" radius={[4, 4, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </Box>
  );
}
