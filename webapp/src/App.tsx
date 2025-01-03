import React from "react";
import {MantineProvider} from "@mantine/core";
import {theme} from "./theme";
import Router from "./Routing/Router";

interface IProps {
    baseApiUrl: string;
}

export const App: React.FunctionComponent<IProps> = (props) => {

    return (
        <MantineProvider theme={theme}>
            <Router baseApiUrl={props.baseApiUrl}/>
        </MantineProvider>
    );
}

export default App;
