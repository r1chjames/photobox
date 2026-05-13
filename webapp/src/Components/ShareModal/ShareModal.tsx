import React, { useState } from 'react';
import { Button, Group, Modal, Select, Text, TextInput, CopyButton, ActionIcon, Tooltip } from '@mantine/core';
import { IconCopy, IconCheck } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { ISharesAdapter } from '../../Adapters/ISharesAdapter';

interface ShareModalProps {
    opened: boolean;
    onClose: () => void;
    resourceType: 'photo' | 'album';
    resourceId: string;
    resourceName: string;
    sharesAdapter: ISharesAdapter;
}

export const ShareModal: React.FC<ShareModalProps> = ({ opened, onClose, resourceType, resourceId, resourceName, sharesAdapter }) => {
    const [expiry, setExpiry] = useState<string | null>('7days');
    const [password, setPassword] = useState('');
    const [shareUrl, setShareUrl] = useState<string | null>(null);
    const [isCreating, setIsCreating] = useState(false);

    const handleCreate = async () => {
        setIsCreating(true);
        try {
            let expiryDate: string | undefined;
            if (expiry === '7days') {
                expiryDate = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString();
            } else if (expiry === '30days') {
                expiryDate = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString();
            }
            const share = await sharesAdapter.createShare(resourceType, resourceId, expiryDate, password || undefined);
            setShareUrl(share.url);
            notifications.show({
                title: 'Share link created',
                message: 'Your share link is ready to copy',
                color: 'green',
            });
        } catch (e) {
            notifications.show({
                title: 'Failed to create share link',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        } finally {
            setIsCreating(false);
        }
    };

    const handleClose = () => {
        setShareUrl(null);
        setPassword('');
        setExpiry('7days');
        onClose();
    };

    return (
        <Modal opened={opened} onClose={handleClose} title={`Share ${resourceType}`} centered>
            <Text size="sm" c="dimmed" mb="md">
                Create a shareable link for &quot;{resourceName}&quot;
            </Text>

            {!shareUrl ? (
                <>
                    <Select
                        label="Link expiry"
                        value={expiry}
                        onChange={setExpiry}
                        data={[
                            { value: 'never', label: 'Never' },
                            { value: '7days', label: '7 days' },
                            { value: '30days', label: '30 days' },
                        ]}
                        mb="sm"
                    />
                    <TextInput
                        label="Password protection (optional)"
                        placeholder="Leave empty for no password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        type="password"
                        mb="md"
                    />
                    <Group justify="flex-end">
                        <Button variant="default" onClick={handleClose}>Cancel</Button>
                        <Button onClick={handleCreate} loading={isCreating}>Create link</Button>
                    </Group>
                </>
            ) : (
                <>
                    <Text size="sm" fw={500} mb="xs">Share link:</Text>
                    <Group gap="xs">
                        <TextInput value={shareUrl} readOnly style={{ flex: 1 }} />
                        <CopyButton value={shareUrl} timeout={2000}>
                            {({ copied, copy }) => (
                                <Tooltip label={copied ? 'Copied' : 'Copy'}>
                                    <ActionIcon color={copied ? 'teal' : 'gray'} variant="light" onClick={copy} size="lg">
                                        {copied ? <IconCheck size="1rem" /> : <IconCopy size="1rem" />}
                                    </ActionIcon>
                                </Tooltip>
                            )}
                        </CopyButton>
                    </Group>
                    <Group justify="flex-end" mt="md">
                        <Button variant="default" onClick={handleClose}>Done</Button>
                        <Button variant="light" onClick={() => setShareUrl(null)}>Create another</Button>
                    </Group>
                </>
            )}
        </Modal>
    );
};
