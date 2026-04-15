"use client";
import { useState, useCallback, useEffect } from "react";
import { api } from "@/lib/api";
import { useGroup } from "@/components/GroupContext";

export function useStats() {
  const { selectedGroup } = useGroup();
  const [stats, setStats] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState(false);

  const fetchStats = useCallback(async () => {
    if (!selectedGroup) return;
    setLoading(true);
    try {
      const data = await api.get<{ wins: Record<string, number> }>(
        `/api/groups/${selectedGroup.id}/rolls/stats`
      );
      setStats(data?.wins || {});
    } finally {
      setLoading(false);
    }
  }, [selectedGroup]);

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  return { stats, loading, refresh: fetchStats };
}
