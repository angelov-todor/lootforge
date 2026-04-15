# UI Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign the LootForge UI from a single accordion-based dashboard into a mobile-first, multi-page app with bottom tab navigation, a roll-centric experience, full-screen winner reveal, and proper branding.

**Architecture:** Next.js App Router nested layouts. A shared `dashboard/layout.tsx` provides the top bar and bottom nav. Each tab is a separate page under `/dashboard/*`. Existing hooks and API client are reused unchanged. Components are rewritten for mobile-first card-based layouts.

**Tech Stack:** Next.js 14, React 18, Material-UI v9, Recharts, TypeScript

**Spec:** `docs/superpowers/specs/2026-04-15-ui-redesign-design.md`

---

## File Structure

### New Files
- `ui/src/components/BottomNav.tsx` — bottom tab bar with 5 tabs + elevated Roll FAB
- `ui/src/components/Logo.tsx` — anvil+die logo mark and wordmark
- `ui/src/components/WinnerOverlay.tsx` — full-screen winner reveal with animation
- `ui/src/app/dashboard/members/page.tsx` — members page
- `ui/src/app/dashboard/history/page.tsx` — history page
- `ui/src/app/dashboard/stats/page.tsx` — stats page
- `ui/src/app/dashboard/settings/page.tsx` — settings page

### Modified Files
- `ui/src/app/dashboard/layout.tsx` — add GroupProvider, TopBar, BottomNav wrapper
- `ui/src/app/dashboard/page.tsx` — rewrite as roll-station-only page
- `ui/src/app/page.tsx` — update login page with new logo
- `ui/src/components/TopBar.tsx` — slim version with Logo + group selector + avatar
- `ui/src/components/RollStation.tsx` — big-button + chips redesign
- `ui/src/components/MembersPanel.tsx` — card-based mobile list with FAB add
- `ui/src/components/RollHistory.tsx` — card-based list
- `ui/src/components/Statistics.tsx` — responsive mobile layout

### Unchanged Files (reused as-is)
- `ui/src/hooks/useMembers.ts`
- `ui/src/hooks/useRolls.ts`
- `ui/src/hooks/useRollHistory.ts`
- `ui/src/hooks/useStats.ts`
- `ui/src/components/AuthProvider.tsx`
- `ui/src/components/ProtectedRoute.tsx`
- `ui/src/components/GroupContext.tsx`
- `ui/src/components/GroupSelector.tsx`
- `ui/src/components/CreateGroupDialog.tsx`
- `ui/src/lib/api.ts`
- `ui/src/lib/firebase.ts`
- `ui/src/lib/theme.ts`
- `ui/src/types/index.ts`

---

### Task 1: Create Logo Component

**Files:**
- Create: `ui/src/components/Logo.tsx`

- [ ] **Step 1: Create the Logo component**

```tsx
// ui/src/components/Logo.tsx
"use client";
import { Box, Typography } from "@mui/material";

interface LogoProps {
  size?: "small" | "large";
}

export function Logo({ size = "small" }: LogoProps) {
  if (size === "large") {
    return (
      <Box sx={{ textAlign: "center" }}>
        <Box sx={{ fontSize: 48, position: "relative", display: "inline-block" }}>
          <Box
            component="span"
            sx={{ filter: "drop-shadow(0 0 8px rgba(124,58,237,0.5))" }}
          >
            &#9879;
          </Box>
          <Box
            component="span"
            sx={{
              position: "absolute",
              top: -8,
              right: -12,
              fontSize: 24,
              filter: "drop-shadow(0 0 6px rgba(245,158,11,0.5))",
            }}
          >
            &#127922;
          </Box>
        </Box>
        <Typography
          variant="h4"
          sx={{ fontWeight: 800, mt: 1 }}
        >
          <Box
            component="span"
            sx={{
              background: "linear-gradient(135deg, #a855f7, #f59e0b)",
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
            }}
          >
            LOOT
          </Box>
          <Box component="span" sx={{ color: "white" }}>
            FORGE
          </Box>
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
      <Box sx={{ fontSize: 20, position: "relative", display: "inline-block", lineHeight: 1 }}>
        <span>&#9879;</span>
        <Box
          component="span"
          sx={{ position: "absolute", top: -4, right: -6, fontSize: 10 }}
        >
          &#127922;
        </Box>
      </Box>
      <Typography variant="subtitle1" sx={{ fontWeight: 800, lineHeight: 1 }}>
        <Box
          component="span"
          sx={{ color: "primary.main" }}
        >
          LOOT
        </Box>
        <Box component="span" sx={{ color: "white" }}>
          FORGE
        </Box>
      </Typography>
    </Box>
  );
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add ui/src/components/Logo.tsx
git commit -m "feat(ui): add Logo component with anvil+die branding"
```

