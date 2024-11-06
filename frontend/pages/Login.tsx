import * as React from "react";
import Typography from "@mui/material/Typography";
import {Paper, TextField, Button} from "@mui/material";
import {useNavigate} from "react-router-dom";
import originLogo from '../resources/images/origin-logo.webp';
import {useUserStore} from "../store/user";
import { z } from "zod";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";

export default function Login() {
    const navigate = useNavigate();
    const login = useUserStore((store) => store.login);

    const schema = z.object({
        email: z.string().email("Invalid email").min(1, "Email is required"),
        password: z.string().min(6, "Password must be at least 6 characters long"),
    });
    type FormData = z.infer<typeof schema>;

    const {
        // register: function to register input elements
        register,
        reset,
        // handleSubmit: function to handle form submission
        handleSubmit,
        // watch: function to watch values of form inputs
        watch,
        // formState: object containing information about form state
        formState: { errors, touchedFields }, // Destructure errors and touchedFields from formState
    } = useForm<FormData>({ // Call useForm hook with generic type FormData
                            // resolver: specify resolver for form validation using Zod
        resolver: zodResolver(schema), // Pass Zod schema to resolver
        // defaultValues: specify default values for form inputs
        defaultValues: {
            email: "",
            password: "",
        },
    });

    const onSubmit = (data: FormData) => {
        login(data.email, data.password).then((success) => {
            if (success) {
                navigate("/");
            }
            reset();
        });
    };

    return (
        <Paper className="login login-paper">
            <Paper className="login login-logo login-logo-wrapper" elevation={2} >
                <img src={originLogo} alt="Logo" className="login login-logo login-logo-img" />
            </Paper>
            <Typography variant="h4" component="h1" className="login login-header">Login</Typography>
            <form className="login login-form login-form-wrapper" onSubmit={handleSubmit(onSubmit)}>
                <TextField
                    label="Email"
                    className="login login-form login-form-input"
                    defaultValue={""}
                    {...register("email")}
                    error={!!errors.email && touchedFields.email}
                    helperText={errors.email?.message}
                />
                <TextField
                    label="Password"
                    type="password"
                    className="login login-form login-form-input"
                    defaultValue={""}
                    {...register("password")}
                    error={!!errors.password && touchedFields.password}
                    helperText={errors.password?.message}
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