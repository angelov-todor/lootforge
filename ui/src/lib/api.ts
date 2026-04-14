import { getIdToken } from "./firebase";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "";

async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = await getIdToken();
  if (!token) {
    throw new Error("Not authenticated");
  }

  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...options.headers,
    },
  });

  if (res.status === 401) {
    window.location.href = "/";
    throw new Error("Unauthorized");
  }

  if (!res.ok) {
    const error = await res.json().catch(() => ({ message: "Request failed" }));
    throw new Error(error.message || `HTTP ${res.status}`);
  }

  if (res.status === 204) return undefined as T;
  return res.json();
}

export const api = {
  get: <T>(path: string) => apiFetch<T>(path),
  post: <T>(path: string, body: unknown) =>
    apiFetch<T>(path, { method: "POST", body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown) =>
    apiFetch<T>(path, { method: "PUT", body: JSON.stringify(body) }),
  patch: <T>(path: string, body: unknown) =>
    apiFetch<T>(path, { method: "PATCH", body: JSON.stringify(body) }),
  del: (path: string) => apiFetch<void>(path, { method: "DELETE" }),

  // Invite helpers
  createInvite: (groupID: string) =>
    apiFetch<import("@/types").Invite>(`/api/groups/${groupID}/invites`, {
      method: "POST",
      body: "{}",
    }),
  acceptInvite: (token: string) =>
    apiFetch<import("@/types").Group>(`/api/invites/${token}/accept`, {
      method: "POST",
      body: "{}",
    }),
  revokeInvite: (token: string) =>
    apiFetch<void>(`/api/invites/${token}`, { method: "DELETE" }),

  // Group helpers
  createGroup: (name: string, strategy: import("@/types").StrategyConfig) =>
    apiFetch<import("@/types").Group>("/api/groups", {
      method: "POST",
      body: JSON.stringify({ name, strategy }),
    }),
};
