// define the child properties and their types.
type globalAppVariables = {
    baseApiUrl: string;
    // more can go here.
};

declare global {
    // eslint-disable-next-line no-var
    var app: globalAppVariables;
}

declare module "*.module.css";

// set the values.
globalThis.app = {
    baseApiUrl: import.meta.env.VITE_API_URL || "http://localhost:8080/api"
};

// Freeze so these can only be defined in this file.
Object.freeze(globalThis.app);
