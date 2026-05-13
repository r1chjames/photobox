export const isJson = (str: string) => {
    try {
        JSON.parse(str);
        return true;
    } catch {
        return false;
    }
}

export const valueType = (value: unknown) => {
    if (typeof value === 'string' && isJson(value)) return 'json'
    if (typeof value === 'string') return 'string'
    if (typeof value === 'object' || Array.isArray(value)) return 'object'
    return 'unknown'
}
