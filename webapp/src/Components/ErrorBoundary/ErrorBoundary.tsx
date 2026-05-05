import React from 'react';
import {Button, Container, Text, Title, Stack, Group} from '@mantine/core';

interface ErrorBoundaryState {
    hasError: boolean;
    error: Error | null;
}

export class ErrorBoundary extends React.Component<{ children: React.ReactNode }, ErrorBoundaryState> {
    constructor(props: { children: React.ReactNode }) {
        super(props);
        this.state = {hasError: false, error: null};
    }

    static getDerivedStateFromError(error: Error): ErrorBoundaryState {
        return {hasError: true, error};
    }

    handleReset = () => {
        this.setState({hasError: false, error: null});
    };

    render() {
        if (this.state.hasError) {
            return (
                <Container size="sm" py="xl">
                    <Stack align="center" gap="md">
                        <Title order={2}>Something went wrong</Title>
                        <Text c="dimmed">{this.state.error?.message}</Text>
                        <Group gap="sm">
                            <Button variant="default" onClick={() => window.history.back()}>Go back</Button>
                            <Button onClick={this.handleReset}>Try again</Button>
                        </Group>
                    </Stack>
                </Container>
            );
        }

        return this.props.children;
    }
}
