import React, { useEffect, useState, useCallback } from 'react';
import { IUsersAdapter } from '../../Adapters/IUsersAdapter';
import { User } from '../../Models/User';
import { EmptyState } from '../EmptyState/EmptyState';
import { notifications } from '@mantine/notifications';
import { modals } from '@mantine/modals';
import { Table, Title, Select, ActionIcon, Loader, Center, Text, Tooltip, Skeleton } from '@mantine/core';
import { IconTrash, IconUsers } from '@tabler/icons-react';

interface UserManagementProps {
    usersAdapter: IUsersAdapter;
}

export const UserManagement: React.FC<UserManagementProps> = ({ usersAdapter }) => {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);

    const loadUsers = useCallback(async () => {
        try {
            if (typeof usersAdapter.getAllUsers !== 'function') {
                throw new Error('usersAdapter.getAllUsers is not a function');
            }
            const data = await usersAdapter.getAllUsers();
            setUsers(Array.isArray(data) ? data : []);
        } catch (e) {
            notifications.show({
                title: 'Failed to load users',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
            setUsers([]);
        } finally {
            setLoading(false);
        }
    }, [usersAdapter]);

    useEffect(() => {
        loadUsers();
    }, [loadUsers]);

    const handleRoleChange = async (userId: string, newRole: string) => {
        if (!userId) {
            notifications.show({
                title: 'Cannot update user',
                message: 'User ID is missing',
                color: 'red',
            });
            return;
        }
        try {
            await usersAdapter.updateUser(userId, { role: newRole });
            notifications.show({
                title: 'User updated',
                message: `Role changed to ${newRole}`,
                color: 'green',
            });
            setUsers(prev => prev.map(u => u.id === userId ? { ...u, role: newRole } : u));
        } catch (e) {
            notifications.show({
                title: 'Update failed',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    };

    const handleDelete = (user: User) => {
        modals.openConfirmModal({
            title: 'Delete user?',
            children: `Are you sure you want to delete ${user.name || user.email}? This action cannot be undone.`,
            labels: { confirm: 'Delete', cancel: 'Cancel' },
            confirmProps: { color: 'red' },
            onConfirm: async () => {
                try {
                    await usersAdapter.deleteUser(user.id!);
                    notifications.show({
                        title: 'User deleted',
                        message: `${user.name || user.email} has been deleted`,
                        color: 'green',
                    });
                    setUsers(prev => prev.filter(u => u.id !== user.id!));
                } catch (e) {
                    notifications.show({
                        title: 'Delete failed',
                        message: e instanceof Error ? e.message : 'An error occurred',
                        color: 'red',
                    });
                }
            },
        });
    };

    if (loading) {
        return (
            <div>
                <Title size="h4" mb="md">User Management</Title>
                <Table>
                    <Table.Thead>
                        <Table.Tr>
                            <Table.Th>Username</Table.Th>
                            <Table.Th>Email</Table.Th>
                            <Table.Th>Role</Table.Th>
                            <Table.Th>Created</Table.Th>
                            <Table.Th>Actions</Table.Th>
                        </Table.Tr>
                    </Table.Thead>
                    <Table.Tbody>
                        {Array.from({ length: 5 }).map((_, i) => (
                            <Table.Tr key={i}>
                                <Table.Td><Skeleton height={20} width={100} radius="sm" /></Table.Td>
                                <Table.Td><Skeleton height={20} width={160} radius="sm" /></Table.Td>
                                <Table.Td><Skeleton height={28} width={90} radius="sm" /></Table.Td>
                                <Table.Td><Skeleton height={20} width={80} radius="sm" /></Table.Td>
                                <Table.Td><Skeleton height={28} width={28} radius="sm" /></Table.Td>
                            </Table.Tr>
                        ))}
                    </Table.Tbody>
                </Table>
            </div>
        );
    }

    if (!Array.isArray(users) || users.length === 0) {
        return (
            <>
                <Title size="h4" mb="md">User Management</Title>
                <EmptyState
                    title="No users found"
                    description="Registered users will appear here once the backend endpoint is implemented."
                    icon={<IconUsers size="2rem" />}
                />
            </>
        );
    }

    return (
        <div>
            <Title size="h4" mb="md">User Management ({users.length})</Title>
            <Table>
                <Table.Thead>
                    <Table.Tr>
                        <Table.Th>Username</Table.Th>
                        <Table.Th>Email</Table.Th>
                        <Table.Th>Role</Table.Th>
                        <Table.Th>Created</Table.Th>
                        <Table.Th>Actions</Table.Th>
                    </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                    {users.map((user, index) => {
                        const hasId = !!user.id;
                        return (
                            <Table.Tr key={user.id || index}>
                                <Table.Td>{user.name}</Table.Td>
                                <Table.Td>{user.email}</Table.Td>
                                <Table.Td>
                                    {hasId ? (
                                        <Select
                                            size="xs"
                                            value={user.role || 'viewer'}
                                            onChange={(value) => handleRoleChange(user.id!, value!)}
                                            data={[
                                                { value: 'administrator', label: 'Admin' },
                                                { value: 'contributor', label: 'Contributor' },
                                                { value: 'viewer', label: 'Viewer' },
                                            ]}
                                        />
                                    ) : (
                                        <Text size="sm" c="dimmed">{user.role || 'viewer'} (no ID)</Text>
                                    )}
                                </Table.Td>
                                <Table.Td>{user.created_at ? new Date(user.created_at).toLocaleDateString() : '-'}</Table.Td>
                                <Table.Td>
                                    {hasId ? (
                                        <ActionIcon color="red" variant="light" onClick={() => handleDelete(user)} size="sm">
                                            <IconTrash size="0.8rem" />
                                        </ActionIcon>
                                    ) : (
                                        <Tooltip label="Cannot delete user without ID">
                                            <ActionIcon color="red" variant="light" disabled size="sm">
                                                <IconTrash size="0.8rem" />
                                            </ActionIcon>
                                        </Tooltip>
                                    )}
                                </Table.Td>
                            </Table.Tr>
                        );
                    })}
                </Table.Tbody>
            </Table>
        </div>
    );
};
