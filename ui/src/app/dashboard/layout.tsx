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
              pb: "80px",
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
