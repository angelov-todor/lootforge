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

  const fetchRolls = useCallback(
    async (cursor = "") => {
      if (!selectedGroup) return;
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

  return { rolls, loading, hasMore, loadMore, refresh: () => fetchRolls() };
}
