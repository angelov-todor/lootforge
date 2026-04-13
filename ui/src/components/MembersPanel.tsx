"use client";
import { useMemo, useState } from "react";
import { Box, Button, IconButton, Tooltip } from "@mui/material";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import DownloadIcon from "@mui/icons-material/Download";
import {
  MaterialReactTable,
  useMaterialReactTable,
  type MRT_ColumnDef,
  type MRT_Row,
} from "material-react-table";
import { useGroup } from "@/components/GroupContext";
import { useMembers } from "@/hooks/useMembers";
import type { Member } from "@/types";

export function MembersPanel() {
  const { selectedGroup } = useGroup();
  const { members, loading, addMember, updateMember, deleteMember, exportMembers } = useMembers();
  const [validationErrors, setValidationErrors] = useState<Record<string, string | undefined>>({});

  const strategyType = selectedGroup?.strategy.type;

  const strategyColumn = useMemo((): MRT_ColumnDef<Member> | null => {
    if (strategyType === "weighted_luck") {
      return {
        accessorKey: "luck",
        header: "Luck",
        muiEditTextFieldProps: { type: "number" },
      };
    }
    if (strategyType === "round_robin") {
      return {
        accessorKey: "priority",
        header: "Priority",
        muiEditTextFieldProps: { type: "number" },
      };
    }
    if (strategyType === "dkp") {
      return {
        accessorKey: "points",
        header: "Points",
        muiEditTextFieldProps: { type: "number" },
      };
    }
    return null;
  }, [strategyType]);

  const columns = useMemo<MRT_ColumnDef<Member>[]>(() => {
    const base: MRT_ColumnDef<Member>[] = [
      {
        accessorKey: "name",
        header: "Name",
        muiEditTextFieldProps: {
          required: true,
          error: !!validationErrors.name,
          helperText: validationErrors.name,
          onFocus: () => setValidationErrors((prev) => ({ ...prev, name: undefined })),
        },
      },
      {
        accessorKey: "role",
        header: "Role",
        muiEditTextFieldProps: {
          onFocus: () => setValidationErrors((prev) => ({ ...prev, role: undefined })),
        },
      },
    ];
    if (strategyColumn) base.push(strategyColumn);
    return base;
  }, [strategyColumn, validationErrors]);

  const handleCreate = async ({ values, table }: { values: Partial<Member>; row?: MRT_Row<Member>; table: ReturnType<typeof useMaterialReactTable<Member>> }) => {
    if (!values.name) {
      setValidationErrors({ name: "Name is required" });
      return;
    }
    setValidationErrors({});
    await addMember(values);
    table.setCreatingRow(null);
  };

  const handleSave = async ({ values, row, table }: { values: Member; row: MRT_Row<Member>; table: ReturnType<typeof useMaterialReactTable<Member>> }) => {
    await updateMember({ ...row.original, ...values });
    table.setEditingRow(null);
  };

  const table = useMaterialReactTable<Member>({
    columns,
    data: members,
    state: { isLoading: loading },
    createDisplayMode: "row",
    editDisplayMode: "row",
    enableEditing: true,
    getRowId: (row) => row.id,
    onCreatingRowSave: handleCreate,
    onEditingRowSave: handleSave,
    onCreatingRowCancel: () => setValidationErrors({}),
    onEditingRowCancel: () => setValidationErrors({}),
    renderRowActions: ({ row, table }) => (
      <Box sx={{ display: "flex", gap: "0.5rem" }}>
        <Tooltip title="Edit">
          <IconButton onClick={() => table.setEditingRow(row)}>
            <EditIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete">
          <IconButton
            color="error"
            onClick={async () => {
              if (confirm(`Delete "${row.original.name}"?`)) {
                await deleteMember(row.original.id);
              }
            }}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </Box>
    ),
    renderTopToolbarCustomActions: ({ table }) => (
      <Box sx={{ display: "flex", gap: 1 }}>
        <Button
          variant="contained"
          size="small"
          onClick={() => table.setCreatingRow(true)}
        >
          Add Member
        </Button>
        <Button
          variant="outlined"
          size="small"
          startIcon={<DownloadIcon />}
          onClick={exportMembers}
        >
          Export
        </Button>
      </Box>
    ),
  });

  return <MaterialReactTable table={table} />;
}
