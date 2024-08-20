import React from "react";
import * as Sentry from "@sentry/react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Home from "./pages/Home";
import About from "./pages/About";
import { Header } from "./components/Header";
import { Notification } from "./components/Notification";
import {createTheme, ThemeProvider} from "@mui/material/styles";

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
        primary: {
            main: '#2334dc'
        }
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
                        <Route index element={<Home />} />
                        <Route path="/about" element={<About />} />
                    </Routes>
                </BrowserRouter>
            </ThemeProvider>
        );
    }
}

export default Sentry.withProfiler(App);