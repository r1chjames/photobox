import React, { useState } from 'react';
import { Button, Paper, Stack, TextInput, Title } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { IUsersAdapter } from '../../Adapters/IUsersAdapter';

interface IProps {
    usersAdapter: IUsersAdapter;
}

export const AccountSettings: React.FunctionComponent<IProps> = (props) => {
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [saving, setSaving] = useState(false);

    const handleSave = async () => {
        const updates: Record<string, string> = {};
        if (username) updates.username = username;
        if (email) updates.email = email;
        if (password) updates.password = password;

        if (Object.keys(updates).length === 0) {
            notifications.show({
                title: 'No changes',
                message: 'Please fill in at least one field to update',
                color: 'yellow',
            });
            return;
        }

        setSaving(true);
        try {
            await props.usersAdapter.updateSelf(updates);
            notifications.show({
                title: 'Account updated',
                message: 'Your account settings have been saved successfully',
                color: 'green',
            });
            setPassword('');
        } catch (e) {
            notifications.show({
                title: 'Save failed',
                message: e instanceof Error ? e.message : 'Failed to update account settings',
                color: 'red',
            });
        } finally {
            setSaving(false);
        }
    };

    return (
        <div>
            <Title order={2} mb="md">Account Settings</Title>
            <Paper shadow="sm" p="md" withBorder style={{ maxWidth: 500 }}>
                <Stack gap="md">
                    <TextInput
                        label="Username"
                        placeholder="Your username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                    />
                    <TextInput
                        label="Email"
                        type="email"
                        placeholder="Your email address"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                    />
                    <TextInput
                        label="Password"
                        type="password"
                        placeholder="New password (leave blank to keep current)"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <Button onClick={handleSave} loading={saving} style={{ alignSelf: 'flex-end' }}>
                        Save
                    </Button>
                </Stack>
            </Paper>
        </div>
    );
};