---

### Task 2: Create BottomNav Component

**Files:**
- Create: `ui/src/components/BottomNav.tsx`

- [ ] **Step 1: Create the BottomNav component**

```tsx
// ui/src/components/BottomNav.tsx
"use client";
import { usePathname, useRouter } from "next/navigation";
import { Box, ButtonBase, Typography } from "@mui/material";
import PeopleIcon from "@mui/icons-material/People";
import HistoryIcon from "@mui/icons-material/History";
import BarChartIcon from "@mui/icons-material/BarChart";
import SettingsIcon from "@mui/icons-material/Settings";
import CasinoIcon from "@mui/icons-material/Casino";

const tabs = [
  { label: "Members", icon: PeopleIcon, path: "/dashboard/members" },
  { label: "History", icon: HistoryIcon, path: "/dashboard/history" },
  { label: "Roll", icon: CasinoIcon, path: "/dashboard", isCenter: true },
  { label: "Stats", icon: BarChartIcon, path: "/dashboard/stats" },
  { label: "Settings", icon: SettingsIcon, path: "/dashboard/settings" },
];

export function BottomNav() {
  const pathname = usePathname();
  const router = useRouter();

  const isActive = (path: string) => {
    if (path === "/dashboard") return pathname === "/dashboard";
    return pathname.startsWith(path);
  };

  return (
    <Box
      component="nav"
      sx={{
        position: "fixed",
        bottom: 0,
        left: 0,
        right: 0,
        zIndex: 1200,
        bgcolor: "background.paper",
        borderTop: "1px solid",
        borderColor: "divider",
        display: "flex",
        justifyContent: "space-around",
        alignItems: "flex-end",
        pb: "env(safe-area-inset-bottom)",
        height: 64,
      }}
    >
      {tabs.map((tab) => {
        const active = isActive(tab.path);
        const Icon = tab.icon;

        if (tab.isCenter) {
          return (
            <ButtonBase
              key={tab.path}
              onClick={() => router.push(tab.path)}
              sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                mt: "-16px",
              }}
            >
              <Box
                sx={{
                  width: 56,
                  height: 56,
                  borderRadius: "50%",
                  background: "linear-gradient(135deg, #7c3aed, #a855f7)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  boxShadow: active
                    ? "0 4px 20px rgba(124,58,237,0.5)"
                    : "0 4px 15px rgba(124,58,237,0.3)",
                }}
              >
                <Icon sx={{ fontSize: 28, color: "white" }} />
              </Box>
              <Typography
                variant="caption"
                sx={{
                  mt: 0.25,
                  fontSize: "0.65rem",
                  fontWeight: active ? 700 : 400,
                  color: "primary.main",
                }}
              >
                {tab.label}
              </Typography>
            </ButtonBase>
          );
        }

        return (
          <ButtonBase
            key={tab.path}
            onClick={() => router.push(tab.path)}
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              py: 1,
              px: 1.5,
              minWidth: 0,
            }}
          >
            <Icon
              sx={{
                fontSize: 24,
                color: active ? "primary.main" : "text.secondary",
              }}
            />
            <Typography
              variant="caption"
              sx={{
                fontSize: "0.65rem",
                fontWeight: active ? 700 : 400,
                color: active ? "primary.main" : "text.secondary",
              }}
            >
              {tab.label}
            </Typography>
          </ButtonBase>
        );
      })}
    </Box>
  );
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add ui/src/components/BottomNav.tsx
git commit -m "feat(ui): add BottomNav with elevated Roll FAB"
```

