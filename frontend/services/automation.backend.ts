/* eslint-disable @typescript-eslint/no-explicit-any */
import axios, { AxiosInstance } from "axios";
import https from 'https';
import {IUserInfo} from "../store/types";
const UNAUTHORIZED_CODE = 401;

/**
 * Initialization needs to be done before calling any method,
 * @param target
 * @param propertyKey
 * @param descriptor
 */
function wrapInit(target: AutomationAPI, propertyKey: string | symbol, descriptor: PropertyDescriptor): void {
    const originalMethod = descriptor.value;
    const newMethod = async (...args: any[]): Promise<any> => {
        await target.init();
        return originalMethod(...args);
    };
    descriptor.value = newMethod.bind(target);
}


export enum ApiStatus {
    SUCCESS = 'success',
    ERROR = 'error',
}

class AutomationAPI {
    protected static endpoint: AxiosInstance;
    protected static authEndpoint: AxiosInstance;

    public async init(): Promise<void> {
        if (AutomationAPI.endpoint !== undefined) {
            return;
        }

        // At request level
        const httpsAgent = new https.Agent({
            rejectUnauthorized: false,
        });

        const apiURL = "http://localhost:8080/api/v1"
        const authURL = "http://localhost:8080/auth"

        AutomationAPI.endpoint = axios.create({
            baseURL: apiURL,
            httpsAgent,
        });

        AutomationAPI.authEndpoint = axios.create({
            baseURL: authURL,
            httpsAgent,
        });
    }

    public async login(email: string, password: string): Promise<{ message: string, status: ApiStatus, token?: string, user?: IUserInfo }> {
        return AutomationAPI.authEndpoint.post('/login', {
            email,
            password,
        });
    }
}

export default new AutomationAPI();
