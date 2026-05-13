import { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from '../Routing/AuthContext';
import { MemoryRouter } from 'react-router-dom';
import { AdapterProvider } from '../Routing/AdapterContext';

const createTestQueryClient = () => new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
        mutations: {
            retry: false,
        },
    },
});

interface AllTheProvidersProps {
    children: React.ReactNode;
}

const AllTheProviders = ({ children }: AllTheProvidersProps) => {
    const queryClient = createTestQueryClient();

    return (
        <MemoryRouter>
            <AuthProvider>
                <QueryClientProvider client={queryClient}>
                    <AdapterProvider baseApiUrl="http://localhost:8080/api">
                        <MantineProvider>
                            {children}
                        </MantineProvider>
                    </AdapterProvider>
                </QueryClientProvider>
            </AuthProvider>
        </MemoryRouter>
    );
};

const customRender = (
    ui: ReactElement,
    options?: Omit<RenderOptions, 'wrapper'>,
) => render(ui, { wrapper: AllTheProviders, ...options });

export * from '@testing-library/react';
export { customRender as render, createTestQueryClient };
