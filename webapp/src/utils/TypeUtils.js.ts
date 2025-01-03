export const isJson = (str: string) => {
    try {
        return JSON.parse(str) && !!str;
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (e) {
        return false;
    }
}

export const valueType = (value: any) => {
    if (isJson(value)) return 'json'
    if (typeof value === 'object' || Array.isArray(value)) return 'object'
    if (typeof value === 'string') return 'string'
}