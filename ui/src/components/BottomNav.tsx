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