---

### Task 3: Create WinnerOverlay Component

**Files:**
- Create: `ui/src/components/WinnerOverlay.tsx`

- [ ] **Step 1: Create the WinnerOverlay component**

```tsx
// ui/src/components/WinnerOverlay.tsx
"use client";
import { Box, Button, Typography, Backdrop } from "@mui/material";
import type { RollResponse } from "@/types";

const overlayKeyframes = `
@keyframes winnerScale {
  0% { transform: scale(0.3); opacity: 0; }
  60% { transform: scale(1.05); }
  100% { transform: scale(1); opacity: 1; }
}
@keyframes sparkle {
  0%, 100% { opacity: 0.2; transform: scale(0.8); }
  50% { opacity: 0.8; transform: scale(1.2); }
}
`;

interface WinnerOverlayProps {
  result: RollResponse;
  onPass: () => void;
  onDone: () => void;
  loading: boolean;
}

const sparkles = [
  { top: "10%", left: "15%", size: 20, delay: "0s" },
  { top: "8%", right: "18%", size: 16, delay: "0.3s" },
  { bottom: "25%", left: "12%", size: 14, delay: "0.6s" },
  { bottom: "15%", right: "14%", size: 18, delay: "0.2s" },
  { top: "35%", left: "8%", size: 12, delay: "0.8s" },
  { top: "20%", right: "8%", size: 15, delay: "0.5s" },
];

export function WinnerOverlay({ result, onPass, onDone, loading }: WinnerOverlayProps) {
  return (
    <Backdrop
      open
      sx={{
        zIndex: 1300,
        background: "radial-gradient(ellipse at center, rgba(124,58,237,0.25) 0%, rgba(15,23,42,0.95) 70%)",
        flexDirection: "column",
        gap: 1,
      }}
    >
      <style>{overlayKeyframes}</style>

      {/* Sparkles */}
      {sparkles.map((s, i) => (
        <Box
          key={i}
          sx={{
            position: "absolute",
            top: s.top,
            left: s.left,
            right: s.right,
            bottom: s.bottom,
            fontSize: s.size,
            animation: `sparkle 2s ease-in-out ${s.delay} infinite`,
            opacity: 0.3,
          }}
        >
          &#10022;
        </Box>
      ))}

      {/* Content */}
      <Box
        sx={{
          textAlign: "center",
          animation: "winnerScale 0.5s ease-out",
        }}
      >
        <Typography sx={{ fontSize: 48 }}>&#127881;</Typography>
        <Typography
          variant="overline"
          sx={{ letterSpacing: 3, color: "primary.main", fontSize: "0.75rem" }}
        >
          Winner
        </Typography>
        <Typography
          variant="h3"
          sx={{
            fontWeight: "bold",
            color: "white",
            textShadow: "0 0 30px rgba(124,58,237,0.5)",
            my: 1,
          }}
        >
          {result.winner.name}
        </Typography>
        {result.item && (
          <Typography sx={{ color: "secondary.main", fontSize: "1.1rem" }}>
            &#9876; {result.item}
          </Typography>
        )}
        <Box sx={{ display: "flex", gap: 1.5, justifyContent: "center", mt: 3 }}>
          <Button
            variant="outlined"
            color="secondary"
            onClick={onPass}
            disabled={loading}
            sx={{ px: 4, py: 1.2, borderWidth: 2 }}
          >
            &#8634; Pass
          </Button>
          <Button
            variant="contained"
            onClick={onDone}
            sx={{ px: 4, py: 1.2 }}
          >
            &#10003; Done
          </Button>
        </Box>
      </Box>
    </Backdrop>
  );
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add ui/src/components/WinnerOverlay.tsx
git commit -m "feat(ui): add WinnerOverlay with full-screen animation"
```

