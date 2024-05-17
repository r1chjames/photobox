import React from "react";
import Router from "./Routing/Router";

interface IProps {
    baseApiUrl: string;
}

export const App: React.FunctionComponent<IProps> = (props) => {
    return (
        <React.StrictMode>
           <Router baseApiUrl={props.baseApiUrl}/>
        </React.StrictMode>
    );
}

export default App;
