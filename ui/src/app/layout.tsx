import type { Metadata } from "next";
import { AppRouterCacheProvider } from "@mui/material-nextjs/v14-appRouter";
import { ThemeContextProvider } from "@/components/ThemeContext";
import { AuthProvider } from "@/components/AuthProvider";
import "./globals.css";

export const metadata: Metadata = {
  title: "LootForge",
  description: "Fair loot distribution for gaming groups",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <AppRouterCacheProvider>
          <ThemeContextProvider>
            <AuthProvider>{children}</AuthProvider>
          </ThemeContextProvider>
        </AppRouterCacheProvider>
      </body>
    </html>
  );
}
