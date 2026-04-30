import React, { useEffect, useState, useCallback } from 'react';
import { ISharesAdapter, ShareLink } from '../../Adapters/ISharesAdapter';
import { EmptyState } from '../EmptyState/EmptyState';
import { notifications } from '@mantine/notifications';
import { modals } from '@mantine/modals';
import { Group, Table, Title, Text, Badge, ActionIcon, Tooltip, CopyButton } from '@mantine/core';
import { IconCopy, IconCheck, IconLink, IconTrash } from '@tabler/icons-react';

interface ShareManagementProps {
    sharesAdapter: ISharesAdapter;
}

export const ShareManagement: React.FC<ShareManagementProps> = ({ sharesAdapter }) => {
    const [shares, setShares] = useState<ShareLink[]>([]);
    const [loading, setLoading] = useState(true);

    const loadShares = useCallback(async () => {
        try {
            const data = await sharesAdapter.getShares();
            setShares(data);
        } catch (e) {
            notifications.show({
                title: 'Failed to load shares',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        } finally {
            setLoading(false);
        }
    }, [sharesAdapter]);

    useEffect(() => {
        loadShares();
    }, [loadShares]);

    const handleRevoke = (share: ShareLink) => {
        modals.openConfirmModal({
            title: 'Revoke share link?',
            children: 'Anyone with this link will no longer be able to access the shared content.',
            labels: { confirm: 'Revoke', cancel: 'Cancel' },
            confirmProps: { color: 'red' },
            onConfirm: async () => {
                try {
                    await sharesAdapter.revokeShare(share.token);
                    notifications.show({
                        title: 'Link revoked',
                        message: 'The share link has been revoked',
                        color: 'green',
                    });
                    setShares(prev => prev.filter(s => s.token !== share.token));
                } catch (e) {
                    notifications.show({
                        title: 'Failed to revoke',
                        message: e instanceof Error ? e.message : 'An error occurred',
                        color: 'red',
                    });
                }
            },
        });
    };

    if (loading) {
        return <Title size="h4">Share Management</Title>;
    }

    if (shares.length === 0) {
        return (
            <>
                <Title size="h4" mb="md">Share Management</Title>
                <EmptyState
                    title="No active shares"
                    description="Share links you create will appear here. Use the share button on any photo or album to get started."
                    icon={<IconLink size="2rem" />}
                />
            </>
        );
    }

    return (
        <div>
            <Title size="h4" mb="md">Share Management ({shares.length})</Title>
            <Table>
                <Table.Thead>
                    <Table.Tr>
                        <Table.Th>Type</Table.Th>
                        <Table.Th>Resource ID</Table.Th>
                        <Table.Th>Link</Table.Th>
                        <Table.Th>Expiry</Table.Th>
                        <Table.Th>Actions</Table.Th>
                    </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                    {shares.map(share => (
                        <Table.Tr key={share.token}>
                            <Table.Td>
                                <Badge variant="light" color={share.resourceType === 'photo' ? 'blue' : 'teal'}>
                                    {share.resourceType}
                                </Badge>
                            </Table.Td>
                            <Table.Td>
                                <Text size="sm" lineClamp={1} maw={200}>{share.resourceId}</Text>
                            </Table.Td>
                            <Table.Td>
                                <Group gap="xs">
                                    <Text size="sm" lineClamp={1} maw={250}>{share.url}</Text>
                                    <CopyButton value={share.url} timeout={2000}>
                                        {({ copied, copy }) => (
                                            <Tooltip label={copied ? 'Copied' : 'Copy'}>
                                                <ActionIcon color={copied ? 'teal' : 'gray'} variant="light" onClick={copy} size="sm">
                                                    {copied ? <IconCheck size="0.8rem" /> : <IconCopy size="0.8rem" />}
                                                </ActionIcon>
                                            </Tooltip>
                                        )}
                                    </CopyButton>
                                </Group>
                            </Table.Td>
                            <Table.Td>
                                {share.expiry ? new Date(share.expiry).toLocaleDateString() : <Text size="sm" c="dimmed">Never</Text>}
                            </Table.Td>
                            <Table.Td>
                                <ActionIcon color="red" variant="light" onClick={() => handleRevoke(share)} size="sm">
                                    <IconTrash size="0.8rem" />
                                </ActionIcon>
                            </Table.Td>
                        </Table.Tr>
                    ))}
                </Table.Tbody>
            </Table>
        </div>
    );
};
