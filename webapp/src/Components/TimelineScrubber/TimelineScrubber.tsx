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
const HANDLE_SIZE = 12;
// Desktop overlay rail (issue #174): absolute over the grid's right gutter.
// RAIL_WIDTH must match PhotoGrid's TIMELINE_RAIL_INSET so the rail covers
// only reserved padding and never photo content. RAIL_LINE_CENTER is the
// line's x position measured from the strip's right edge.
const RAIL_WIDTH = 48;
const RAIL_LINE_CENTER = 24;

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

  const isMobile = useMediaQuery('(max-width: 50em)');
  const trackRef = useRef<HTMLDivElement>(null);
  const throttleTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastSelectRef = useRef<{ year: number; month: number } | null>(null);

  const loadTimeline = useCallback(async () => {
    try {
      const data = await photosAdapter.getTimeline();
      setEntries(data || []);
    } catch {
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

  // Keyboard scrubbing (issue #174): the rail is a vertical slider. Entries
  // are newest-first with index 0 at the top, so ArrowUp moves toward newer
  // months and Home/End jump to the extremes. Declared before the early
  // returns below so the hook count stays stable across renders (React #310).
  const handleRailKeyDown = useCallback((e: React.KeyboardEvent) => {
    const count = flatMonths.length;
    if (count === 0) return;
    const currentIdx = getIndexForPercent(scrubberPercent);
    let nextIdx = currentIdx;
    switch (e.key) {
      case 'ArrowUp': nextIdx = Math.max(0, currentIdx - 1); break;
      case 'ArrowDown': nextIdx = Math.min(count - 1, currentIdx + 1); break;
      case 'Home': nextIdx = 0; break;
      case 'End': nextIdx = count - 1; break;
      default: return;
    }
    e.preventDefault();
    setScrubberPercent(getPercentForIndex(nextIdx));
    const entry = flatMonths[nextIdx];
    if (entry) {
      lastSelectRef.current = null; // keyboard steps must not be deduped away
      onSelectMonth(entry.year, entry.month);
    }
  }, [flatMonths, getIndexForPercent, getPercentForIndex, scrubberPercent, onSelectMonth]);

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
  const renderTrackLabels = (labelRight: string) => {
    if (!isDragging && !isHovering) return null;

    // Limit to ~12 labels to avoid crowding
    const maxLabels = 12;
    const skip = Math.max(1, Math.ceil(entries.length / maxLabels));

    return entries.map((entry, idx) => {
      if (idx % skip !== 0 && idx !== entries.length - 1) return null;
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
            right: labelRight,
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
            left: '50%',
            transform: 'translateX(-50%)',
            top: 0,
            bottom: 0,
            width: TRACK_WIDTH,
            borderRight: `${TRACK_WIDTH}px dotted var(--mantine-color-default-border)`,
            opacity: isDragging || isHovering ? 0.8 : 0.4,
            transition: 'opacity 0.15s',
          }}
        />

        {/* Month labels */}
        {renderTrackLabels('calc(50% + 8px)')}

        {/* Draggable handle - small circle, prominent only during interaction */}
        <div
          style={{
            position: 'absolute',
            left: '50%',
            top: `${scrubberPercent}%`,
            transform: 'translate(-50%, -50%)',
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

  const isInteracting = isDragging || isHovering;
  const railTrackIndex = getIndexForPercent(scrubberPercent);
  const railEntry = entries[railTrackIndex];

  // Desktop: overlaid right-edge rail (issue #174). The wrapper is
  // pointer-events:none so photos underneath stay interactive; only the
  // narrow strip itself captures input. Month labels float left over the
  // grid with a translucent chip background.
  return (
    <div
      className="timeline-rail"
      style={{
        position: 'absolute',
        top: 8,
        bottom: 8,
        right: 0,
        width: RAIL_WIDTH,
        pointerEvents: 'none',
        zIndex: 20,
      }}
    >
      <div
        role="slider"
        aria-label="Timeline scrubber"
        aria-valuemin={0}
        aria-valuemax={Math.max(0, entries.length - 1)}
        aria-valuenow={Math.max(0, entries.length - 1 - railTrackIndex)}
        aria-valuetext={railEntry ? `${monthName(railEntry.month)} ${railEntry.year}` : undefined}
        tabIndex={0}
        onKeyDown={handleRailKeyDown}
        ref={trackRef}
        style={{
          position: 'relative',
          height: '100%',
          width: '100%',
          cursor: isDragging ? 'grabbing' : 'pointer',
          pointerEvents: 'auto',
          borderRadius: 8,
        }}
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
            right: RAIL_LINE_CENTER - TRACK_WIDTH / 2,
            top: 0,
            bottom: 0,
            width: TRACK_WIDTH,
            borderRight: `${TRACK_WIDTH}px dotted var(--mantine-color-default-border)`,
            opacity: isInteracting ? 0.8 : 0.4,
            transition: 'opacity 0.15s',
          }}
        />

        {/* Month labels - float left of the line over the photos */}
        {renderTrackLabels(`${RAIL_LINE_CENTER + 8}px`)}

        {/* Draggable handle */}
        <div
          style={{
            position: 'absolute',
            right: RAIL_LINE_CENTER - HANDLE_SIZE / 2,
            top: `${scrubberPercent}%`,
            transform: 'translateY(-50%)',
            width: HANDLE_SIZE,
            height: HANDLE_SIZE,
            background: isDragging ? 'var(--mantine-primary-color-filled)' : (isHovering ? 'var(--mantine-primary-color-filled-hover)' : 'var(--mantine-color-default-border)'),
            borderRadius: '50%',
            cursor: 'grab',
            boxShadow: isInteracting ? '0 1px 4px rgba(0,0,0,0.3)' : 'none',
            opacity: isInteracting ? 1 : 0.6,
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

      {/* Floating month chip - visible while scrubbing or when a filter is active.
          Floats over the grid (left of the rail) with a translucent background;
          pointer-events auto only on this element so photos stay clickable. */}
      {displayMonth && (isInteracting || (activeYear !== undefined && activeMonth !== undefined)) && (
        <div
          style={{
            position: 'absolute',
            top: `${scrubberPercent}%`,
            right: RAIL_WIDTH + 8,
            transform: 'translateY(-50%)',
            display: 'flex',
            alignItems: 'center',
            gap: 4,
            padding: '2px 6px',
            background: 'rgba(18, 19, 22, 0.85)',
            borderRadius: 10,
            whiteSpace: 'nowrap',
            pointerEvents: 'auto',
            zIndex: 10,
          }}
        >
          <Text size="xs" fw={700} c="white">
            {monthName(displayMonth.month)} {displayMonth.year}
          </Text>
          {(activeYear !== undefined && activeMonth !== undefined) && (
            <ActionIcon
              size="xs"
              variant="subtle"
              color="white"
              onClick={(e) => {
                e.stopPropagation();
                onClear();
              }}
            >
              <IconX size={10} />
            </ActionIcon>
          )}
        </div>
      )}
    </div>
  );
};
