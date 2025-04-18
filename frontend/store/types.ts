import {OverridableStringUnion} from "@mui/types";
import {AlertColor, AlertPropsColorOverrides} from "@mui/material";

export interface Notification {
    message: string
    level: OverridableStringUnion<AlertColor, AlertPropsColorOverrides>
    title: string
    id: string
}

export interface NotifyEvent {
    message: string
    level: OverridableStringUnion<AlertColor, AlertPropsColorOverrides>
    title: string
}

export interface IErrorDetail {
    field: string
    message: string
}

export interface IUser {
    id: string
    displayName: string
    email: string
    profilePictureUrl: string
}

export interface IUserInfo  extends IUser {
    token: string
}

export interface IBackupJob {
    id: string,
    name: string,
    identifier: string,
    configuration: any
    schedule: string
    schedulerId: string
}

export interface IBackupJobCreate {
    name: string,
    identifier: string,
    configuration: any
    schedule: string
}