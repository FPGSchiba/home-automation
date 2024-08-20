import * as React from 'react';
import { useNotificationStore } from "../store";
import {Alert, Snackbar} from "@mui/material";
import Typography from "@mui/material/Typography";
import IconButton from "@mui/material/IconButton";
import CloseIcon from '@mui/icons-material/Close';

export function Notification() {
    const notifications = useNotificationStore((state) => state.notifications);
    const closeNotification = useNotificationStore((state) => state.closeNotification);

    return (
        <Snackbar
            anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
            open={notifications.length > 0}
        >
            <div>
                { notifications.map(function (notification) {
                    return (
                        <Alert severity={notification.level} key={notification.id} style={{overflow: "auto"}}>
                            <Typography>{notification.title}</Typography>
                            <Typography>{notification.message}</Typography>
                            <IconButton
                                size="small"
                                aria-label="close"
                                color="inherit"
                                onClick={() => closeNotification(notification.id)}
                            >
                                <CloseIcon fontSize="small"/>
                            </IconButton>
                        </Alert>
                    )
                })}
            </div>
        </Snackbar>
    )
}