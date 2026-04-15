"use client";
import { useState, useCallback, useEffect } from "react";
import { api } from "@/lib/api";
import { useGroup } from "@/components/GroupContext";
import type { RollSession } from "@/types";

export function useRollHistory() {
  const { selectedGroup } = useGroup();
  const [rolls, setRolls] = useState<RollSession[]>([]);
  const [loading, setLoading] = useState(false);
  const [nextCursor, setNextCursor] = useState("");
  const [error, setError] = useState<string | null>(null);

  const fetchRolls = useCallback(
    async (cursor = "") => {
      if (!selectedGroup) return;
      setError(null);
      setLoading(true);
      try {
        const params = new URLSearchParams();
        if (cursor) params.set("cursor", cursor);
        const data = await api.get<{ rolls: RollSession[]; nextCursor: string }>(
          `/api/groups/${selectedGroup.id}/rolls?${params}`
        );
        if (cursor) {
          setRolls((prev) => [...prev, ...(data.rolls || [])]);
        } else {
          setRolls(data.rolls || []);
        }
        setNextCursor(data.nextCursor || "");
      } catch (e) {
        console.error("Failed to fetch roll history", e);
        setError("Failed to load roll history");
      } finally {
        setLoading(false);
      }
    },
    [selectedGroup]
  );

  useEffect(() => {
    fetchRolls();
  }, [fetchRolls]);

  const loadMore = () => {
    if (nextCursor) fetchRolls(nextCursor);
  };
  const hasMore = nextCursor !== "";

  return { rolls, loading, error, hasMore, loadMore, refresh: () => fetchRolls() };
}
