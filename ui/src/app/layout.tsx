import type { Metadata } from "next";
import { AppRouterCacheProvider } from "@mui/material-nextjs/v14-appRouter";
import { ThemeContextProvider } from "@/components/ThemeContext";
import "./globals.css";

export const metadata: Metadata = {
  title: "LootForge",
  description: "Fair loot distribution for gaming groups",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://apis.google.com" />
        <link rel="preconnect" href="https://www.googleapis.com" />
      </head>
      <body>
        <AppRouterCacheProvider>
          <ThemeContextProvider>
            {children}
          </ThemeContextProvider>
        </AppRouterCacheProvider>
      </body>
    </html>
  );
}
