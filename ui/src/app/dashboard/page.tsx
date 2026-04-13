"use client";
export const dynamic = "force-dynamic";
import { Box, Accordion, AccordionSummary, AccordionDetails, Typography } from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { GroupProvider } from "@/components/GroupContext";
import { TopBar } from "@/components/TopBar";
import { MembersPanel } from "@/components/MembersPanel";
import { RollStation } from "@/components/RollStation";
import { RollHistory } from "@/components/RollHistory";
import { Statistics } from "@/components/Statistics";

export default function DashboardPage() {
  return (
    <GroupProvider>
      <TopBar />
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
    </GroupProvider>
  );
}
