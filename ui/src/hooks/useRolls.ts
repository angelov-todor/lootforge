"use client";
import { useState } from "react";
import { api } from "@/lib/api";
import { useGroup } from "@/components/GroupContext";
import type { RollResponse } from "@/types";

export function useRolls() {
  const { selectedGroup } = useGroup();
  const [rollResult, setRollResult] = useState<RollResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const executeRoll = async (participantIDs: string[], item: string) => {
    if (!selectedGroup) return;
    setError(null);
    setLoading(true);
    try {
      const result = await api.post<RollResponse>(
        `/api/groups/${selectedGroup.id}/rolls`,
        { participantIDs, item }
      );
      setRollResult(result);
    } catch (e) {
      console.error("Failed to execute roll", e);
      setError("Failed to execute roll");
    } finally {
      setLoading(false);
    }
  };

  const clearResult = () => setRollResult(null);

  return { executeRoll, rollResult, loading, error, clearResult };
}
