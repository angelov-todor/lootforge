"use client";
export const dynamic = "force-dynamic";
import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Box, Accordion, AccordionSummary, AccordionDetails, Typography, Snackbar, Alert } from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { GroupProvider, useGroup } from "@/components/GroupContext";
import { TopBar } from "@/components/TopBar";
import { MembersPanel } from "@/components/MembersPanel";
import { RollStation } from "@/components/RollStation";
import { RollHistory } from "@/components/RollHistory";
import { Statistics } from "@/components/Statistics";
import { api } from "@/lib/api";

function InviteAcceptor() {
  const searchParams = useSearchParams();
  const { refreshGroups, setSelectedGroup } = useGroup();
  const [msg, setMsg] = useState<{ text: string; severity: "success" | "error" } | null>(null);

  useEffect(() => {
    const token = searchParams.get("invite");
    if (!token) return;

    // Remove invite param from URL without reload
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

function DashboardContent() {
  return (
    <>
      <TopBar />
      <InviteAcceptor />
      <Box sx={{ maxWidth: 1200, mx: "auto", p: 2 }}>
        <Accordion defaultExpanded>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography variant="h6">Members</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <MembersPanel />
          </AccordionDetails>
        </Accordion>
        <Accordion defaultExpanded>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography variant="h6">Roll Station</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <RollStation />
          </AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography variant="h6">Roll History</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <RollHistory />
          </AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography variant="h6">Statistics</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <Statistics />
          </AccordionDetails>
        </Accordion>
      </Box>
    </>
  );
}

export default function DashboardPage() {
  return (
    <GroupProvider>
      <DashboardContent />
    </GroupProvider>
  );
}
