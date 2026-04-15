"use client";
import { useState, useCallback, useEffect } from "react";
import { api } from "@/lib/api";
import { useGroup } from "@/components/GroupContext";
import type { Member } from "@/types";
import { saveAs } from "file-saver";

export function useMembers() {
  const { selectedGroup } = useGroup();
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refreshMembers = useCallback(async () => {
    if (!selectedGroup) return;
    setError(null);
    setLoading(true);
    try {
      const data = await api.get<Member[]>(`/api/groups/${selectedGroup.id}/members`);
      setMembers(data || []);
    } catch (e) {
      console.error("Failed to fetch members", e);
      setError("Failed to load members");
    } finally {
      setLoading(false);
    }
  }, [selectedGroup]);

  useEffect(() => {
    refreshMembers();
  }, [refreshMembers]);

  const addMember = async (member: Partial<Member>) => {
    if (!selectedGroup) return;
    await api.post(`/api/groups/${selectedGroup.id}/members`, member);
    await refreshMembers();
  };

  const updateMember = async (member: Member) => {
    if (!selectedGroup) return;
    await api.put(`/api/groups/${selectedGroup.id}/members/${member.id}`, member);
    await refreshMembers();
  };

  const deleteMember = async (id: string) => {
    if (!selectedGroup) return;
    await api.del(`/api/groups/${selectedGroup.id}/members/${id}`);
    await refreshMembers();
  };

  const exportMembers = () => {
    const blob = new Blob([JSON.stringify(members, null, 2)], { type: "application/json" });
    saveAs(blob, `lootforge-members-${selectedGroup?.name || "export"}.json`);
  };

  return { members, loading, error, addMember, updateMember, deleteMember, refreshMembers, exportMembers };
}
