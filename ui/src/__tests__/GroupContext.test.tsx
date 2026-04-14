import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor, cleanup } from '@testing-library/react';
import { GroupProvider, useGroup } from '@/components/GroupContext';
import type { Group } from '@/types';

// Mock the api module
vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    del: vi.fn(),
    createInvite: vi.fn(),
    acceptInvite: vi.fn(),
    revokeInvite: vi.fn(),
    createGroup: vi.fn(),
  },
}));

import { api } from '@/lib/api';

const mockGroups: Group[] = [
  {
    id: 'g1',
    name: 'Group One',
    ownerID: 'u1',
    strategy: { type: 'pure_random' },
    createdAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'g2',
    name: 'Group Two',
    ownerID: 'u1',
    strategy: { type: 'dkp', dkp: { winCost: 50 } },
    createdAt: '2024-01-02T00:00:00Z',
  },
];

function TestConsumer() {
  const { groups, selectedGroup } = useGroup();
  return (
    <div>
      <span data-testid="count">{groups.length}</span>
      <span data-testid="selected">{selectedGroup?.name ?? 'none'}</span>
    </div>
  );
}

describe('GroupContext', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  afterEach(() => {
    cleanup();
  });

  it('fetches groups and selects the first one', async () => {
    vi.mocked(api.get).mockResolvedValueOnce(mockGroups);

    render(
      <GroupProvider>
        <TestConsumer />
      </GroupProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId('count')).toHaveTextContent('2');
    });

    expect(screen.getByTestId('selected')).toHaveTextContent('Group One');
  });

  it('handles empty groups', async () => {
    vi.mocked(api.get).mockResolvedValueOnce([]);

    render(
      <GroupProvider>
        <TestConsumer />
      </GroupProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId('count')).toHaveTextContent('0');
    });

    expect(screen.getByTestId('selected')).toHaveTextContent('none');
  });

  it('handles API error gracefully', async () => {
    vi.mocked(api.get).mockRejectedValueOnce(new Error('Network error'));

    render(
      <GroupProvider>
        <TestConsumer />
      </GroupProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId('count')).toHaveTextContent('0');
    });
  });
});
