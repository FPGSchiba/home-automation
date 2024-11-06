import React from 'react';
import {Navigate, Outlet} from 'react-router-dom';
import {useUserStore} from "../store/user";

const PrivateRoute = () => {
    const authenticated = useUserStore((state) => state.authenticated);

    // If authorized, return an outlet that will render child elements
    // If not, return element that will navigate to login page
    return authenticated ? <Outlet /> : <Navigate to={"/login"} />
}

export default PrivateRoute;
