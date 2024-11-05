import * as React from "react";
import Typography from "@mui/material/Typography";
import {Paper, TextField, Button} from "@mui/material";
import {useNavigate} from "react-router-dom";
import originLogo from '../resources/images/origin-logo.webp';

export default function Login() {
    const navigate = useNavigate();
    const submit = (event: { preventDefault: () => void; }) => {
        event.preventDefault();
        console.log("Submit");
    }

    return (
        <Paper className="login login-paper">
            <Paper className="login login-logo login-logo-wrapper" elevation={2} >
                <img src={originLogo} alt="Logo" className="login login-logo login-logo-img" />
            </Paper>
            <Typography variant="h4" component="h1" className="login login-header">Login</Typography>
            <form className="login login-form login-form-wrapper" onSubmit={submit}>
                <TextField
                    error={false}
                    id="username"
                    label="Email"
                    helperText="Your Email address"
                    className="login login-form login-form-input"
                />
                <TextField
                    error={false}
                    id="password"
                    label="Password"
                    type="password"
                    helperText="Your Password"
                    className="login login-form login-form-input"
                />
                <Button
                    variant="contained"
                    color="primary"
                    className="login login-form login-form-button"
                    type="submit"
                >
                    Sign In
                </Button>
            </form>
            <Typography
                variant="body1"
                component="p"
                className="login login-form login-form-reset"
                onClick={() => navigate("/reset-password")}
            >Forgot password?</Typography>
        </Paper>
    );
}