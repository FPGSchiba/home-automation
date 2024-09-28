import {IUserInfo, Notification, NotifyEvent} from './types';
import {v4 as uuidv4} from 'uuid';
import { create } from 'zustand';
import automationBackend from "../services/automation.backend";

type NotificationState = {
    notifications: Notification[]
}

type NotificationActions = {
    closeNotification: (id: string) => void
    notify: (notification: NotifyEvent) => void
}

type UserState = {
    user?: IUserInfo
    token?: string
}

type UserActions = {
    setUser: (user: IUserInfo, token: string) => void
    logout: () => void
    login: (email: string, password: string) => Promise<void>
}

const useNotificationStore = create<NotificationState & NotificationActions>((set) => ({
    notifications: [],
    closeNotification: (id) => set((state) => ({notifications: state.notifications.filter(item => item.id != id)})),
    notify: (notification) => set((state) => ({notifications: [...state.notifications, {...notification, id: uuidv4()}]}))
}))

const useUserStore = create<UserState & UserActions>((set) => ({
    user: undefined,
    token: undefined,
    setUser: (user, token) => set({user, token}),
    logout: () => set({user: undefined, token: undefined}),
    login: async (email, password) => {
        const response = await automationBackend.login(email, password);
        if (response.status == 'success') {
            set({user: response.user, token: response.token})
        } else {
            useNotificationStore.getState().notify({message: response.message, level: 'error', title: 'Login Error'});
        }
    }
}))

export { useNotificationStore, useUserStore }