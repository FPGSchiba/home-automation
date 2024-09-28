import React from "react";
import * as Sentry from "@sentry/react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Home from "./pages/Home";
import About from "./pages/About";
import { Header } from "./components/Header";
import { Notification } from "./components/Notification";
import {createTheme, ThemeProvider} from "@mui/material/styles";
import PrivateRoute from "./components/PrivateRoutes";
import Login from "./pages/Login";

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
        primary: {
            main: '#2334dc'
        },
        secondary: {
            main: '#1a2380'
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
                            <Route path="/about" element={<About />} />
                        </Route>
                        <Route path="/login" element={<Login />} />
                    </Routes>
                </BrowserRouter>
            </ThemeProvider>
        );
    }
}

export default Sentry.withProfiler(App);