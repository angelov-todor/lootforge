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
