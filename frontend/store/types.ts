import {OverridableStringUnion} from "@mui/types";
import {AlertColor, AlertPropsColorOverrides} from "@mui/material";

export type Notification = {
    message: string
    level: OverridableStringUnion<AlertColor, AlertPropsColorOverrides>
    title: string
    id: string
}

export type NotifyEvent = {
    message: string
    level: OverridableStringUnion<AlertColor, AlertPropsColorOverrides>
    title: string
}
