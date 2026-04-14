import { describe, it, expect } from 'vitest';
import type {
  Group,
  Member,
  RollSession,
  StrategyConfig,
  Invite,
} from '@/types';

describe('Type interfaces', () => {
  it('Group should have required fields', () => {
    const group: Group = {
      id: 'g1',
      name: 'Test Group',
      ownerID: 'u1',
      strategy: { type: 'pure_random' },
      createdAt: '2024-01-01T00:00:00Z',
    };
    expect(group.id).toBe('g1');
    expect(group.strategy.type).toBe('pure_random');
  });

  it('Member should have stat fields', () => {
    const member: Member = {
      id: 'm1',
      groupID: 'g1',
      name: 'Alice',
      role: 'member',
      luck: 50,
      priority: 0,
      points: 100,
      createdAt: '2024-01-01T00:00:00Z',
    };
    expect(member.luck).toBe(50);
    expect(member.points).toBe(100);
  });

  it('RollSession should reference a winner', () => {
    const roll: RollSession = {
      id: 'r1',
      groupID: 'g1',
      strategy: { type: 'dkp', dkp: { winCost: 50 } },
      participantIDs: ['m1', 'm2'],
      winnerID: 'm1',
      item: 'Sword',
      createdAt: '2024-01-01T00:00:00Z',
    };
    expect(roll.winnerID).toBe('m1');
    expect(roll.strategy.dkp?.winCost).toBe(50);
  });

  it('StrategyConfig supports all strategy types', () => {
    const strategies: StrategyConfig[] = [
      { type: 'pure_random' },
      { type: 'weighted_luck', weightedLuck: { increment: 10, decrement: 5, defaultLuck: 50 } },
      { type: 'round_robin' },
      { type: 'dkp', dkp: { winCost: 100 } },
    ];
    expect(strategies).toHaveLength(4);
  });

  it('Invite should have token and expiry', () => {
    const invite: Invite = {
      token: 'abc123',
      groupID: 'g1',
      createdBy: 'u1',
      expiresAt: '2024-12-31T23:59:59Z',
    };
    expect(invite.token).toBe('abc123');
    expect(invite.expiresAt).toContain('2024');
  });
});
