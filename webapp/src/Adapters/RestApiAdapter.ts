import axios, {AxiosRequestConfig, AxiosError} from 'axios';

export interface IRestApiAdapter {

    postApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>): Promise<any>

    putApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>): Promise<any>

    deleteApiCall(path: string, headers: Record<string, string>): Promise<any>

    getApiCall(path: string, headers: Record<string, string>, params: Record<string, unknown>): Promise<any>

    getPaginatedApiCall(path: string, headers: Record<string, string>, params: Record<string, unknown>): Promise<any>

    getBinaryApiCall(path: string, headers: Record<string, string>, params: Record<string, unknown>): Promise<any>

    postBinaryApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>): Promise<any>

    authHeader(): Record<string, string>

    getBaseApiPath(): string
}

let globalOnUnauthorized: (() => void) | undefined;

function handleUnauthorized(error: AxiosError<{ error?: string }>) {
    if (error.response?.status === 401) {
        localStorage.removeItem("token");
        globalOnUnauthorized?.();
    }
}

export class RestApiAdapter implements IRestApiAdapter {

    private readonly baseApiPath: string;

    constructor(baseApiPath: string, onUnauthorized?: () => void) {
        this.baseApiPath = baseApiPath;
        if (onUnauthorized) {
            globalOnUnauthorized = onUnauthorized;
        }
    }

    async postApiCall<T>(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
        return this.apiCall<T>('post', headers, body, path, {});
    }

    async putApiCall<T>(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
        return this.apiCall<T>('put', headers, body, path, {});
    }

    async deleteApiCall<T>(path: string, headers: Record<string, string>) {
        return this.apiCall<T>('delete', headers, {}, path, {});
    }

    async getApiCall<T>(path: string, headers: Record<string, string>, params: Record<string, any>) {
        return this.apiCall<T>('get', headers, {}, path, params);
    }

    async getPaginatedApiCall<T>(path: string, headers: Record<string, string>, params: Record<string, any>) {
        return this.paginatedApiCall<T>('get', headers, {}, path, params);
    }

    async getBinaryApiCall(path: string, headers: Record<string, string>, params: Record<string, any>) {
        return this.binaryApiCall('get', headers, {}, path, params);
    }

    async postBinaryApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
        return this.binaryApiCall('post', headers, body, path, {});
    }

    private async binaryApiCall(method: string, headers: Record<string, string>, body: Record<string, any>, url: string, params: Record<string, any>) {
        const parsedUrl = `${this.baseApiPath}/${url}`;
        const options: AxiosRequestConfig = {
            method,
            headers,
            data: JSON.stringify(body),
            url: parsedUrl,
            responseType: 'blob',
            params,
        };
        try {
            const response = await axios(options);
            return response.data;
        } catch (error) {
            handleUnauthorized(error as AxiosError<{ error?: string }>);
            throw error;
        }
    }

    private async apiCall<T>(method: string, headers: Record<string, string>, body: Record<string, any>, url: string, params: Record<string, any>) {
        const parsedUrl = `${this.baseApiPath}/${url}`;
        const options: AxiosRequestConfig = {
            method,
            headers,
            data: JSON.stringify(body),
            url: parsedUrl,
            params,
        };
        try {
            const response = await axios<ServerResponse<T>>(options);
            return response.data.data;
        } catch (error) {
            handleUnauthorized(error as AxiosError<{ error?: string }>);
            throw error;
        }
    }

    private async paginatedApiCall<T>(method: string, headers: Record<string, string>, body: Record<string, any>, url: string, params: Record<string, any>) {
        const parsedUrl = `${this.baseApiPath}/${url}`;
        const options: AxiosRequestConfig = {
            method,
            headers,
            data: JSON.stringify(body),
            url: parsedUrl,
            params,
        };
        try {
            const response = await axios<PaginatedServerResponse<T>>(options);
            return response.data.data;
        } catch (error) {
            handleUnauthorized(error as AxiosError<{ error?: string }>);
            throw error;
        }
    }

    authHeader() {
        return {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    }

    getBaseApiPath() {
        return this.baseApiPath;
    }
}

interface ServerResponse<T> {
    success: boolean;
    message: string;
    data: T
}

interface PaginatedMetadata {
    fromId: string;
    toId: string;
    count: number;
    nextPage: string;
}

interface PaginatedServerResponse<T> extends ServerResponse<T>{
    metaData: PaginatedMetadata
}