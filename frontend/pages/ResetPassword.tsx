import * as React from "react";
import {Button, Paper, TextField} from "@mui/material";
import originLogo from "../resources/images/origin-logo.webp";
import Typography from "@mui/material/Typography";
import {z} from "zod";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {useUserStore} from "../store/user";
import {useNotificationStore} from "../store/notification";
import {useNavigate} from "react-router-dom";

export default function ResetPassword() {
    const resetPassword = useUserStore((store) => store.resetPassword);
    const notify = useNotificationStore((store) => store.notify);
    const navigate = useNavigate();

    const schema = z.object({
        email: z.string().email("Invalid email").min(1, "Email is required"),
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
        },
    });

    const onSubmit = (data: FormData) => {
        notify({message: "Not yet implemented: Reset password", title: "Not implemented", level: "warning"});
        resetPassword(data.email).then((success) => {
            reset();
        });
    };

    return <Paper className="reset reset-paper">
        <Paper className="reset reset-logo reset-logo-wrapper" elevation={2} >
            <img src={originLogo} alt="Logo" className="reset reset-logo reset-logo-img" />
        </Paper>
        <Typography variant="h4" component="h1" className="reset reset-header">Reset Password</Typography>
        <form className="reset reset-form reset-form-wrapper" onSubmit={handleSubmit(onSubmit)}>
            <TextField
                label="Email"
                className="reset reset-form reset-form-input"
                defaultValue={""}
                {...register("email")}
                error={!!errors.email && touchedFields.email}
                helperText={errors.email?.message}
            />
            <Button
                variant="contained"
                color="primary"
                className="reset reset-form reset-form-button"
                type="submit"
            >
                Submit
            </Button>
        </form>
        <Button
            variant="contained"
            color="primary"
            className="reset reset-login"
            onClick={() => navigate("/login")}
        >
            Back to Login
        </Button>
    </Paper>;
}