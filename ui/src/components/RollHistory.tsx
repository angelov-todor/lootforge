"use client";
import { Box, Button, Card, CardContent, Chip, CircularProgress, Typography } from "@mui/material";
import { useRollHistory } from "@/hooks/useRollHistory";

function formatDate(iso: string): string {
  try {
    return new Date(iso).toLocaleString();
  } catch {
    return iso;
  }
}

export function RollHistory() {
  const { rolls, loading, hasMore, loadMore } = useRollHistory();

  if (loading && rolls.length === 0) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", py: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ px: 2, pt: 2, pb: 4 }}>
      <Typography variant="h6" sx={{ mb: 2, fontWeight: "bold" }}>
        Roll History
      </Typography>

      {!loading && rolls.length === 0 && (
        <Typography color="text.secondary" sx={{ textAlign: "center", py: 4 }}>
          No rolls recorded yet.
        </Typography>
      )}

      <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
        {rolls.map((roll) => (
          <Card key={roll.id} variant="outlined">
            <CardContent sx={{ py: 1.5, "&:last-child": { pb: 1.5 } }}>
              <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                  {roll.item || "No item"}
                </Typography>
                <Chip label={roll.strategy.type.replace("_", " ")} size="small" variant="outlined" />
              </Box>
              <Box sx={{ display: "flex", justifyContent: "space-between", mt: 0.5 }}>
                <Typography variant="body2" color="text.secondary">
                  Winner: <Box component="span" sx={{ color: "primary.main", fontWeight: 600 }}>{roll.winnerID}</Box>
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {formatDate(roll.createdAt)}
                </Typography>
              </Box>
            </CardContent>
          </Card>
        ))}
      </Box>

      {hasMore && (
        <Box sx={{ mt: 2, display: "flex", justifyContent: "center" }}>
          <Button variant="outlined" onClick={loadMore} disabled={loading}>
            {loading ? "Loading..." : "Load More"}
          </Button>
        </Box>
      )}
    </Box>
  );
}
