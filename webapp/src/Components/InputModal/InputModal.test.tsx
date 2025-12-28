import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '../../test/test-utils';
import userEvent from '@testing-library/user-event';
import { InputModal } from './InputModal';

describe('InputModal', () => {
    it('should render modal when isOpen is true', () => {
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        render(
            <InputModal
                isOpen={true}
                title="Test Modal"
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            >
                <div>Modal Content</div>
            </InputModal>
        );

        expect(screen.getByText('Test Modal')).toBeInTheDocument();
        expect(screen.getByText('Modal Content')).toBeInTheDocument();
    });

    it('should not render modal when isOpen is false', () => {
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        render(
            <InputModal
                isOpen={false}
                title="Hidden Modal"
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            >
                <div>Hidden Content</div>
            </InputModal>
        );

        expect(screen.queryByText('Hidden Modal')).not.toBeInTheDocument();
        expect(screen.queryByText('Hidden Content')).not.toBeInTheDocument();
    });

    it('should call handleSave when Save button is clicked', async () => {
        const user = userEvent.setup();
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        render(
            <InputModal
                isOpen={true}
                title="Save Test"
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            >
                <div>Content</div>
            </InputModal>
        );

        const saveButton = screen.getByText('Save');
        await user.click(saveButton);

        expect(mockHandleSave).toHaveBeenCalledTimes(1);
    });

    it('should call handleClose when Close button is clicked', async () => {
        const user = userEvent.setup();
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        render(
            <InputModal
                isOpen={true}
                title="Close Test"
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            >
                <div>Content</div>
            </InputModal>
        );

        const closeButton = screen.getByText('Close');
        await user.click(closeButton);

        expect(mockHandleClose).toHaveBeenCalledTimes(1);
    });

    it('should render children content correctly', () => {
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        render(
            <InputModal
                isOpen={true}
                title="Children Test"
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            >
                <input placeholder="Test Input" />
                <p>Additional content</p>
            </InputModal>
        );

        expect(screen.getByPlaceholderText('Test Input')).toBeInTheDocument();
        expect(screen.getByText('Additional content')).toBeInTheDocument();
    });

    it('should display both Save and Close buttons', () => {
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        render(
            <InputModal
                isOpen={true}
                title="Buttons Test"
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            >
                <div>Content</div>
            </InputModal>
        );

        expect(screen.getByText('Save')).toBeInTheDocument();
        expect(screen.getByText('Close')).toBeInTheDocument();
    });

    it('should render with different titles', () => {
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        const { rerender } = render(
            <InputModal
                isOpen={true}
                title="First Title"
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            >
                <div>Content</div>
            </InputModal>
        );

        expect(screen.getByText('First Title')).toBeInTheDocument();

        rerender(
            <InputModal
                isOpen={true}
                title="Second Title"
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            >
                <div>Content</div>
            </InputModal>
        );

        expect(screen.getByText('Second Title')).toBeInTheDocument();
        expect(screen.queryByText('First Title')).not.toBeInTheDocument();
    });
});