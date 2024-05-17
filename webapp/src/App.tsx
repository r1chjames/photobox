import React from "react";
import {AppShell, Burger, MantineProvider} from "@mantine/core";
import {useDisclosure} from "@mantine/hooks";
import {QueryClient, QueryClientProvider} from "react-query";
import {theme} from "./theme";
import Router from "./Routing/Router";

interface IProps {
    baseApiUrl: string;
}

export const App: React.FunctionComponent<IProps> = (props) => {
    const queryClient = new QueryClient();
    const [opened, {toggle}] = useDisclosure(false);

    return (
        <Router baseApiUrl={props.baseApiUrl}>
            <QueryClientProvider client={queryClient}>
                <MantineProvider theme={theme}>

                    <AppShell
                        header={{height: 60}}
                        navbar={{
                            width: 300,
                            breakpoint: 'sm',
                            collapsed: {mobile: !opened},
                        }}
                        padding="md"
                    >
                        <AppShell.Header>
                            <Burger
                                opened={opened}
                                onClick={toggle}
                                hiddenFrom="sm"
                                size="sm"
                            />
                            <div>Logo</div>
                        </AppShell.Header>

                        <AppShell.Navbar p="md">Navbar</AppShell.Navbar>

                        <AppShell.Main>
                            Something
                        </AppShell.Main>
                    </AppShell>
                </MantineProvider>
            </QueryClientProvider>
        </Router>
    );
}

export default App;
