export type StrategyType = "weighted_luck" | "round_robin" | "dkp" | "pure_random";

export interface WeightedLuckConfig {
  increment: number;
  decrement: number;
  defaultLuck: number;
}

export interface DKPConfig {
  winCost: number;
}

export interface StrategyConfig {
  type: StrategyType;
  weightedLuck?: WeightedLuckConfig;
  dkp?: DKPConfig;
}

export interface User {
  id: string;
  email: string;
  displayName: string;
  photoURL: string;
  createdAt: string;
  lastLoginAt: string;
}

export interface Group {
  id: string;
  name: string;
  ownerID: string;
  strategy: StrategyConfig;
  createdAt: string;
}

export interface Member {
  id: string;
  groupID: string;
  name: string;
  role: string;
  luck: number;
  priority: number;
  points: number;
  createdAt: string;
}

export interface RollSession {
  id: string;
  groupID: string;
  strategy: StrategyConfig;
  participantIDs: string[];
  winnerID: string;
  item: string;
  createdAt: string;
}

export interface RollRequest {
  participantIDs: string[];
  item?: string;
}

export interface RollResponse {
  id: string;
  winner: Member;
  participants: Member[];
  item: string;
  strategy: StrategyConfig;
}

export interface Invite {
  token: string;
  groupID: string;
  createdBy: string;
  expiresAt: string;
}
