import React from "react";
import {MantineProvider} from "@mantine/core";
import {QueryClient, QueryClientProvider} from "react-query";
import {theme} from "./theme";
import Router from "./Routing/Router";
import {AppBar} from "./Components/AppBar/AppBar";

interface IProps {
    baseApiUrl: string;
}

export const App: React.FunctionComponent<IProps> = (props) => {
    const queryClient = new QueryClient();


    return (
        <Router baseApiUrl={props.baseApiUrl}>
            <QueryClientProvider client={queryClient}>
                <MantineProvider theme={theme}>
                    <AppBar activeLink={"Albums"}/>
                </MantineProvider>
            </QueryClientProvider>
        </Router>
    );
}

export default App;
