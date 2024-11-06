import {IUser, IUserInfo} from "./types";
import {create} from "zustand";
import {checkUserInfo, eraseCookies, getUserInfoFromCookies, setUserInfoToCookies} from "../services/util";
import automationBackend, {ApiStatus} from "../services/automation.backend";
import { useNotificationStore } from "./notification";

type UserState = {
    user?: IUser
    token?: string
    authenticated: boolean
}

type UserActions = {
    setUser: (user: IUser, token: string) => void
    logout: () => void
    login: (email: string, password: string) => Promise<boolean>
    resetPassword: (email: string) => Promise<boolean>
}

const useUserStore = create<UserState & UserActions>((set) => ({
    user: checkUserInfo() ? {
        id: getUserInfoFromCookies().user.id,
        displayName: getUserInfoFromCookies().user.displayName,
        email: getUserInfoFromCookies().user.email,
        profilePictureUrl: getUserInfoFromCookies().user.profilePictureUrl
    } as IUser : undefined,
    token: checkUserInfo() ? getUserInfoFromCookies().user.token : undefined,
    authenticated: checkUserInfo(),
    setUser: (user, token) => set({user, token}),
    logout: () => {
        set({user: undefined, token: undefined, authenticated: false});
        eraseCookies();
    },
    login: async (email, password): Promise<boolean> => {
        const response = await automationBackend.login(email, password);
        if (response.status == ApiStatus.SUCCESS) {
            setUserInfoToCookies({user: response.user});
            set({user: response.user, token: response.token, authenticated: true});
            return true;
        } else {
            useNotificationStore.getState().notify({message: response.message, level: 'error', title: 'Login Error'});
            return false;
        }
    },
    resetPassword: async (email): Promise<boolean> => {
        console.log("Not yet implemented: Reset password");
        return false;
    }
}))

export {useUserStore}