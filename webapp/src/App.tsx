 import { MantineProvider } from '@mantine/core';
import {Dashboard} from "./Components/Dashboard/Dashboard";
import React from "react";

interface IProps {
    baseApiUrl: string;
}

export const App: React.FunctionComponent<IProps> = (props) => {
    return (
        <MantineProvider theme={{ loader: 'bars' }}>
            <Dashboard baseApiUrl={props.baseApiUrl}/>
        </MantineProvider>
    );
}

export default App;