---

### Task 4: Rewrite Dashboard Layout

**Files:**
- Modify: `ui/src/app/dashboard/layout.tsx`

- [ ] **Step 1: Rewrite the dashboard layout with TopBar + BottomNav**

Replace the entire contents of `ui/src/app/dashboard/layout.tsx` with:

```tsx
"use client";
import { Box } from "@mui/material";
import { AuthProvider } from "@/components/AuthProvider";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import { GroupProvider } from "@/components/GroupContext";
import { TopBar } from "@/components/TopBar";
import { BottomNav } from "@/components/BottomNav";

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <ProtectedRoute>
        <GroupProvider>
          <TopBar />
          <Box
            component="main"
            sx={{
              pb: "80px", // space for bottom nav
              minHeight: "calc(100vh - 64px)",
            }}
          >
            {children}
          </Box>
          <BottomNav />
        </GroupProvider>
      </ProtectedRoute>
    </AuthProvider>
  );
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add ui/src/app/dashboard/layout.tsx
git commit -m "feat(ui): rewrite dashboard layout with bottom nav"
```

---

### Task 5: Rewrite TopBar (Slim Version)

**Files:**
- Modify: `ui/src/components/TopBar.tsx`

- [ ] **Step 1: Rewrite TopBar as slim header with Logo**

Replace the entire contents of `ui/src/components/TopBar.tsx` with:

```tsx
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
```

- [ ] **Step 2: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add ui/src/components/TopBar.tsx
git commit -m "feat(ui): slim TopBar with Logo and group selector"
```

---

### Task 6: Rewrite Roll Station Page (Big Button + Chips)

**Files:**
- Modify: `ui/src/app/dashboard/page.tsx`
- Modify: `ui/src/components/RollStation.tsx`

- [ ] **Step 1: Rewrite RollStation component with big-button + chips design**

Replace the entire contents of `ui/src/components/RollStation.tsx` with:

```tsx
"use client";
import { useState, useEffect } from "react";
import { Box, Chip, CircularProgress, TextField, Typography } from "@mui/material";
import { useMembers } from "@/hooks/useMembers";
import { useRolls } from "@/hooks/useRolls";
import { WinnerOverlay } from "@/components/WinnerOverlay";

