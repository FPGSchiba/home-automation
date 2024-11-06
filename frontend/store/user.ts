import {IUser} from "./types";
import {create} from "zustand";
import {checkUserInfo, eraseCookies, getUserInfoFromCookies, setUserInfoToCookies} from "../services/util";
import automationBackend, {ApiStatus} from "../services/automation.backend";
import { useNotificationStore } from "./notification";
import {FrontendPermissions, remapPermissions} from "../services/permissions";

type UserState = {
    user?: IUser
    token?: string
    authenticated: boolean
    permissions: FrontendPermissions[]
}

type PermissionRoute = {
    path: string;
    methods: string[];
    jsonFilters: any;
}

export type Permission = {
    routes: PermissionRoute[];
}

type UserActions = {
    setUser: (user: IUser, token: string) => void
    logout: () => void
    login: (email: string, password: string) => Promise<boolean>
    resetPassword: (email: string) => Promise<boolean>
    listPermissions: () => Promise<boolean>
}

const loadUserStore = (): {
    user: IUser,
    token: string,
    authenticated: boolean,
    permissions: FrontendPermissions[]
} => {
    const authenticated = checkUserInfo();
    const user = authenticated ? getUserInfoFromCookies().user : undefined;
    const token = authenticated ? getUserInfoFromCookies().token : undefined;
    automationBackend.setToken(token);

    // Initially set permissions to an empty array
    const permissions: FrontendPermissions[] = [];

    // Fetch permissions asynchronously and update the store
    automationBackend.listPermissions().then(permissionsResponse => {
        if (permissionsResponse.status === ApiStatus.SUCCESS) {
            useUserStore.setState({ permissions: remapPermissions(permissionsResponse.permissions) });
        }
    });

    return {user, token, authenticated, permissions};
}

const defaultValues = loadUserStore();

const useUserStore = create<UserState & UserActions>((set) => ({
        user: defaultValues.user,
        token: defaultValues.token,
        authenticated: defaultValues.authenticated,
        permissions: defaultValues.permissions,
        setUser: (user, token) => set({user, token}),
        logout: () => {
            set({user: undefined, token: undefined, authenticated: false});
            eraseCookies();
        },
        login: async (email, password): Promise<boolean> => {
            const response = await automationBackend.login(email, password);
            if (response.status == ApiStatus.SUCCESS) {
                setUserInfoToCookies({user: response.user, token: response.token});
                set({user: response.user, token: response.token, authenticated: true});
                automationBackend.setToken(response.token);
                return useUserStore.getState().listPermissions();
            } else {
                useNotificationStore.getState().notify({message: response.message, level: 'error', title: 'Login Error'});
                return false;
            }
        },
        resetPassword: async (email): Promise<boolean> => {
            console.log("Not yet implemented: Reset password");
            return false;
        },
        listPermissions: async (): Promise<boolean> => {
            const response = await automationBackend.listPermissions();
            if (response.status == ApiStatus.SUCCESS) {
                set({permissions: remapPermissions(response.permissions)});
                return true;
            } else {
                useNotificationStore.getState().notify({message: response.message, level: 'error', title: 'Permission Error'});
                return false;
            }
        }
    }));

export {useUserStore};