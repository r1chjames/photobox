import axios, {AxiosRequestConfig} from 'axios';

export interface IRestApiAdapter {

    postApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>): Promise<any>

    putApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>): Promise<any>

    getApiCall(path: string, headers: Record<string, string>, params: Record<string, unknown>): Promise<any>

    getPaginatedApiCall(path: string, headers: Record<string, string>, params: Record<string, unknown>): Promise<any>

    getBinaryApiCall(path: string, headers: Record<string, string>, params: Record<string, unknown>): Promise<any>

    authHeader(): Record<string, string>
}

export class RestApiAdapter implements IRestApiAdapter {

    private readonly baseApiPath: string;

    constructor(baseApiPath: string) {
        this.baseApiPath = baseApiPath;
    }

    async postApiCall<T>(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
        return this.apiCall<T>('post', headers, body, path, {});
    }

    async putApiCall<T>(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
        return this.apiCall<T>('put', headers, body, path, {});
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
        return await axios(options)
            .then(function (response) {
                return response.data;
            })
            .catch(function (error) {
                if (error.response) {
                    if (error.response.data.error === "Access Token Not Valid") {
                        localStorage.removeItem("token");
                    }
                }
            });
    };

    private async apiCall <T>(method: string, headers: Record<string, string>, body: Record<string, any>, url: string, params: Record<string, any>) {
        const parsedUrl = `${this.baseApiPath}/${url}`;
        const options: AxiosRequestConfig = {
            method,
            headers,
            data: JSON.stringify(body),
            url: parsedUrl,
            params,
        };
        return await axios<ServerResponse<T>>(options)
            .then(function (response) {
                return response.data.data;
            })
            .catch(function (error) {
                if (error.response) {
                    if (error.response.data.error === "Access Token Not Valid") {
                        localStorage.removeItem("token");
                    }
                }
            });
    };

    private async paginatedApiCall <T>(method: string, headers: Record<string, string>, body: Record<string, any>, url: string, params: Record<string, any>) {
        const parsedUrl = `${this.baseApiPath}/${url}`;
        const options: AxiosRequestConfig = {
            method,
            headers,
            data: JSON.stringify(body),
            url: parsedUrl,
            params,
        };
        return await axios<PaginatedServerResponse<T>>(options)
            .then(function (response) {
                return response.data.data;
            })
            .catch(function (error) {
                if (error.response) {
                    if (error.response.data.error === "Access Token Not Valid") {
                        localStorage.removeItem("token");
                    }
                }
            });
    };


    authHeader() {
        return {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
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