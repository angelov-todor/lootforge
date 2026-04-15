"use client";
import { useState } from "react";
import {
  Box,
  Button,
  Card,
  CardContent,
  Divider,
  IconButton,
  Snackbar,
  Alert,
  Switch,
  TextField,
  Typography,
} from "@mui/material";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import DarkModeIcon from "@mui/icons-material/DarkMode";
import LightModeIcon from "@mui/icons-material/LightMode";
import { useGroup } from "@/components/GroupContext";
import { useThemeMode } from "@/components/ThemeContext";
import { api } from "@/lib/api";
import type { Invite } from "@/types";

export default function SettingsPage() {
  const { selectedGroup, refreshGroups } = useGroup();
  const [invite, setInvite] = useState<Invite | null>(null);
  const [copied, setCopied] = useState(false);
  const [loading, setLoading] = useState(false);
  const { mode, toggleTheme } = useThemeMode();

  if (!selectedGroup) {
    return (
      <Box sx={{ px: 2, pt: 4, textAlign: "center" }}>
        <Typography color="text.secondary">Select a group first.</Typography>
      </Box>
    );
  }

  const handleGenerateInvite = async () => {
    setLoading(true);
    try {
      const inv = await api.createInvite(selectedGroup.id);
      setInvite(inv);
    } finally {
      setLoading(false);
    }
  };

  const inviteLink = invite
    ? `${window.location.origin}/dashboard?invite=${invite.token}`
    : "";

  const handleCopy = () => {
    navigator.clipboard.writeText(inviteLink);
    setCopied(true);
  };

  const handleDelete = async () => {
    if (!confirm(`Delete group "${selectedGroup.name}"? This cannot be undone.`)) return;
    await api.del(`/api/groups/${selectedGroup.id}`);
    await refreshGroups();
  };

  return (
    <Box sx={{ px: 2, pt: 2, pb: 4 }}>
      <Typography variant="h6" sx={{ mb: 2, fontWeight: "bold" }}>
        Settings
      </Typography>

      {/* Theme toggle */}
      <Card variant="outlined" sx={{ mb: 2 }}>
        <CardContent sx={{ display: "flex", alignItems: "center", justifyContent: "space-between", py: 1.5, "&:last-child": { pb: 1.5 } }}>
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            {mode === "dark" ? <DarkModeIcon /> : <LightModeIcon />}
            <Typography>
              {mode === "dark" ? "Dark" : "Light"} Mode
            </Typography>
          </Box>
          <Switch checked={mode === "dark"} onChange={toggleTheme} />
        </CardContent>
      </Card>

      <Card variant="outlined" sx={{ mb: 2 }}>
        <CardContent>
          <Typography variant="overline" color="text.secondary">
            Group
          </Typography>
          <Typography variant="h6" sx={{ fontWeight: 600 }}>
            {selectedGroup.name}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
            Strategy: {selectedGroup.strategy.type.replace("_", " ")}
          </Typography>
        </CardContent>
      </Card>

      <Card variant="outlined" sx={{ mb: 2 }}>
        <CardContent>
          <Typography variant="overline" color="text.secondary">
            Invite Members
          </Typography>
          {!invite ? (
            <Box sx={{ mt: 1 }}>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1.5 }}>
                Generate a link to invite others to this group.
              </Typography>
              <Button
                variant="contained"
                onClick={handleGenerateInvite}
                disabled={loading}
                fullWidth
              >
                Generate Invite Link
              </Button>
            </Box>
          ) : (
            <Box sx={{ mt: 1 }}>
              <Box sx={{ display: "flex", gap: 1, alignItems: "center" }}>
                <TextField
                  value={inviteLink}
                  size="small"
                  fullWidth
                  slotProps={{ input: { readOnly: true } }}
                />
                <IconButton onClick={handleCopy}>
                  <ContentCopyIcon />
                </IconButton>
              </Box>
              <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: "block" }}>
                Expires: {new Date(invite.expiresAt).toLocaleDateString()}
              </Typography>
            </Box>
          )}
        </CardContent>
      </Card>

      <Card variant="outlined" sx={{ borderColor: "error.main" }}>
        <CardContent>
          <Typography variant="overline" color="error">
            Danger Zone
          </Typography>
          <Divider sx={{ my: 1 }} />
          <Button
            variant="outlined"
            color="error"
            onClick={handleDelete}
            fullWidth
          >
            Delete Group
          </Button>
        </CardContent>
      </Card>

      <Snackbar open={copied} autoHideDuration={2000} onClose={() => setCopied(false)}>
        <Alert severity="success" variant="filled">Link copied!</Alert>
      </Snackbar>
    </Box>
  );
}
