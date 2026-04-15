# LootForge UI Redesign — Mobile-First Multi-Page App

**Date:** 2026-04-15
**Status:** Approved

## Problem

The current UI is a single dashboard page with all features in collapsible accordions. This is not mobile-friendly, buries the core rolling experience, lacks a logo, and doesn't feel like a purpose-built gaming tool.

## Goals

- Make rolling the central, unmissable experience
- Mobile-first design with touch-friendly controls
- Separate pages for distinct features instead of accordions
- Add branding (logo, consistent identity)
- Improve discoverability and navigation

## Design Decisions

### Navigation: Bottom Tab Bar

Fixed bottom navigation with 5 tabs:

| Tab | Icon | Route | Notes |
|-----|------|-------|-------|
| Members | People icon | `/dashboard/members` | Member CRUD |
| History | Scroll/list icon | `/dashboard/history` | Roll history with pagination |
| **Roll** | Die icon (elevated FAB) | `/dashboard` | Center tab, visually prominent |
| Stats | Chart icon | `/dashboard/stats` | Win statistics, fairness |
| Settings | Gear icon | `/dashboard/settings` | Group config, invites, strategy |

The Roll tab uses an elevated circular button (56px) with a purple gradient and glow shadow, visually breaking out of the tab bar to signal it's the primary action.

**Top bar** is slim: logo (left), group selector dropdown (center), user avatar with menu (right).

### Route Structure

```
/                         — Login page (public)
/dashboard                — Roll Station (default, protected)
/dashboard/members        — Members management (protected)
/dashboard/history        — Roll history (protected)
/dashboard/stats          — Statistics (protected)
/dashboard/settings       — Group settings & invites (protected)
```

All `/dashboard/*` routes share a layout with the top bar and bottom tab bar. The `GroupProvider` and `AuthProvider` wrap the dashboard layout.

### Roll Station (Core Screen)

The default landing page after login. Layout from top to bottom:

1. **Item name input** — centered, optional, placeholder "Item name (optional)"
2. **"Participants" label** — small uppercase
3. **Member chips** — horizontal wrap, tap to toggle. Selected = purple fill, unselected = gray. Shows checkmark on selected.
4. **Select All / Clear** — small text links below chips
5. **Roll button** — large circle (140px), centered, purple gradient with glow shadow. Die emoji + "ROLL" text. Dominates the visual hierarchy.
6. **Selection count** — "4 of 5 members selected" below the button

The entire screen is designed for one-handed mobile use. The roll button sits in the natural thumb zone.

### Winner Reveal

Full-screen overlay triggered after a successful roll:

- **Background:** radial gradient from purple center to dark edges, dims the page behind
- **Content (centered):** Party popper emoji, "WINNER" label (small caps, purple), winner name (32px bold white with text shadow), item name (amber, with sword emoji)
- **Actions:** two buttons side by side:
  - "Pass" — outlined amber button, re-rolls with remaining participants (winner excluded)
  - "Done" — solid purple button, clears result and returns to roll station
- **Decorative sparkle dots** scattered at varying opacity for depth
- **Animation:** scale-up entrance with a slight bounce, confetti/particle CSS animation

### Logo & Branding

**Mark:** Anvil emoji (&#9879;) with a small die emoji (&#127922;) overlaid at top-right. Purple glow drop shadow on the anvil, amber glow on the die.

**Wordmark:** "LOOT" in purple-to-amber gradient, "FORGE" in white. Bold 800 weight.

**Top bar variant:** compact — small anvil+die icon (20px) next to "LOOTFORGE" text (14px).

**Login page:** large logo centered, same gradient background as current (`#0f172a` to `#1e1b4b`), with the Google sign-in button below.

### Members Page (`/dashboard/members`)

Mobile-optimized list replacing the dense Material React Table:

- Each member is a card showing: name, role badge, strategy-specific stat (luck/priority/points)
- Swipe or tap actions for edit/delete
- FAB "+" button at bottom-right for adding members
- Add/edit uses a bottom sheet or simple dialog (name, role, initial stat value)

### History Page (`/dashboard/history`)

- Card-based list of past rolls (not a table)
- Each card: date, item name, winner name, strategy type badge
- "Load More" button at bottom for pagination
- Empty state with illustration/message

### Stats Page (`/dashboard/stats`)

- Summary cards at top: Total Wins, Unique Winners, Fairness Score
- Bar chart (Recharts) showing wins per member
- Responsive — chart fills width on mobile

### Settings Page (`/dashboard/settings`)

- Group name (editable)
- Strategy configuration (type selector + dynamic fields)
- Invite section: generate link button, copy to clipboard
- Danger zone: delete group

### Theme (unchanged)

- Dark mode with `#0f172a` background, `#1e293b` paper
- Primary: `#7c3aed` (purple)
- Secondary: `#f59e0b` (amber)
- Buttons: no text transform, 8px border radius

### Mobile-First Principles

- Minimum 48px touch targets
- Bottom nav in thumb zone
- Cards instead of tables
- No nested scrolling or accordions
- Full-width layouts on mobile, max-width 600px centered on desktop

## Components to Create/Modify

### New Components
- `BottomNav.tsx` — bottom tab bar with 5 tabs, elevated Roll FAB
- `DashboardLayout.tsx` — shared layout (top bar + bottom nav + content area)
- `WinnerOverlay.tsx` — full-screen winner reveal with animation
- `MemberCard.tsx` — individual member card for the members list
- `RollHistoryCard.tsx` — individual roll history card
- `SettingsPage.tsx` — group settings page

### Modified Components
- `TopBar.tsx` — slim version with logo, group selector, avatar only
- `RollStation.tsx` — redesign to big-button + chips layout
- `MembersPanel.tsx` — convert from table to card list
- `RollHistory.tsx` — convert from table to card list
- `Statistics.tsx` — keep mostly the same, ensure responsive
- `GroupSelector.tsx` — compact for top bar use

### Pages to Create
- `dashboard/layout.tsx` — update with `BottomNav` and `DashboardLayout`
- `dashboard/members/page.tsx` — members page
- `dashboard/history/page.tsx` — history page
- `dashboard/stats/page.tsx` — stats page
- `dashboard/settings/page.tsx` — settings page

### Removed
- Accordion-based single-page dashboard layout
- `InviteDialog.tsx` — functionality moves to Settings page

## Out of Scope

- Offline support / PWA
- Push notifications
- Real-time multiplayer (WebSocket)
- Custom avatar uploads
- Dark/light mode toggle (staying dark only)
