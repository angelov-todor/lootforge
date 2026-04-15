"use client";
import { useState } from "react";
import {
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Fab,
  IconButton,
  ListItemIcon,
  ListItemText,
  Menu,
  MenuItem,
  TextField,
  Typography,
  CircularProgress,
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import { useGroup } from "@/components/GroupContext";
import { useMembers } from "@/hooks/useMembers";
import type { Member } from "@/types";

interface MemberFormData {
  name: string;
  role: string;
  luck: number;
  priority: number;
  points: number;
}

const emptyForm: MemberFormData = { name: "", role: "member", luck: 0, priority: 0, points: 0 };

export function MembersPanel() {
  const { selectedGroup } = useGroup();
  const { members, loading, addMember, updateMember, deleteMember } = useMembers();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<MemberFormData>(emptyForm);
  const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);
  const [menuMember, setMenuMember] = useState<Member | null>(null);

  const strategyType = selectedGroup?.strategy.type;

  const strategyLabel =
    strategyType === "weighted_luck" ? "Luck" :
    strategyType === "round_robin" ? "Priority" :
    strategyType === "dkp" ? "Points" : null;

  const getStrategyValue = (m: Member) =>
    strategyType === "weighted_luck" ? m.luck :
    strategyType === "round_robin" ? m.priority :
    strategyType === "dkp" ? m.points : null;

  const openAdd = () => {
    setEditingId(null);
    setForm(emptyForm);
    setDialogOpen(true);
  };

  const openEdit = (m: Member) => {
    setEditingId(m.id);
    setForm({ name: m.name, role: m.role, luck: m.luck, priority: m.priority, points: m.points });
    setDialogOpen(true);
  };

  const handleSave = async () => {
    if (!form.name.trim()) return;
    if (editingId) {
      const existing = members.find((m) => m.id === editingId);
      if (existing) {
        await updateMember({ ...existing, ...form });
      }
    } else {
      await addMember(form);
    }
    setDialogOpen(false);
  };

  const handleDelete = async (m: Member) => {
    if (confirm(`Delete "${m.name}"?`)) {
      await deleteMember(m.id);
    }
  };

  if (loading && members.length === 0) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", py: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ px: 2, pt: 2, pb: 10 }}>
      <Typography variant="h6" sx={{ mb: 2, fontWeight: "bold" }}>
        Members
      </Typography>

      {members.length === 0 && (
        <Typography color="text.secondary" sx={{ textAlign: "center", py: 4 }}>
          No members yet. Tap + to add one.
        </Typography>
      )}

      <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
        {members.map((m) => (
          <Card key={m.id} variant="outlined">
            <CardContent sx={{ display: "flex", alignItems: "center", py: 1.5, "&:last-child": { pb: 1.5 } }}>
              <Box sx={{ flex: 1 }}>
                <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                  {m.name}
                </Typography>
                <Box sx={{ display: "flex", gap: 1, alignItems: "center", mt: 0.5 }}>
                  <Chip label={m.role || "member"} size="small" variant="outlined" />
                  {strategyLabel && (
                    <Typography variant="caption" color="text.secondary">
                      {strategyLabel}: {getStrategyValue(m)}
                    </Typography>
                  )}
                </Box>
              </Box>
              <IconButton
                onClick={(e) => { setMenuAnchor(e.currentTarget); setMenuMember(m); }}
                sx={{ ml: 0.5 }}
              >
                <MoreVertIcon />
              </IconButton>
            </CardContent>
          </Card>
        ))}
      </Box>

      {/* Context menu */}
      <Menu
        anchorEl={menuAnchor}
        open={Boolean(menuAnchor)}
        onClose={() => setMenuAnchor(null)}
      >
        <MenuItem onClick={() => { if (menuMember) openEdit(menuMember); setMenuAnchor(null); }}>
          <ListItemIcon><EditIcon fontSize="small" /></ListItemIcon>
          <ListItemText>Edit</ListItemText>
        </MenuItem>
        <MenuItem onClick={() => { if (menuMember) handleDelete(menuMember); setMenuAnchor(null); }}>
          <ListItemIcon><DeleteIcon fontSize="small" color="error" /></ListItemIcon>
          <ListItemText sx={{ color: "error.main" }}>Delete</ListItemText>
        </MenuItem>
      </Menu>

      {/* FAB */}
      <Fab
        color="primary"
        onClick={openAdd}
        sx={{ position: "fixed", bottom: 80, right: 16 }}
      >
        <AddIcon />
      </Fab>

      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} fullWidth maxWidth="xs">
        <DialogTitle>{editingId ? "Edit Member" : "Add Member"}</DialogTitle>
        <DialogContent sx={{ display: "flex", flexDirection: "column", gap: 2, pt: "16px !important" }}>
          <TextField
            label="Name"
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            required
            autoFocus
          />
          <TextField
            label="Role"
            value={form.role}
            onChange={(e) => setForm((f) => ({ ...f, role: e.target.value }))}
          />
          {strategyType === "weighted_luck" && (
            <TextField
              label="Luck"
              type="number"
              value={form.luck}
              onChange={(e) => setForm((f) => ({ ...f, luck: Number(e.target.value) }))}
            />
          )}
          {strategyType === "round_robin" && (
            <TextField
              label="Priority"
              type="number"
              value={form.priority}
              onChange={(e) => setForm((f) => ({ ...f, priority: Number(e.target.value) }))}
            />
          )}
          {strategyType === "dkp" && (
            <TextField
              label="Points"
              type="number"
              value={form.points}
              onChange={(e) => setForm((f) => ({ ...f, points: Number(e.target.value) }))}
            />
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>Cancel</Button>
          <Button variant="contained" onClick={handleSave} disabled={!form.name.trim()}>
            {editingId ? "Save" : "Add"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
