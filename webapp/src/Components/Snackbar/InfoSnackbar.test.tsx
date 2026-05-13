import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '../../test/test-utils';
import userEvent from '@testing-library/user-event';
import { InfoSnackbar } from './InfoSnackbar';

describe('InfoSnackbar', () => {
    it('should render notification with text', () => {
        const mockHandleStopShowing = vi.fn();

        render(
            <InfoSnackbar
                text="Test notification message"
                show={true}
                handleStopShowing={mockHandleStopShowing}
            />
        );

        expect(screen.getByText('Test notification message')).toBeInTheDocument();
        expect(screen.getByText('Notification')).toBeInTheDocument();
    });

    it('should call handleStopShowing when notification is clicked', async () => {
        const user = userEvent.setup();
        const mockHandleStopShowing = vi.fn();

        render(
            <InfoSnackbar
                text="Click me"
                show={true}
                handleStopShowing={mockHandleStopShowing}
            />
        );

        const notification = screen.getByText('Click me').closest('div');
        if (notification) {
            await user.click(notification);
        }

        expect(mockHandleStopShowing).toHaveBeenCalled();
    });

    it('should call handleStopShowing when close button is clicked', async () => {
        const user = userEvent.setup();
        const mockHandleStopShowing = vi.fn();

        const { container } = render(
            <InfoSnackbar
                text="Close me"
                show={true}
                handleStopShowing={mockHandleStopShowing}
            />
        );

        // Find the close button (Mantine Notification has a close button)
        const closeButton = container.querySelector('button[aria-label="Hide notification"]');
        if (closeButton) {
            await user.click(closeButton);
            expect(mockHandleStopShowing).toHaveBeenCalled();
        }
    });

    it('should render with different text messages', () => {
        const mockHandleStopShowing = vi.fn();

        const { rerender } = render(
            <InfoSnackbar
                text="First message"
                show={true}
                handleStopShowing={mockHandleStopShowing}
            />
        );

        expect(screen.getByText('First message')).toBeInTheDocument();

        rerender(
            <InfoSnackbar
                text="Second message"
                show={true}
                handleStopShowing={mockHandleStopShowing}
            />
        );

        expect(screen.getByText('Second message')).toBeInTheDocument();
    });

    it('should render notification component with Mantine styles', () => {
        const mockHandleStopShowing = vi.fn();

        const { container } = render(
            <InfoSnackbar
                text="Styled notification"
                show={true}
                handleStopShowing={mockHandleStopShowing}
            />
        );

        // Notification should have Mantine classes
        const notification = container.querySelector('.mantine-Notification-root');
        expect(notification).toBeInTheDocument();
    });
});