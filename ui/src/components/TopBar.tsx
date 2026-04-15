"use client";
import { useState } from "react";
import {
  AppBar,
  Toolbar,
  Box,
  Avatar,
  IconButton,
  Menu,
  MenuItem,
} from "@mui/material";
import { useAuth } from "@/components/AuthProvider";
import { GroupSelector } from "@/components/GroupSelector";
import { Logo } from "@/components/Logo";

export function TopBar() {
  const { user, signOut } = useAuth();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleClose = () => setAnchorEl(null);

  const handleSignOut = async () => {
    handleClose();
    await signOut();
  };

  return (
    <AppBar position="sticky" color="default" elevation={1}>
      <Toolbar sx={{ minHeight: { xs: 56 }, px: { xs: 1, sm: 2 } }}>
        <Logo />
        <Box sx={{ flexGrow: 1, display: "flex", justifyContent: "center" }}>
          <GroupSelector />
        </Box>
        <IconButton onClick={(e) => setAnchorEl(e.currentTarget)} size="small">
          <Avatar
            src={user?.photoURL ?? undefined}
            alt={user?.displayName ?? "User"}
            sx={{ width: 32, height: 32 }}
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
  );
}
