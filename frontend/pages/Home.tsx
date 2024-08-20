import * as React from "react";
import { Link } from "react-router-dom";
import Button from "@mui/material/Button";
import {useNotificationStore} from "../store";

export default function Home() {
    const notify = useNotificationStore((store) => store.notify);

    return (
        <div>
            <Link to="/about">About page</Link>
            <Button variant="contained" onClick={() => notify({message: "testing", title: "testing", level: "info", id: "test"})}>Contained</Button>
        </div>
    );
}