export function RollStation() {
  const { members } = useMembers();
  const { executeRoll, rollResult, loading, clearResult } = useRolls();

  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [item, setItem] = useState("");

  useEffect(() => {
    if (!rollResult) {
      setSelected(new Set(members.map((m) => m.id)));
    }
  }, [members, rollResult]);

  const toggleMember = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const selectAll = () => setSelected(new Set(members.map((m) => m.id)));
  const clearAll = () => setSelected(new Set());

  const handleRoll = async () => {
    await executeRoll(Array.from(selected), item);
  };

  const handlePass = async () => {
    if (!rollResult) return;
    const remaining = Array.from(selected).filter((id) => id !== rollResult.winner.id);
    if (remaining.length === 0) return;
    await executeRoll(remaining, item);
  };

  const handleDone = () => {
    clearResult();
    setItem("");
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        px: 2,
        pt: 3,
        gap: 2,
        minHeight: "calc(100vh - 200px)",
      }}
    >
      {/* Winner overlay */}
      {rollResult && (
        <WinnerOverlay
          result={rollResult}
          onPass={handlePass}
          onDone={handleDone}
          loading={loading}
        />
      )}

      {/* Item input */}
      <TextField
        placeholder="Item name (optional)"
        value={item}
        onChange={(e) => setItem(e.target.value)}
        size="small"
        disabled={!!rollResult}
        sx={{
          width: "80%",
          maxWidth: 320,
          "& .MuiInputBase-input": { textAlign: "center" },
        }}
      />

      {/* Participants label */}
      <Typography
        variant="overline"
        sx={{ color: "text.secondary", letterSpacing: 1 }}
      >
        Participants
      </Typography>

      {/* Member chips */}
      <Box sx={{ display: "flex", flexWrap: "wrap", justifyContent: "center", gap: 0.75 }}>
        {members.map((m) => (
          <Chip
            key={m.id}
            label={m.name}
            onClick={() => toggleMember(m.id)}
            disabled={!!rollResult}
            color={selected.has(m.id) ? "primary" : "default"}
            variant={selected.has(m.id) ? "filled" : "outlined"}
            sx={{ fontWeight: selected.has(m.id) ? 600 : 400 }}
          />
        ))}
      </Box>

      {/* Select All / Clear */}
      <Box sx={{ display: "flex", gap: 1.5 }}>
        <Typography
          variant="caption"
          sx={{ color: "primary.main", cursor: "pointer" }}
          onClick={selectAll}
        >
          Select All
        </Typography>
        <Typography variant="caption" sx={{ color: "text.secondary" }}>
          |
        </Typography>
        <Typography
          variant="caption"
          sx={{ color: "text.secondary", cursor: "pointer" }}
          onClick={clearAll}
        >
          Clear
        </Typography>
      </Box>

      {/* Big Roll Button */}
      <Box sx={{ flex: 1, display: "flex", alignItems: "center", justifyContent: "center" }}>
        <Box
          component="button"
          onClick={handleRoll}
          disabled={selected.size === 0 || loading || !!rollResult}
          sx={{
            width: 140,
            height: 140,
            borderRadius: "50%",
            background: "linear-gradient(135deg, #7c3aed, #a855f7)",
            border: "none",
            cursor: selected.size === 0 || loading ? "not-allowed" : "pointer",
            opacity: selected.size === 0 || loading ? 0.5 : 1,
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            boxShadow: "0 0 40px rgba(124,58,237,0.3)",
            transition: "transform 0.15s, box-shadow 0.15s",
            "&:hover:not(:disabled)": {
              transform: "scale(1.05)",
              boxShadow: "0 0 50px rgba(124,58,237,0.5)",
            },
            "&:active:not(:disabled)": {
              transform: "scale(0.95)",
            },
          }}
        >
          {loading ? (
            <CircularProgress size={36} sx={{ color: "white" }} />
          ) : (
            <>
              <Box sx={{ fontSize: 36, lineHeight: 1 }}>&#127922;</Box>
              <Typography
                sx={{ color: "white", fontWeight: "bold", fontSize: 18, mt: 0.5 }}
              >
                ROLL
              </Typography>
            </>
          )}
        </Box>
      </Box>

      {/* Selection count */}
      <Typography variant="caption" sx={{ color: "text.secondary", pb: 2 }}>
        {selected.size} of {members.length} members selected
      </Typography>
    </Box>
  );
}
```

- [ ] **Step 2: Rewrite dashboard page to show only RollStation + InviteAcceptor**

Replace the entire contents of `ui/src/app/dashboard/page.tsx` with:

```tsx
"use client";
export const dynamic = "force-dynamic";
import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Snackbar, Alert } from "@mui/material";
import { useGroup } from "@/components/GroupContext";
import { RollStation } from "@/components/RollStation";
import { api } from "@/lib/api";

