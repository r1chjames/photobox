import React from "react";
import {MantineProvider} from "@mantine/core";
import {Notifications} from "@mantine/notifications";
import {ModalsProvider} from "@mantine/modals";
import {theme} from "./theme";
import Router from "./Routing/Router";
import {AdapterProvider} from "./Routing/AdapterContext";
import {ErrorBoundary} from "./Components/ErrorBoundary/ErrorBoundary";
import '@mantine/notifications/styles.css';

interface IProps {
    baseApiUrl: string;
}

export const App: React.FunctionComponent<IProps> = (props) => {

    return (
        <MantineProvider theme={theme}>
            <Notifications position="top-right" zIndex={1000}/>
            <ModalsProvider>
                <ErrorBoundary>
                    <AdapterProvider baseApiUrl={props.baseApiUrl}>
                        <Router/>
                    </AdapterProvider>
                </ErrorBoundary>
            </ModalsProvider>
        </MantineProvider>
    );
}

export default App;
