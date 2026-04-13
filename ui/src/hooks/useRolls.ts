"use client";
import { useState } from "react";
import { api } from "@/lib/api";
import { useGroup } from "@/components/GroupContext";
import type { RollResponse } from "@/types";

export function useRolls() {
  const { selectedGroup } = useGroup();
  const [rollResult, setRollResult] = useState<RollResponse | null>(null);
  const [loading, setLoading] = useState(false);

  const executeRoll = async (participantIDs: string[], item: string) => {
    if (!selectedGroup) return;
    setLoading(true);
    try {
      const result = await api.post<RollResponse>(
        `/api/groups/${selectedGroup.id}/rolls`,
        { participantIDs, item }
      );
      setRollResult(result);
    } finally {
      setLoading(false);
    }
  };

  const clearResult = () => setRollResult(null);

  return { executeRoll, rollResult, loading, clearResult };
}
