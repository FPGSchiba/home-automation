import * as React from "react";
import {Paper} from "@mui/material";
import {TabContext, TabList, TabPanel} from "@mui/lab";
import Tab from "@mui/material/Tab";
import DashboardIcon from '@mui/icons-material/Dashboard';
import GroupIcon from '@mui/icons-material/Group';
import TheaterComedyIcon from '@mui/icons-material/TheaterComedy';
import {Outlet, useNavigate} from "react-router-dom";
import {useEffect} from "react";

export default function UserNavigation() {
    const [value, setValue] = React.useState('/');
    const navigate = useNavigate();

    const handleChange = (event: React.SyntheticEvent, newValue: string) => {
        setValue(newValue);
        navigate(newValue);
    };

    useEffect(() => {
        const navigationPart = location.pathname.replace(/\/users/, '').replaceAll('/', '');
        switch (navigationPart) {
            case '':
                setValue('/users');
                break;
            case '/list':
                setValue('/users/list');
                break;
            case '/roles':
                setValue('/users/roles');
                break;
            default: // Probably a User id, could check here, but should not be necessary
                setValue('/users/list');
                break;

        }
    }, [location]);

    return <div className="users users-wrapper">
        <TabContext value={value}>
            <Paper className="users users-navigation users-navigation-wrapper">
                <TabList onChange={handleChange} aria-label="lab API tabs example" orientation="vertical">
                    <Tab className="users users-navigation users-navigation-tab" icon={<DashboardIcon />} iconPosition="start" label="Dashboard" value="/users" />
                    <Tab className="users users-navigation users-navigation-tab" icon={<GroupIcon />} iconPosition="start" label="Users" value="/users/list" />
                    <Tab className="users users-navigation users-navigation-tab" icon={<TheaterComedyIcon />} iconPosition="start" label="Roles" value="/users/roles" />
                </TabList>
            </Paper>
            <Paper className="users users-content users-content-wrapper">
                <TabPanel value="/users"><Outlet /></TabPanel>
                <TabPanel value="/users/list"><Outlet /></TabPanel>
                <TabPanel value="/users/roles"><Outlet /></TabPanel>
            </Paper>
        </TabContext>
    </div>;
}