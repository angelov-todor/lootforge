"use client";
import { useMemo } from "react";
import { Box, Card, CardContent, CircularProgress, Typography } from "@mui/material";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";
import { useStats } from "@/hooks/useStats";
import { useMembers } from "@/hooks/useMembers";

function stdDev(values: number[]): number {
  if (values.length === 0) return 0;
  const mean = values.reduce((a, b) => a + b, 0) / values.length;
  const variance = values.reduce((sum, v) => sum + (v - mean) ** 2, 0) / values.length;
  return Math.sqrt(variance);
}

export function Statistics() {
  const { stats, loading } = useStats();
  const { members } = useMembers();

  const memberNames = useMemo(() => {
    const map: Record<string, string> = {};
    for (const m of members) {
      map[m.id] = m.name;
    }
    return map;
  }, [members]);

  const entries = useMemo(() => Object.entries(stats), [stats]);
  const totalWins = useMemo(() => entries.reduce((sum, [, v]) => sum + v, 0), [entries]);
  const uniqueWinners = entries.length;
  const fairness = useMemo(() => {
    if (totalWins === 0) return 0;
    const percentages = entries.map(([, v]) => (v / totalWins) * 100);
    return stdDev(percentages);
  }, [entries, totalWins]);
  const chartData = useMemo(
    () => entries.map(([id, wins]) => ({ name: memberNames[id] || id, wins })),
    [entries, memberNames],
  );

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", py: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ px: 2, pt: 2, pb: 4 }}>
      <Typography variant="h6" sx={{ mb: 2, fontWeight: "bold" }}>
        Statistics
      </Typography>

      {entries.length === 0 ? (
        <Typography color="text.secondary" sx={{ textAlign: "center", py: 4 }}>
          No statistics available yet.
        </Typography>
      ) : (
        <>
          <Box sx={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 1, mb: 3 }}>
            <Card variant="outlined">
              <CardContent sx={{ textAlign: "center", py: 1.5, "&:last-child": { pb: 1.5 } }}>
                <Typography variant="overline" color="text.secondary" sx={{ fontSize: "0.65rem" }}>
                  Total Wins
                </Typography>
                <Typography variant="h5" sx={{ fontWeight: "bold" }}>
                  {totalWins}
                </Typography>
              </CardContent>
            </Card>
            <Card variant="outlined">
              <CardContent sx={{ textAlign: "center", py: 1.5, "&:last-child": { pb: 1.5 } }}>
                <Typography variant="overline" color="text.secondary" sx={{ fontSize: "0.65rem" }}>
                  Winners
                </Typography>
                <Typography variant="h5" sx={{ fontWeight: "bold" }}>
                  {uniqueWinners}
                </Typography>
              </CardContent>
            </Card>
            <Card variant="outlined">
              <CardContent sx={{ textAlign: "center", py: 1.5, "&:last-child": { pb: 1.5 } }}>
                <Typography variant="overline" color="text.secondary" sx={{ fontSize: "0.65rem" }}>
                  Fairness
                </Typography>
                <Typography variant="h5" sx={{ fontWeight: "bold" }}>
                  {fairness.toFixed(1)}
                </Typography>
              </CardContent>
            </Card>
          </Box>

          <Typography variant="subtitle2" sx={{ mb: 1 }}>
            Wins per Member
          </Typography>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={chartData} margin={{ top: 8, right: 8, left: -16, bottom: 8 }}>
              <XAxis dataKey="name" tick={{ fontSize: 11, fill: "#94a3b8" }} />
              <YAxis allowDecimals={false} tick={{ fontSize: 11, fill: "#94a3b8" }} />
              <Tooltip
                cursor={{ fill: "rgba(124,58,237,0.15)" }}
                contentStyle={{
                  backgroundColor: "#1e293b",
                  border: "1px solid #334155",
                  borderRadius: 8,
                  color: "#e2e8f0",
                }}
              />
              <Bar dataKey="wins" fill="#7c3aed" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </>
      )}
    </Box>
  );
}
