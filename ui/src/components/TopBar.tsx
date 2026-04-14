"use client";
import { useState } from "react";
import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  Avatar,
  IconButton,
  Menu,
  MenuItem,
  Tooltip,
} from "@mui/material";
import PersonAddIcon from "@mui/icons-material/PersonAdd";
import { useAuth } from "@/components/AuthProvider";
import { useGroup } from "@/components/GroupContext";
import { GroupSelector } from "@/components/GroupSelector";
import { InviteDialog } from "@/components/InviteDialog";

export function TopBar() {
  const { user, signOut } = useAuth();
  const { selectedGroup } = useGroup();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [inviteOpen, setInviteOpen] = useState(false);

  const handleAvatarClick = (e: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(e.currentTarget);
  };

  const handleClose = () => setAnchorEl(null);

  const handleSignOut = async () => {
    handleClose();
    await signOut();
  };

  return (
    <>
      <AppBar position="sticky" color="default" elevation={1}>
        <Toolbar>
          {/* Left: Brand */}
          <Typography variant="h6" color="primary" sx={{ mr: 2, fontWeight: "bold" }}>
            LootForge
          </Typography>

          {/* Center: Group Selector */}
          <Box sx={{ flexGrow: 1, display: "flex", justifyContent: "center", gap: 1, alignItems: "center" }}>
            <GroupSelector />
            {selectedGroup && (
              <Tooltip title="Invite to group">
                <IconButton size="small" onClick={() => setInviteOpen(true)}>
                  <PersonAddIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            )}
          </Box>

          {/* Right: User avatar */}
          <IconButton onClick={handleAvatarClick} size="small">
            <Avatar
              src={user?.photoURL ?? undefined}
              alt={user?.displayName ?? "User"}
              sx={{ width: 36, height: 36 }}
            />
          </IconButton>
          <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleClose}>
            <MenuItem disabled sx={{ fontSize: "0.85rem", color: "text.secondary" }}>
              {user?.displayName}
            </MenuItem>
            <MenuItem onClick={handleSignOut}>Sign Out</MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
      {selectedGroup && (
        <InviteDialog
          groupID={selectedGroup.id}
          open={inviteOpen}
          onClose={() => setInviteOpen(false)}
        />
      )}
    </>
  );
}
