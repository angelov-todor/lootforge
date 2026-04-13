"use client";
import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from "react";
import { api } from "@/lib/api";
import type { Group } from "@/types";

const STORAGE_KEY = "lootforge_selected_group";

interface GroupContextType {
  groups: Group[];
  selectedGroup: Group | null;
  setSelectedGroup: (group: Group | null) => void;
  refreshGroups: () => Promise<void>;
}

const GroupContext = createContext<GroupContextType>({
  groups: [],
  selectedGroup: null,
  setSelectedGroup: () => {},
  refreshGroups: async () => {},
});

export function GroupProvider({ children }: { children: ReactNode }) {
  const [groups, setGroups] = useState<Group[]>([]);
  const [selectedGroup, setSelectedGroupState] = useState<Group | null>(null);

  const refreshGroups = useCallback(async () => {
    try {
      const data = await api.get<Group[]>("/api/groups");
      const fetched = data || [];
      setGroups(fetched);

      // Restore from localStorage or pick first
      setSelectedGroupState((prev) => {
        if (prev) {
          const match = fetched.find((g) => g.id === prev.id);
          return match ?? fetched[0] ?? null;
        }
        const savedId = typeof window !== "undefined" ? localStorage.getItem(STORAGE_KEY) : null;
        if (savedId) {
          const match = fetched.find((g) => g.id === savedId);
          if (match) return match;
        }
        return fetched[0] ?? null;
      });
    } catch (e) {
      console.error("Failed to fetch groups", e);
    }
  }, []);

  // Initial load: restore from localStorage before fetching
  useEffect(() => {
    refreshGroups();
  }, [refreshGroups]);

  const setSelectedGroup = useCallback((group: Group | null) => {
    setSelectedGroupState(group);
    if (group) {
      localStorage.setItem(STORAGE_KEY, group.id);
    } else {
      localStorage.removeItem(STORAGE_KEY);
    }
  }, []);

  return (
    <GroupContext.Provider value={{ groups, selectedGroup, setSelectedGroup, refreshGroups }}>
      {children}
    </GroupContext.Provider>
  );
}

export function useGroup() {
  return useContext(GroupContext);
}
