import React from "react";
import * as Sentry from "@sentry/react";
import {BrowserRouter, Navigate, Route, Routes} from "react-router-dom";
import Home from "./pages/Home";
import { Header } from "./components/Header";
import { Notification } from "./components/Notification";
import {createTheme, ThemeProvider} from "@mui/material/styles";
import PrivateRoute from "./components/PrivateRoutes";
import Login from "./pages/Login";
import ResetPassword from "./pages/ResetPassword";
import NotFound from "./pages/NotFound";
import Settings from "./pages/Settings";
import Profile from "./pages/Profile";
import UsersDashboard from "./pages/Users/UsersDashboard";
import RolesDashboard from "./pages/Users/RolesDashboard";
import FinanceDashboard from "./pages/Finance/FinanceDashboard";
import MealDashboard from "./pages/Meals/MealDashboard";
import UserNavigation from "./pages/Users/UserNavigation";
import UserDetails from "./pages/Users/UserDetails";
import UsersList from "./pages/Users/UsersList";
import BackupDashboard from "./pages/Backups/BackupDashboard";
import BackupList from "./pages/Backups/BackupList";
import BackupSettings from "./pages/Backups/BackupSettings";
import BackupJobList from "./pages/Backups/BackupJobList";
import BackupJobDetails from "./pages/Backups/BackupJobDetails";
import BackupDetails from "./pages/Backups/BackupDetails";
import BackupNavigation from "./pages/Backups/BackupNavigation";


const darkTheme = createTheme({
    palette: {
        mode: 'dark',
        primary: {
            main: '#2334dc'
        },
        secondary: {
            main: '#22c356'
        },
        text: {
            primary: '#ffffff',
            secondary: '#ffffff'
        },
        background: {
            default: '#121212',
            paper: '#333333'
        },
        divider: '#ffffff',
        success: {
            main: '#4caf50',

        },
        error: {
            main: '#f44336',
        },
        warning: {
            main: '#ff9800',
        },
        info: {
            main: '#2196f3',
        },
    },
});

class App extends React.Component {
    render() {
        return (
            <ThemeProvider theme={darkTheme}>
                <BrowserRouter>
                    <Header />
                    <Notification />
                    <Routes>
                        <Route path={"/"} element={<PrivateRoute />}>
                            <Route index element={<Home />} />
                            <Route path="/home" element={<Home />} />
                            <Route path={"/settings"} element={<Settings />} />
                            <Route path={"/profile"} element={<Profile />} />
                            <Route path={"/users"} element={<UserNavigation />} >
                                <Route path={""} element={<UsersDashboard />} />
                                <Route path={"roles"} element={<RolesDashboard />} />
                                <Route path={"list"} element={<UsersList />} />
                                <Route path={"user/:id"} element={<UserDetails />} />
                            </Route>
                            <Route path={"/finance"} element={<FinanceDashboard />} />
                            <Route path={"/meal"} element={<MealDashboard />} />
                            <Route path={"/backups"} element={<BackupNavigation />} >
                                <Route path={""} element={<BackupDashboard />} />
                                <Route path={"list"} element={<BackupList />} />
                                <Route path={"settings"} element={<BackupSettings />} />
                                <Route path={"jobs"} element={<BackupJobList />} />
                                <Route path={"jobs/:id"} element={<BackupJobDetails />} />
                                <Route path={":id"} element={<BackupDetails />} />
                            </Route>
                        </Route>
                        <Route path="/login" element={<Login />} />
                        <Route path="/reset-password" element={<ResetPassword />} />
                        <Route path="/not-found" element={<NotFound />} />
                        <Route path="*" element={<Navigate  to="/not-found" />} />
                    </Routes>
                </BrowserRouter>
            </ThemeProvider>
        );
    }
}

export default Sentry.withProfiler(App);