import React from 'react';
import {Flex, Text, ThemeIcon, Button} from '@mantine/core';
import {IconPhotoOff} from '@tabler/icons-react';

interface EmptyStateProps {
    title: string;
    description?: string;
    icon?: React.ReactNode;
    action?: {
        label: string;
        onClick: () => void;
    };
}

export const EmptyState: React.FC<EmptyStateProps> = ({ title, description, icon, action }) => {
    return (
        <Flex justify="center" align="center" direction="column" wrap="wrap" w="100%" py="xl">
            <ThemeIcon radius="md" size="xl" color="gray" variant="light">
                {icon || <IconPhotoOff size="2rem" />}
            </ThemeIcon>
            <Text size="lg" fw={500} mt="md">
                {title}
            </Text>
            {description && (
                <Text size="sm" c="dimmed" ta="center" maw={400} mt="xs">
                    {description}
                </Text>
            )}
            {action && (
                <Button mt="md" onClick={action.onClick}>
                    {action.label}
                </Button>
            )}
        </Flex>
    );
};
