import React, { useEffect, useState, useCallback } from 'react';
import { IPhotosAdapter, TimelineEntry } from '../../Adapters/IPhotosAdapter';
import { ActionIcon, Badge, Tooltip, ScrollArea, Text } from '@mantine/core';
import { IconClock, IconX } from '@tabler/icons-react';

interface TimelineScrubberProps {
  photosAdapter: IPhotosAdapter;
  onSelectMonth: (year: number, month: number) => void;
  onClear: () => void;
  activeYear?: number;
  activeMonth?: number;
}

export const TimelineScrubber: React.FC<TimelineScrubberProps> = ({
  photosAdapter,
  onSelectMonth,
  onClear,
  activeYear,
  activeMonth,
}) => {
  const [entries, setEntries] = useState<TimelineEntry[]>([]);
  const [loading, setLoading] = useState(true);

  const loadTimeline = useCallback(async () => {
    try {
      const data = await photosAdapter.getTimeline();
      setEntries(data || []);
    } catch (e) {
      setEntries([]);
    } finally {
      setLoading(false);
    }
  }, [photosAdapter]);

  useEffect(() => {
    loadTimeline();
  }, [loadTimeline]);

  if (loading || entries.length === 0) {
    return null;
  }

  // Group by year
  const grouped = entries.reduce<Record<number, TimelineEntry[]>>((acc, entry) => {
    if (!acc[entry.year]) acc[entry.year] = [];
    acc[entry.year].push(entry);
    return acc;
  }, {});

  const years = Object.keys(grouped).map(Number).sort((a, b) => b - a);

  const monthName = (month: number) => {
    return new Date(2000, month - 1, 1).toLocaleString('default', { month: 'short' });
  };

  return (
    <div
      style={{
        position: 'absolute',
        top: 60,
        right: 8,
        width: 120,
        maxHeight: 'calc(100% - 80px)',
        zIndex: 5,
        background: 'var(--mantine-color-body)',
        border: '1px solid var(--mantine-color-default-border)',
        borderRadius: 8,
        padding: '8px 4px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
      }}
    >
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 4, padding: '0 4px' }}>
        <Text size="xs" fw={700} c="dimmed">
          <IconClock size={12} style={{ display: 'inline', verticalAlign: 'middle', marginRight: 4 }} />
          Timeline
        </Text>
        {(activeYear !== undefined && activeMonth !== undefined) && (
          <Tooltip label="Clear filter">
            <ActionIcon size="xs" variant="subtle" onClick={onClear}>
              <IconX size={12} />
            </ActionIcon>
          </Tooltip>
        )}
      </div>
      <ScrollArea style={{ maxHeight: 'calc(100vh - 200px)' }}>
        {years.map(year => (
          <div key={year} style={{ marginBottom: 8 }}>
            <Text size="xs" fw={700} c="var(--mantine-primary-color-filled)" style={{ paddingLeft: 4 }}>
              {year}
            </Text>
            {grouped[year].map(entry => {
              const isActive = activeYear === entry.year && activeMonth === entry.month;
              return (
                <div
                  key={`${entry.year}-${entry.month}`}
                  onClick={() => onSelectMonth(entry.year, entry.month)}
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: '2px 4px',
                    borderRadius: 4,
                    cursor: 'pointer',
                    background: isActive ? 'var(--mantine-primary-color-light)' : 'transparent',
                  }}
                >
                  <Text size="xs" c={isActive ? 'var(--mantine-primary-color-filled)' : undefined}>
                    {monthName(entry.month)}
                  </Text>
                  <Badge size="xs" variant={isActive ? 'filled' : 'light'} color="gray">
                    {entry.count}
                  </Badge>
                </div>
              );
            })}
          </div>
        ))}
      </ScrollArea>
    </div>
  );
};
