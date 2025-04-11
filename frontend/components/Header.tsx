import * as React from 'react';
import {useEffect} from 'react';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';
import {useLocation, useNavigate} from "react-router-dom";
import {Avatar, Menu, MenuItem, Tooltip} from "@mui/material";
import {deepOrange} from "@mui/material/colors";
import {useUserStore} from "../store/user";
import Box from '@mui/material/Box';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import {FrontendPermissions, hasPermission} from "../services/permissions";
import vngdLogo from '../resources/images/vanguardSemperDucimusLogo.png';

const nonHeaderLocations = ["/login", "/reset-password"];

export function Header() {
    const displayName = useUserStore((store) => store.user?.displayName);
    const permissions = useUserStore((store) => store.permissions);
    const location = useLocation();
    const navigate = useNavigate();
    const [anchorElUser, setAnchorElUser] = React.useState<null | HTMLElement>(null);
    const [value, setValue] = React.useState("/");

    const handleChange = (event: React.SyntheticEvent, newValue: string) => {
        // @ts-ignore
        setValue(newValue);
        navigate(newValue);
    };

    const settings = ['Profile', 'Logout'];

    const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorElUser(event.currentTarget);
    };

    const handleCloseUserMenu = (setting: string) => {
        setAnchorElUser(null);
        switch (setting) {
            case 'Profile':
                navigate('/profile');
                break;
            case 'Logout':
                useUserStore.getState().logout();
                break;
        }
    };

    useEffect(() => {
        const groupPart = /\/\w*\/?/.exec(location.pathname);
        if (groupPart) {
            switch (groupPart[0].replaceAll('/', '')) {
                case 'home':
                    setValue('/home');
                    break;
                case '':
                    setValue('/home');
                    break;
                case 'finance':
                    setValue('/finance');
                    break;
                case 'meal':
                    setValue('/meal');
                    break;
                case 'settings':
                    setValue('/settings');
                    break;
                case 'users':
                    setValue('/users');
                    break;
                case 'backups':
                    setValue('/backups');
                    break;
                default:
                    setValue('/');
                    break;
            }
        }
    }, [location]);

    return (
        nonHeaderLocations.includes(location.pathname) ? null :
            <Box sx={{ flexGrow: 1 }} className={"header header-container"}>
                <AppBar position="static" className={"header header-bar"}>
                    <Toolbar>
                        <img src={vngdLogo} className="header header-logo"  alt="Vanguard Logo"/>
                        <Box sx={{ width: '100%' }}>
                            <Tabs
                                value={value}
                                onChange={handleChange}
                                indicatorColor="secondary"
                                textColor="secondary"
                                aria-label="secondary tabs example"
                            >
                                <Tab value="/home" label="Home" />
                                {hasPermission(FrontendPermissions.backupsRead, permissions) ? <Tab value="/backups" label="Backups" /> : <></>}
                                {hasPermission(FrontendPermissions.financeRead, permissions) ? <Tab value="/finance" label="Finance" /> : <></>}
                                {hasPermission(FrontendPermissions.mealRead, permissions) ? <Tab value="/meal" label="Meals" /> : <></>}
                                {hasPermission(FrontendPermissions.settingsRead, permissions) ? <Tab value="/settings" label="Settings"/> : <></>}
                                {hasPermission(FrontendPermissions.usersRead, permissions) ? <Tab value="/users" label="Users"/> : <></>}
                            </Tabs>
                        </Box>
                        <Box sx={{ flexGrow: 0 }}>
                            <Tooltip title="Open settings">
                                <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                                    <Avatar sx={{ bgcolor: deepOrange[400] }} alt={displayName || ""} >{displayName ? displayName[0] : "NoN"}</Avatar>
                                </IconButton>
                            </Tooltip>
                            <Menu
                                sx={{ mt: '45px' }}
                                id="menu-appbar"
                                anchorEl={anchorElUser}
                                anchorOrigin={{
                                    vertical: 'top',
                                    horizontal: 'right',
                                }}
                                keepMounted
                                transformOrigin={{
                                    vertical: 'top',
                                    horizontal: 'right',
                                }}
                                open={Boolean(anchorElUser)}
                                onClose={handleCloseUserMenu}
                            >
                                {settings.map((setting) => (
                                    <MenuItem key={setting} onClick={() => {handleCloseUserMenu(setting)}}>
                                        <Typography sx={{ textAlign: 'center' }}>{setting}</Typography>
                                    </MenuItem>
                                ))}
                            </Menu>
                        </Box>
                    </Toolbar>
                </AppBar>
            </Box>
    )
}