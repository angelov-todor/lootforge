"use client";
import { useState, useCallback, useEffect } from "react";
import { api } from "@/lib/api";
import { useGroup } from "@/components/GroupContext";

export function useStats() {
  const { selectedGroup } = useGroup();
  const [stats, setStats] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = useCallback(async () => {
    if (!selectedGroup) return;
    setError(null);
    setLoading(true);
    try {
      const data = await api.get<{ wins: Record<string, number> }>(
        `/api/groups/${selectedGroup.id}/rolls/stats`
      );
      setStats(data?.wins || {});
    } catch (e) {
      console.error("Failed to fetch stats", e);
      setError("Failed to load stats");
    } finally {
      setLoading(false);
    }
  }, [selectedGroup]);

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  return { stats, loading, error, refresh: fetchStats };
}
