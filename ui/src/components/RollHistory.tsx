"use client";
import {
  Box,
  Button,
  CircularProgress,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";
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

  if (!loading && rolls.length === 0) {
    return (
      <Typography color="text.secondary" sx={{ py: 2 }}>
        No rolls recorded yet.
      </Typography>
    );
  }

  return (
    <Box>
      <TableContainer component={Paper} variant="outlined">
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Date</TableCell>
              <TableCell>Item</TableCell>
              <TableCell>Winner</TableCell>
              <TableCell>Strategy</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rolls.map((roll) => (
              <TableRow key={roll.id} hover>
                <TableCell>{formatDate(roll.createdAt)}</TableCell>
                <TableCell>{roll.item || "—"}</TableCell>
                <TableCell>{roll.winnerID}</TableCell>
                <TableCell>{roll.strategy.type}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      {hasMore && (
        <Box sx={{ mt: 2, display: "flex", justifyContent: "center" }}>
          <Button variant="outlined" onClick={loadMore} disabled={loading}>
            {loading ? "Loading…" : "Load More"}
          </Button>
        </Box>
      )}
    </Box>
  );
}
