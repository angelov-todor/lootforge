"use client";
import { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  IconButton,
  Typography,
  Box,
  Snackbar,
  Alert,
} from "@mui/material";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import { api } from "@/lib/api";
import type { Invite } from "@/types";

interface InviteDialogProps {
  groupID: string;
  open: boolean;
  onClose: () => void;
}

export function InviteDialog({ groupID, open, onClose }: InviteDialogProps) {
  const [invite, setInvite] = useState<Invite | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  const handleGenerate = async () => {
    setLoading(true);
    setError("");
    try {
      const inv = await api.createInvite(groupID);
      setInvite(inv);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to create invite");
    } finally {
      setLoading(false);
    }
  };

  const inviteLink = invite
    ? `${window.location.origin}/dashboard?invite=${invite.token}`
    : "";

  const handleCopy = async () => {
    await navigator.clipboard.writeText(inviteLink);
    setCopied(true);
  };

  const handleClose = () => {
    setInvite(null);
    setError("");
    onClose();
  };

  return (
    <>
      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>Invite to Group</DialogTitle>
        <DialogContent>
          {!invite ? (
            <Typography sx={{ mt: 1 }}>
              Generate a link that others can use to join this group. The link expires in 7 days.
            </Typography>
          ) : (
            <Box sx={{ mt: 1 }}>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                Share this link with the person you want to invite:
              </Typography>
              <Box sx={{ display: "flex", gap: 1, alignItems: "center" }}>
                <TextField
                  fullWidth
                  size="small"
                  value={inviteLink}
                  slotProps={{ input: { readOnly: true } }}
                />
                <IconButton onClick={handleCopy} size="small" title="Copy link">
                  <ContentCopyIcon fontSize="small" />
                </IconButton>
              </Box>
              <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: "block" }}>
                Expires: {new Date(invite.expiresAt).toLocaleDateString()}
              </Typography>
            </Box>
          )}
          {error && (
            <Typography color="error" sx={{ mt: 1 }}>
              {error}
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          {!invite ? (
            <>
              <Button onClick={handleClose}>Cancel</Button>
              <Button onClick={handleGenerate} variant="contained" disabled={loading}>
                {loading ? "Generating…" : "Generate Link"}
              </Button>
            </>
          ) : (
            <Button onClick={handleClose}>Done</Button>
          )}
        </DialogActions>
      </Dialog>
      <Snackbar open={copied} autoHideDuration={2000} onClose={() => setCopied(false)}>
        <Alert severity="success" variant="filled" onClose={() => setCopied(false)}>
          Link copied to clipboard
        </Alert>
      </Snackbar>
    </>
  );
}
