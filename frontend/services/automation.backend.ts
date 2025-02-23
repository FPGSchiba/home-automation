/* eslint-disable @typescript-eslint/no-explicit-any */
import axios, { AxiosInstance } from "axios";
import https from 'https';
import {IUserInfo} from "../store/types";
import config from "../conf.yaml";
import {Permission} from "../store/user";

export enum ApiStatus {
    SUCCESS = 'success',
    ERROR = 'error',
}

class AutomationAPI {
    protected static userEndpoint: AxiosInstance;
    protected static authEndpoint: AxiosInstance;
    protected static mealEndpoint: AxiosInstance;
    protected static financeEndpoint: AxiosInstance;
    private token: string;
    private static instance: AutomationAPI;


    static getInstance(): AutomationAPI {
        if (!AutomationAPI.instance) {
            AutomationAPI.instance = new AutomationAPI().init();
        }
        return AutomationAPI.instance;
    }

    public init(): this {
        if (AutomationAPI.userEndpoint !== undefined &&
            AutomationAPI.authEndpoint !== undefined &&
            AutomationAPI.mealEndpoint !== undefined &&
            AutomationAPI.financeEndpoint !== undefined) {
            return;
        }

        // At request level
        const httpsAgent = new https.Agent({
            rejectUnauthorized: false,
        });

        const apiPath = config.frontend["api-path"];
        const authPath = config.frontend["auth-path"];
        const userApiHost = config.frontend["user-api-host"];
        const mealApiHost = config.frontend["meal-api-host"];
        const financeApiHost = config.frontend["finance-api-host"];

        AutomationAPI.userEndpoint = axios.create({
            baseURL: userApiHost + apiPath,
            httpsAgent,
        });

        AutomationAPI.authEndpoint = axios.create({
            baseURL: userApiHost + authPath,
            httpsAgent,
        });

        AutomationAPI.mealEndpoint = axios.create({
            baseURL: mealApiHost + apiPath,
            httpsAgent,
        });

        AutomationAPI.financeEndpoint = axios.create({
            baseURL: financeApiHost + apiPath,
            httpsAgent,
        });

        return this;
    }

    public setToken(token: string): void {
        this.token = token;
    }

    public async login(email: string, password: string): Promise<{ message: string, status: ApiStatus, token?: string, user?: IUserInfo }> {
        try {
            const response = await AutomationAPI.authEndpoint.post('/login', {
                email,
                password,
            });
            return response.data;
        }
        catch (reason) {
            return {message: reason.response.data.message, status: ApiStatus.ERROR};
        }
    }

    public async listPermissions(): Promise<{ message: string, status: ApiStatus, permissions?: Permission[] }> {
        try {
            const response = await AutomationAPI.userEndpoint.get('/permissions', {
                headers: {
                    Authorization: `Bearer ${this.token}`,
                },
            });
            return response.data;
        }
        catch (reason) {
            return {message: reason.response.data.message, status: ApiStatus.ERROR};
        }
    }
}

export default AutomationAPI.getInstance();
