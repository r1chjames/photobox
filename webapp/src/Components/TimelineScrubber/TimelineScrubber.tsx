import React, { useEffect, useState, useCallback, useRef } from 'react';
import { IPhotosAdapter, TimelineEntry } from '../../Adapters/IPhotosAdapter';
import { ActionIcon, Badge, Drawer, Tooltip, Text } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { IconClock, IconX } from '@tabler/icons-react';

interface TimelineScrubberProps {
  photosAdapter: IPhotosAdapter;
  onSelectMonth: (year: number, month: number) => void;
  onClear: () => void;
  activeYear?: number;
  activeMonth?: number;
}

const THROTTLE_MS = 100;
const TRACK_WIDTH = 2;
const CONTAINER_WIDTH = 24;
const HANDLE_SIZE = 12;

export const TimelineScrubber: React.FC<TimelineScrubberProps> = ({
  photosAdapter,
  onSelectMonth,
  onClear,
  activeYear,
  activeMonth,
}) => {
  const [entries, setEntries] = useState<TimelineEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [isHovering, setIsHovering] = useState(false);
  const [scrubberPercent, setScrubberPercent] = useState(0);
  const [hoveredMonth, setHoveredMonth] = useState<{ year: number; month: number } | null>(null);

  const isMobile = useMediaQuery('(max-width: 48em)');
  const trackRef = useRef<HTMLDivElement>(null);
  const throttleTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastSelectRef = useRef<{ year: number; month: number } | null>(null);

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

  // Build flat ordered list of months (newest first, matching entries order)
  const flatMonths = entries; // entries are already ordered year DESC, month DESC

  // Compute position percentage for a given entry index
  const getPercentForIndex = useCallback(
    (index: number): number => {
      if (flatMonths.length <= 1) return 50;
      return (index / (flatMonths.length - 1)) * 100;
    },
    [flatMonths.length]
  );

  // Find the entry index closest to a given percentage
  const getIndexForPercent = useCallback(
    (percent: number): number => {
      if (flatMonths.length === 0) return 0;
      const clamped = Math.max(0, Math.min(100, percent));
      const index = Math.round((clamped / 100) * (flatMonths.length - 1));
      return Math.max(0, Math.min(flatMonths.length - 1, index));
    },
    [flatMonths.length]
  );

  // Initialize scrubber position to active month if set
  useEffect(() => {
    if (activeYear !== undefined && activeMonth !== undefined) {
      const idx = flatMonths.findIndex(
        (e) => e.year === activeYear && e.month === activeMonth
      );
      if (idx >= 0) {
        setScrubberPercent(getPercentForIndex(idx));
      }
    }
  }, [activeYear, activeMonth, flatMonths, getPercentForIndex]);

  // Throttled select
  const throttledSelect = useCallback(
    (year: number, month: number) => {
      if (
        lastSelectRef.current?.year === year &&
        lastSelectRef.current?.month === month
      ) {
        return;
      }
      lastSelectRef.current = { year, month };

      if (throttleTimerRef.current) {
        clearTimeout(throttleTimerRef.current);
      }
      throttleTimerRef.current = setTimeout(() => {
        onSelectMonth(year, month);
      }, THROTTLE_MS);
    },
    [onSelectMonth]
  );

  // Compute percent from mouse/touch Y position relative to track
  const computePercentFromEvent = useCallback(
    (clientY: number): number => {
      if (!trackRef.current) return 0;
      const rect = trackRef.current.getBoundingClientRect();
      const y = clientY - rect.top;
      const percent = (y / rect.height) * 100;
      return Math.max(0, Math.min(100, percent));
    },
    []
  );

  // Handle drag start
  const handleDragStart = useCallback(
    (clientY: number) => {
      setIsDragging(true);
      const percent = computePercentFromEvent(clientY);
      setScrubberPercent(percent);
      const idx = getIndexForPercent(percent);
      const entry = flatMonths[idx];
      if (entry) {
        throttledSelect(entry.year, entry.month);
      }
    },
    [computePercentFromEvent, getIndexForPercent, flatMonths, throttledSelect]
  );

  // Handle drag move
  const handleDragMove = useCallback(
    (clientY: number) => {
      if (!isDragging) return;
      const percent = computePercentFromEvent(clientY);
      setScrubberPercent(percent);
      const idx = getIndexForPercent(percent);
      const entry = flatMonths[idx];
      if (entry) {
        setHoveredMonth({ year: entry.year, month: entry.month });
        throttledSelect(entry.year, entry.month);
      }
    },
    [isDragging, computePercentFromEvent, getIndexForPercent, flatMonths, throttledSelect]
  );

  // Handle drag end
  const handleDragEnd = useCallback(() => {
    setIsDragging(false);
    if (throttleTimerRef.current) {
      clearTimeout(throttleTimerRef.current);
      // Fire immediately on drag end if there was a pending call
      if (lastSelectRef.current) {
        onSelectMonth(lastSelectRef.current.year, lastSelectRef.current.month);
      }
    }
  }, [onSelectMonth]);

  // Global mouse/touch listeners for dragging
  useEffect(() => {
    if (!isDragging) return;

    const onMouseMove = (e: MouseEvent) => {
      e.preventDefault();
      handleDragMove(e.clientY);
    };
    const onTouchMove = (e: TouchEvent) => {
      e.preventDefault();
      handleDragMove(e.touches[0].clientY);
    };
    const onMouseUp = () => handleDragEnd();
    const onTouchEnd = () => handleDragEnd();

    document.addEventListener('mousemove', onMouseMove);
    document.addEventListener('touchmove', onTouchMove, { passive: false });
    document.addEventListener('mouseup', onMouseUp);
    document.addEventListener('touchend', onTouchEnd);

    return () => {
      document.removeEventListener('mousemove', onMouseMove);
      document.removeEventListener('touchmove', onTouchMove);
      document.removeEventListener('mouseup', onMouseUp);
      document.removeEventListener('touchend', onTouchEnd);
    };
  }, [isDragging, handleDragMove, handleDragEnd]);

  if (loading || entries.length === 0) {
    return null;
  }

  const monthName = (month: number) => {
    return new Date(2000, month - 1, 1).toLocaleString('default', { month: 'short' });
  };

  const handleSelectMonth = (year: number, month: number) => {
    onSelectMonth(year, month);
    if (isMobile) {
      setDrawerOpen(false);
    }
  };

  // Determine displayed month (hovered during drag, otherwise active)
  const displayMonth = hoveredMonth ?? (activeYear !== undefined && activeMonth !== undefined
    ? { year: activeYear, month: activeMonth }
    : null);

  // Build the track content: month labels positioned along the track
  // Only visible during drag or hover
  const renderTrackLabels = () => {
    if (!isDragging && !isHovering) return null;

    return entries.map((entry, idx) => {
      const percent = getPercentForIndex(idx);
      const isActive = activeYear === entry.year && activeMonth === entry.month;
      const isHovered = hoveredMonth?.year === entry.year && hoveredMonth?.month === entry.month;

      return (
        <div
          key={`${entry.year}-${entry.month}`}
          onClick={(e) => {
            e.stopPropagation();
            if (!isDragging) {
              handleSelectMonth(entry.year, entry.month);
            }
          }}
          style={{
            position: 'absolute',
            top: `${percent}%`,
            right: CONTAINER_WIDTH + 4,
            transform: 'translateY(-50%)',
            display: 'flex',
            alignItems: 'center',
            gap: 4,
            cursor: isDragging ? 'default' : 'pointer',
            padding: '2px 6px',
            borderRadius: 4,
            background: isActive || isHovered ? 'var(--mantine-primary-color-light)' : 'var(--mantine-color-body)',
            boxShadow: '0 1px 4px rgba(0,0,0,0.15)',
            opacity: isDragging ? (isHovered ? 1 : 0.5) : (isHovered ? 1 : 0.6),
            transition: isDragging ? 'none' : 'opacity 0.15s, background 0.15s',
            whiteSpace: 'nowrap',
            pointerEvents: isDragging ? 'none' : 'auto',
          }}
        >
          <Text
            size="xs"
            fw={isActive || isHovered ? 700 : 500}
            c={isActive || isHovered ? 'var(--mantine-primary-color-filled)' : 'dimmed'}
          >
            {monthName(entry.month)} {entry.year}
          </Text>
          <Badge size="xs" variant={isActive ? 'filled' : 'light'} color="gray">
            {entry.count}
          </Badge>
        </div>
      );
    });
  };

  const scrubberContent = (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      {/* Header - only visible during interaction on desktop */}
      {(isDragging || isHovering) && (
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 4, padding: '0 4px', flexShrink: 0 }}>
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
      )}

      {/* Display current month during drag/hover - floating tag */}
      {displayMonth && (isDragging || isHovering) && (
        <div style={{ textAlign: 'center', padding: '4px 0', flexShrink: 0 }}>
          <Text size="sm" fw={700} c="var(--mantine-primary-color-filled)">
            {monthName(displayMonth.month)} {displayMonth.year}
          </Text>
        </div>
      )}

      {/* Track area */}
      <div
        style={{
          position: 'relative',
          flex: 1,
          minHeight: 200,
          cursor: isDragging ? 'grabbing' : 'pointer',
        }}
        ref={trackRef}
        onMouseDown={(e) => {
          e.preventDefault();
          handleDragStart(e.clientY);
        }}
        onTouchStart={(e) => {
          handleDragStart(e.touches[0].clientY);
        }}
        onMouseEnter={() => setIsHovering(true)}
        onMouseLeave={() => {
          setIsHovering(false);
          setHoveredMonth(null);
        }}
      >
        {/* Track line - slim dotted line */}
        <div
          style={{
            position: 'absolute',
            right: Math.floor(CONTAINER_WIDTH / 2),
            top: 0,
            bottom: 0,
            width: TRACK_WIDTH,
            borderRight: `${TRACK_WIDTH}px dotted var(--mantine-color-default-border)`,
            opacity: isDragging || isHovering ? 0.8 : 0.4,
            transition: 'opacity 0.15s',
          }}
        />

        {/* Month labels */}
        {renderTrackLabels()}

        {/* Draggable handle - small circle, prominent only during interaction */}
        <div
          style={{
            position: 'absolute',
            right: Math.floor(CONTAINER_WIDTH / 2) - Math.floor(HANDLE_SIZE / 2),
            top: `${scrubberPercent}%`,
            transform: 'translateY(-50%)',
            width: HANDLE_SIZE,
            height: HANDLE_SIZE,
            background: isDragging ? 'var(--mantine-primary-color-filled)' : (isHovering ? 'var(--mantine-primary-color-filled-hover)' : 'var(--mantine-color-default-border)'),
            borderRadius: '50%',
            cursor: 'grab',
            boxShadow: isDragging || isHovering ? '0 1px 4px rgba(0,0,0,0.3)' : 'none',
            opacity: isDragging || isHovering ? 1 : 0.5,
            transition: isDragging ? 'none' : 'top 0.2s ease-out, opacity 0.15s, background 0.15s',
            zIndex: 10,
          }}
          onMouseDown={(e) => {
            e.stopPropagation();
            e.preventDefault();
            handleDragStart(e.clientY);
          }}
          onTouchStart={(e) => {
            e.stopPropagation();
            handleDragStart(e.touches[0].clientY);
          }}
        />
      </div>
    </div>
  );

  // Mobile: FAB + Drawer
  if (isMobile) {
    return (
      <>
        <ActionIcon
          size="lg"
          radius="xl"
          variant="filled"
          color="var(--mantine-primary-color-filled)"
          style={{
            position: 'fixed',
            bottom: 20,
            right: 20,
            zIndex: 100,
            boxShadow: '0 2px 8px rgba(0,0,0,0.2)',
          }}
          onClick={() => setDrawerOpen(true)}
        >
          <IconClock size={20} />
        </ActionIcon>
        <Drawer
          opened={drawerOpen}
          onClose={() => setDrawerOpen(false)}
          position="bottom"
          size="100%"
          radius="lg"
          padding="md"
          title={
            <Text size="sm" fw={700}>
              <IconClock size={16} style={{ display: 'inline', verticalAlign: 'middle', marginRight: 6 }} />
              Timeline
            </Text>
          }
        >
          <div style={{ height: 'calc(100vh - 200px)' }}>
            {scrubberContent}
          </div>
        </Drawer>
      </>
    );
  }

  // Desktop: slim permanent line at the far right edge
  return (
    <div
      style={{
        position: 'absolute',
        top: 60,
        right: 0,
        width: CONTAINER_WIDTH,
        maxHeight: 'calc(100% - 80px)',
        height: 'calc(100% - 80px)',
        zIndex: 5,
        background: isDragging || isHovering ? 'var(--mantine-color-body)' : 'transparent',
        border: isDragging || isHovering ? '1px solid var(--mantine-color-default-border)' : 'none',
        borderRadius: isDragging || isHovering ? 8 : 0,
        padding: isDragging || isHovering ? '8px' : 0,
        boxShadow: isDragging || isHovering ? '0 2px 8px rgba(0,0,0,0.1)' : 'none',
        transition: 'background 0.15s, border 0.15s, padding 0.15s, box-shadow 0.15s',
      }}
    >
      {scrubberContent}
    </div>
  );
};
