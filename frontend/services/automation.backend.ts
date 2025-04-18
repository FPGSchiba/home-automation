/* eslint-disable @typescript-eslint/no-explicit-any */
import axios, { AxiosInstance } from "axios";
import https from 'https';
import {IBackupJob, IBackupJobCreate, IErrorDetail, IUserInfo} from "../store/types";
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
    protected static backupEndpoint: AxiosInstance;
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
            AutomationAPI.financeEndpoint !== undefined &&
            AutomationAPI.backupEndpoint !== undefined) {
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
        const backupApiHost = config.frontend["backup-api-host"];

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

        AutomationAPI.backupEndpoint = axios.create({
            baseURL: backupApiHost + apiPath,
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

    public async getFinanceVersion(): Promise<{ message: string, status: ApiStatus, version?: string }> {
        try {
            const response = await AutomationAPI.financeEndpoint.get('/', {
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

    public async getBackupJobs(): Promise<{ message: string, status: ApiStatus, jobs?: IBackupJob[] }> {
        try {
            const reponse = await AutomationAPI.backupEndpoint.get('/jobs', {
                headers: {
                    Authorization: `Bearer ${this.token}`,
                },
            });
            return reponse.data;
        } catch (reason) {
            return {message: reason.response.data.message, status: ApiStatus.ERROR};
        }
    }

    public async deleteBackupJob(id: string): Promise<{ message: string, status: ApiStatus }> {
        try {
            const response = await AutomationAPI.backupEndpoint.delete(`/jobs/${id}`, {
                headers: {
                    Authorization: `Bearer ${this.token}`,
                },
            });
            return response.data;
        } catch (reason) {
            return {message: reason.response.data.message, status: ApiStatus.ERROR};
        }
    }

    public async createBackupJob(job: IBackupJobCreate): Promise<{ message: string, status: ApiStatus, job?: IBackupJob, errors?: IErrorDetail[] }> {
        try {
            const response = await AutomationAPI.backupEndpoint.post('/jobs', job, {
                headers: {
                    Authorization: `Bearer ${this.token}`,
                },
            });
            return response.data;
        } catch (reason) {
            return {message: reason.response.data.message, status: ApiStatus.ERROR, errors: reason.response.data.errors};
        }
    }
}

export default AutomationAPI.getInstance();