function InviteAcceptor() {
  const searchParams = useSearchParams();
  const { refreshGroups, setSelectedGroup } = useGroup();
  const [msg, setMsg] = useState<{ text: string; severity: "success" | "error" } | null>(null);

  useEffect(() => {
    const token = searchParams.get("invite");
    if (!token) return;

    const url = new URL(window.location.href);
    url.searchParams.delete("invite");
    window.history.replaceState({}, "", url.pathname);

    api.acceptInvite(token).then(async (group) => {
      await refreshGroups();
      setSelectedGroup(group);
      setMsg({ text: `Joined group "${group.name}"!`, severity: "success" });
    }).catch((e: unknown) => {
      setMsg({ text: e instanceof Error ? e.message : "Failed to accept invite", severity: "error" });
    });
  }, [searchParams, refreshGroups, setSelectedGroup]);

  return (
    <Snackbar open={!!msg} autoHideDuration={4000} onClose={() => setMsg(null)}>
      <Alert severity={msg?.severity} variant="filled" onClose={() => setMsg(null)}>
        {msg?.text}
      </Alert>
    </Snackbar>
  );
}

export default function DashboardPage() {
  return (
    <>
      <InviteAcceptor />
      <RollStation />
    </>
  );
}
```

- [ ] **Step 3: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add ui/src/components/RollStation.tsx ui/src/app/dashboard/page.tsx
git commit -m "feat(ui): redesign Roll Station with big button and chips"
```

---

### Task 7: Create Members Page

**Files:**
- Create: `ui/src/app/dashboard/members/page.tsx`
- Modify: `ui/src/components/MembersPanel.tsx`

- [ ] **Step 1: Rewrite MembersPanel as mobile-friendly card list**

Replace the entire contents of `ui/src/components/MembersPanel.tsx` with:

```tsx
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
  TextField,
  Typography,
  CircularProgress,
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
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
              <IconButton size="small" onClick={() => openEdit(m)}>
                <EditIcon fontSize="small" />
              </IconButton>
              <IconButton size="small" color="error" onClick={() => handleDelete(m)}>
                <DeleteIcon fontSize="small" />
              </IconButton>
            </CardContent>
          </Card>
        ))}
      </Box>

      {/* FAB */}
      <Fab
        color="primary"
        onClick={openAdd}
        sx={{ position: "fixed", bottom: 80, right: 16 }}
      >
        <AddIcon />
      </Fab>

      {/* Add/Edit Dialog */}
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
```

- [ ] **Step 2: Create the members page**

```tsx
// ui/src/app/dashboard/members/page.tsx
"use client";
import { MembersPanel } from "@/components/MembersPanel";

export default function MembersPage() {
  return <MembersPanel />;
}
```

- [ ] **Step 3: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add ui/src/components/MembersPanel.tsx ui/src/app/dashboard/members/page.tsx
git commit -m "feat(ui): mobile-friendly members page with card list"
```

---

### Task 8: Create History Page

**Files:**
- Create: `ui/src/app/dashboard/history/page.tsx`
- Modify: `ui/src/components/RollHistory.tsx`

- [ ] **Step 1: Rewrite RollHistory as card-based list**

Replace the entire contents of `ui/src/components/RollHistory.tsx` with:

```tsx
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
```

- [ ] **Step 2: Create the history page**

```tsx
// ui/src/app/dashboard/history/page.tsx
"use client";
import { RollHistory } from "@/components/RollHistory";

export default function HistoryPage() {
  return <RollHistory />;
}
```

- [ ] **Step 3: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add ui/src/components/RollHistory.tsx ui/src/app/dashboard/history/page.tsx
git commit -m "feat(ui): card-based roll history page"
```

---

### Task 9: Create Stats Page

**Files:**
- Create: `ui/src/app/dashboard/stats/page.tsx`
- Modify: `ui/src/components/Statistics.tsx`

- [ ] **Step 1: Update Statistics component for responsive mobile layout**

Replace the entire contents of `ui/src/components/Statistics.tsx` with:

