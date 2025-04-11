import * as React from "react";
import {Paper} from "@mui/material";
import {TabContext, TabList, TabPanel} from "@mui/lab";
import Tab from "@mui/material/Tab";
import DashboardIcon from '@mui/icons-material/Dashboard';
import SettingsIcon from '@mui/icons-material/Settings';
import ManageHistoryIcon from '@mui/icons-material/ManageHistory';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import {Outlet, useNavigate} from "react-router-dom";
import {useEffect} from "react";

export default function BackupNavigation() {
    const [value, setValue] = React.useState('/');
    const navigate = useNavigate();

    const handleChange = (event: React.SyntheticEvent, newValue: string) => {
        setValue(newValue);
        navigate(newValue);
    };

    useEffect(() => {
        const navigationPart = location.pathname.replace(/\/backups/, '').replaceAll('/', '');
        switch (navigationPart) {
            case '':
                setValue('/backups');
                break;
            case 'list':
                setValue('/backups/list');
                break;
            case 'jobs':
                setValue('/backups/jobs');
                break;
            case 'settings':
                setValue('/backups/settings');
                break;
            default: // Probably a User id, could check here, but should not be necessary
                if (navigationPart.startsWith('jobs')) {
                    setValue('/backups/jobs');
                } else {
                    setValue('/backups/list');
                }
                break;

        }
    }, [location]);

    return <div className="users users-wrapper">
        <TabContext value={value}>
            <Paper className="users users-navigation users-navigation-wrapper">
                <TabList onChange={handleChange} aria-label="lab API tabs example" orientation="vertical">
                    <Tab className="users users-navigation users-navigation-tab" icon={<DashboardIcon />} iconPosition="start" label="Dashboard" value="/backups" />
                    <Tab className="users users-navigation users-navigation-tab" icon={<CloudDownloadIcon />} iconPosition="start" label="Backups" value="/backups/list" />
                    <Tab className="users users-navigation users-navigation-tab" icon={<ManageHistoryIcon />} iconPosition="start" label="Backup Jobs" value="/backups/jobs" />
                    <Tab className="users users-navigation users-navigation-tab" icon={<SettingsIcon />} iconPosition="start" label="Settings" value="/backups/settings" />
                </TabList>
            </Paper>
            <Paper className="users users-content users-content-wrapper">
                <TabPanel value="/backups"><Outlet /></TabPanel>
                <TabPanel value="/backups/list"><Outlet /></TabPanel>
                <TabPanel value="/backups/jobs"><Outlet /></TabPanel>
                <TabPanel value="/backups/settings"><Outlet /></TabPanel>
            </Paper>
        </TabContext>
    </div>;
}