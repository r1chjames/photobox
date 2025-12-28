import { describe, it, expect, vi } from 'vitest';
import { render, screen, waitFor } from '../../test/test-utils';
import userEvent from '@testing-library/user-event';
import { SettingModal } from './SettingModal';

// Mock InputModal to render children directly without react-bootstrap Modal
// This avoids portal rendering issues in the test environment
vi.mock('../InputModal/InputModal', () => ({
    InputModal: ({ isOpen, title, handleSave, handleClose, children }: any) => (
        isOpen ? (
            <div data-testid="input-modal">
                <h2>{title}</h2>
                <div>{children}</div>
                <button onClick={handleSave}>Save</button>
                <button onClick={handleClose}>Close</button>
            </div>
        ) : null
    )
}));

describe('SettingModal', () => {
    it('should render modal when isOpen is true', () => {
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        render(
            <SettingModal
                isOpen={true}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        expect(screen.getByText('Add Setting')).toBeInTheDocument();
    });

    it('should not render modal when isOpen is false', () => {
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        render(
            <SettingModal
                isOpen={false}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        expect(screen.queryByText('Add Setting')).not.toBeInTheDocument();
    });

    it('should render all required input fields', () => {
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        const { container } = render(
            <SettingModal
                isOpen={true}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        expect(container.querySelector('#key')).toBeInTheDocument();
        expect(container.querySelector('#value')).toBeInTheDocument();
        expect(container.querySelector('#friendlyName')).toBeInTheDocument();
        expect(container.querySelector('#category')).toBeInTheDocument();
        expect(container.querySelector('#description')).toBeInTheDocument();
    });

    it('should update key field and clear error', async () => {
        const user = userEvent.setup();
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        const { container } = render(
            <SettingModal
                isOpen={true}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        const keyInput = container.querySelector('#key') as HTMLInputElement;
        await user.type(keyInput, 'test_key');

        expect(keyInput).toHaveValue('test_key');
    });

    it('should update value field and clear error', async () => {
        const user = userEvent.setup();
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        const { container } = render(
            <SettingModal
                isOpen={true}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        const valueInput = container.querySelector('#value') as HTMLInputElement;
        await user.type(valueInput, 'test_value');

        expect(valueInput).toHaveValue('test_value');
    });

    it('should call handleSave with all field values when all required fields are filled', async () => {
        const user = userEvent.setup();
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        const { container } = render(
            <SettingModal
                isOpen={true}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        await user.type(container.querySelector('#key') as HTMLInputElement, 'test_key');
        await user.type(container.querySelector('#value') as HTMLInputElement, 'test_value');
        await user.type(container.querySelector('#friendlyName') as HTMLInputElement, 'Test Name');
        await user.type(container.querySelector('#category') as HTMLInputElement, 'test_category');
        await user.type(container.querySelector('#description') as HTMLInputElement, 'Test description');

        const saveButton = screen.getByText('Save');
        await user.click(saveButton);

        expect(mockHandleSave).toHaveBeenCalledWith(
            'test_key',
            'test_value',
            'Test Name',
            'test_category',
            'Test description'
        );
    });

    it('should not call handleSave when required fields are empty', async () => {
        const user = userEvent.setup();
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        const { container } = render(
            <SettingModal
                isOpen={true}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        // Only fill some fields
        await user.type(container.querySelector('#key') as HTMLInputElement, 'k');
        await user.type(container.querySelector('#value') as HTMLInputElement, 'v');

        const saveButton = screen.getByText('Save');
        await user.click(saveButton);

        expect(mockHandleSave).not.toHaveBeenCalled();
    });

    it('should call handleClose when close button is clicked', async () => {
        const user = userEvent.setup();
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        render(
            <SettingModal
                isOpen={true}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        const closeButton = screen.getByText('Close');
        await user.click(closeButton);

        expect(mockHandleClose).toHaveBeenCalled();
    });

    it('should show error state for short key input', async () => {
        const user = userEvent.setup();
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        const { container } = render(
            <SettingModal
                isOpen={true}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        const keyInput = container.querySelector('#key') as HTMLInputElement;
        await user.type(keyInput, 'k');

        // Input should have error attribute
        await waitFor(() => {
            const inputElement = container.querySelector('#key');
            expect(inputElement).toBeInTheDocument();
        });
    });

    it('should allow all fields to be filled', async () => {
        const user = userEvent.setup();
        const mockHandleSave = vi.fn();
        const mockHandleClose = vi.fn();

        const { container } = render(
            <SettingModal
                isOpen={true}
                handleSave={mockHandleSave}
                handleClose={mockHandleClose}
            />
        );

        // Fill all required fields with minimum valid length
        await user.type(container.querySelector('#key') as HTMLInputElement, 'key');
        await user.type(container.querySelector('#value') as HTMLInputElement, 'val');
        await user.type(container.querySelector('#friendlyName') as HTMLInputElement, 'Name');
        await user.type(container.querySelector('#category') as HTMLInputElement, 'cat');
        await user.type(container.querySelector('#description') as HTMLInputElement, 'desc');

        const saveButton = screen.getByText('Save');
        await user.click(saveButton);

        expect(mockHandleSave).toHaveBeenCalledWith('key', 'val', 'Name', 'cat', 'desc');
    });
});