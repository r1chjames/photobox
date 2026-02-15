// eslint-disable-next-line
var app: globalAppVariables;

// define the child properties and their types.
type globalAppVariables = {
    baseApiUrl: string;
    // more can go here.
};

// set the values.
globalThis.app = {
    baseApiUrl: import.meta.env.VITE_API_URL || "http://localhost:8080/api"
};

declare module "*.module.css";
// Freeze so these can only be defined in this file.
Object.freeze(globalThis.app);
