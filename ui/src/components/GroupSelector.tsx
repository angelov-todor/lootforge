"use client";
import { Select, MenuItem, FormControl, type SelectChangeEvent } from "@mui/material";
import { useGroup } from "@/components/GroupContext";

export function GroupSelector() {
  const { groups, selectedGroup, setSelectedGroup } = useGroup();

  const handleChange = (e: SelectChangeEvent<string>) => {
    const value = e.target.value;
    if (value === "__create__") return; // placeholder — functional in Phase 6
    const group = groups.find((g) => g.id === value);
    if (group) setSelectedGroup(group);
  };

  return (
    <FormControl size="small" sx={{ minWidth: 200 }}>
      <Select
        value={selectedGroup?.id ?? ""}
        onChange={handleChange}
        displayEmpty
        sx={{ color: "inherit" }}
      >
        {groups.map((g) => (
          <MenuItem key={g.id} value={g.id}>
            {g.name}
          </MenuItem>
        ))}
        <MenuItem value="__create__" divider sx={{ fontStyle: "italic", color: "text.secondary" }}>
          Create New Group
        </MenuItem>
      </Select>
    </FormControl>
  );
}
