import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '../../test/test-utils';
import userEvent from '@testing-library/user-event';
import { SettingsView } from './SettingsView';
import { ISettingsAdapter } from '../../Adapters/ISettingsAdapter';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { Setting } from '../../Models/Setting';

// Mock SettingModal
vi.mock('../SettingModal/SettingModal', () => ({
    SettingModal: ({ isOpen, handleSave, handleClose }: any) => (
        isOpen ? (
            <div data-testid="setting-modal">
                <button onClick={() => handleSave('key', 'value', 'name', 'cat', 'desc')}>Save</button>
                <button onClick={handleClose}>Close</button>
            </div>
        ) : null
    )
}));

// Mock InfoSnackbar
vi.mock('../Snackbar/InfoSnackbar', () => ({
    InfoSnackbar: ({ text, show }: any) => (
        show ? <div data-testid="snackbar">{text}</div> : null
    )
}));

describe('SettingsView', () => {
    let mockSettingsAdapter: ISettingsAdapter;
    let mockPhotosAdapter: IPhotosAdapter;
    let mockSettings: Setting[];

    beforeEach(() => {
        mockSettings = [
            { key: 'setting1', value: 'value1', friendlyName: 'Setting One', category: 'Category1', description: 'Description 1' },
            { key: 'setting2', value: 'value2', friendlyName: 'Setting Two', category: 'Category2', description: 'Description 2' },
        ];

        mockSettingsAdapter = {
            getAllSettings: vi.fn().mockResolvedValue(mockSettings),
            updateSettings: vi.fn().mockResolvedValue(undefined),
        } as unknown as ISettingsAdapter;

        mockPhotosAdapter = {
            index: vi.fn().mockResolvedValue(undefined),
        } as unknown as IPhotosAdapter;
    });

    it('should render settings table', async () => {
        render(
            <SettingsView
                settingsAdapter={mockSettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Setting')).toBeInTheDocument();
            expect(screen.getByText('Value')).toBeInTheDocument();
            expect(screen.getByText('Friendly Name')).toBeInTheDocument();
            expect(screen.getByText('Category')).toBeInTheDocument();
            expect(screen.getByText('Description')).toBeInTheDocument();
        });
    });

    it('should fetch and display all settings', async () => {
        render(
            <SettingsView
                settingsAdapter={mockSettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('setting1')).toBeInTheDocument();
            expect(screen.getByText('value1')).toBeInTheDocument();
            expect(screen.getByText('Setting One')).toBeInTheDocument();
        });

        expect(mockSettingsAdapter.getAllSettings).toHaveBeenCalled();
    });

    it('should render Index button', async () => {
        render(
            <SettingsView
                settingsAdapter={mockSettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Index')).toBeInTheDocument();
        });
    });

    it('should call photosAdapter.index when Index button is clicked', async () => {
        const user = userEvent.setup();

        render(
            <SettingsView
                settingsAdapter={mockSettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Index')).toBeInTheDocument();
        });

        const indexButton = screen.getByText('Index');
        await user.click(indexButton);

        expect(mockPhotosAdapter.index).toHaveBeenCalled();
    });

    it('should render edit and add buttons', async () => {
        const { container } = render(
            <SettingsView
                settingsAdapter={mockSettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            const actionIcons = container.querySelectorAll('.mantine-ActionIcon-root');
            expect(actionIcons.length).toBeGreaterThanOrEqual(2);
        });
    });

    it('should not show modal initially', () => {
        render(
            <SettingsView
                settingsAdapter={mockSettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        expect(screen.queryByTestId('setting-modal')).not.toBeInTheDocument();
    });

    it('should not show snackbar initially', () => {
        render(
            <SettingsView
                settingsAdapter={mockSettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        expect(screen.queryByTestId('snackbar')).not.toBeInTheDocument();
    });

    it('should display settings in table rows', async () => {
        render(
            <SettingsView
                settingsAdapter={mockSettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('setting2')).toBeInTheDocument();
            expect(screen.getByText('value2')).toBeInTheDocument();
        });
    });

    it('should render table with proper structure', async () => {
        const { container } = render(
            <SettingsView
                settingsAdapter={mockSettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            const table = container.querySelector('.mantine-Table-table');
            expect(table).toBeInTheDocument();
        });
    });

    it('should handle empty settings list', async () => {
        const emptySettingsAdapter = {
            getAllSettings: vi.fn().mockResolvedValue([]),
            updateSettings: vi.fn().mockResolvedValue(undefined),
        } as unknown as ISettingsAdapter;

        const { container } = render(
            <SettingsView
                settingsAdapter={emptySettingsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            const tbody = container.querySelector('.mantine-Table-tbody');
            expect(tbody).toBeInTheDocument();
        });
    });
});