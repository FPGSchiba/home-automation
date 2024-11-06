import * as React from "react";
import { Link } from "react-router-dom";
import Button from "@mui/material/Button";
import {useNotificationStore} from "../store/notification";
import Typography from "@mui/material/Typography";

export default function Home() {
    const notify = useNotificationStore((store) => store.notify);

    return (
        <div className="home">
            <Typography className="home home-title" variant="h1" component="h1">Home Page</Typography>
            <Link to="/about">About page</Link>
            <Button variant="contained" onClick={() => notify({message: "testing", title: "testing", level: "info"})}>Contained</Button>
        </div>
    );
}