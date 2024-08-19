import React from "react";
import * as Sentry from "@sentry/react";
import {BrowserRouter, Route, Routes} from "react-router-dom";
import Home from "./pages/Home";
import About from "./pages/About";
import {Header} from "./components/Header";

class App extends React.Component {
    render() {
        return (
            <BrowserRouter>
                <Header />
                <Routes>
                    <Route index element={<Home />} />
                    <Route path="/about" element={<About />} />
                </Routes>
            </BrowserRouter>
        );
    }
}

export default Sentry.withProfiler(App);