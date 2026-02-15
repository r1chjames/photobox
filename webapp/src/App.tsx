import React from "react";
import {MantineProvider} from "@mantine/core";
import {theme} from "./theme";
import Router from "./Routing/Router";
import {AdapterProvider} from "./Routing/AdapterContext";
import {ErrorBoundary} from "./Components/ErrorBoundary/ErrorBoundary";

interface IProps {
    baseApiUrl: string;
}

export const App: React.FunctionComponent<IProps> = (props) => {

    return (
        <MantineProvider theme={theme}>
            <ErrorBoundary>
                <AdapterProvider baseApiUrl={props.baseApiUrl}>
                    <Router/>
                </AdapterProvider>
            </ErrorBoundary>
        </MantineProvider>
    );
}

export default App;