```tsx
"use client";
import { useMemo } from "react";
import { Box, Card, CardContent, CircularProgress, Typography } from "@mui/material";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";
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
  const chartData = useMemo(() => entries.map(([id, wins]) => ({ id, wins })), [entries]);

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
          {/* Summary cards */}
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

          {/* Chart */}
          <Typography variant="subtitle2" sx={{ mb: 1 }}>
            Wins per Member
          </Typography>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={chartData} margin={{ top: 8, right: 8, left: -16, bottom: 8 }}>
              <XAxis dataKey="id" tick={{ fontSize: 11 }} />
              <YAxis allowDecimals={false} tick={{ fontSize: 11 }} />
              <Tooltip />
              <Bar dataKey="wins" fill="#7c3aed" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </>
      )}
    </Box>
  );
}
```

- [ ] **Step 2: Create the stats page**

```tsx
// ui/src/app/dashboard/stats/page.tsx
"use client";
import { Statistics } from "@/components/Statistics";

export default function StatsPage() {
  return <Statistics />;
}
```

- [ ] **Step 3: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add ui/src/components/Statistics.tsx ui/src/app/dashboard/stats/page.tsx
git commit -m "feat(ui): responsive stats page with summary cards"
```

---

### Task 10: Create Settings Page

**Files:**
- Create: `ui/src/app/dashboard/settings/page.tsx`

- [ ] **Step 1: Create the settings page with group config and invites**

```tsx
// ui/src/app/dashboard/settings/page.tsx
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
  TextField,
  Typography,
} from "@mui/material";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import { useGroup } from "@/components/GroupContext";
import { api } from "@/lib/api";
import type { Invite } from "@/types";

export default function SettingsPage() {
  const { selectedGroup, refreshGroups } = useGroup();
  const [invite, setInvite] = useState<Invite | null>(null);
  const [copied, setCopied] = useState(false);
  const [loading, setLoading] = useState(false);

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

      {/* Group info */}
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

      {/* Invite section */}
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

      {/* Danger zone */}
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
```

- [ ] **Step 2: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add ui/src/app/dashboard/settings/page.tsx
git commit -m "feat(ui): settings page with invites and group config"
```

---

### Task 11: Update Login Page with Logo

**Files:**
- Modify: `ui/src/app/page.tsx`

- [ ] **Step 1: Update login page to use Logo component**

Replace the entire contents of `ui/src/app/page.tsx` with:

```tsx
"use client";
export const dynamic = "force-dynamic";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Box, Button, CircularProgress } from "@mui/material";
import GoogleIcon from "@mui/icons-material/Google";
import { useAuth } from "@/components/AuthProvider";
import { Logo } from "@/components/Logo";

export default function LoginPage() {
  const { user, loading, signIn } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && user) {
      router.replace("/dashboard");
    }
  }, [user, loading, router]);

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", minHeight: "100vh" }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "100vh",
        background: "linear-gradient(135deg, #0f172a 0%, #1e1b4b 100%)",
        gap: 4,
        px: 3,
      }}
    >
      <Logo size="large" />
      <Button
        variant="contained"
        size="large"
        startIcon={<GoogleIcon />}
        onClick={signIn}
        sx={{ py: 1.5, px: 4, maxWidth: 320, width: "100%" }}
      >
        Sign in with Google
      </Button>
    </Box>
  );
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd ui && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add ui/src/app/page.tsx
git commit -m "feat(ui): update login page with Logo branding"
```

---

### Task 12: Remove Unused InviteDialog Import from TopBar

**Files:**
- Verify: `ui/src/components/TopBar.tsx` no longer imports InviteDialog (already done in Task 5)

- [ ] **Step 1: Verify the app builds end-to-end**

Run: `cd ui && npm run build`
Expected: Build succeeds. If there are errors related to unused imports or missing references, fix them.

- [ ] **Step 2: Final commit with any fixes**

```bash
git add -A
git commit -m "feat(ui): complete mobile-first multi-page redesign"
```

- [ ] **Step 3: Push to trigger deploy**

```bash
git push
```
