import * as React from "react";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import {Avatar, Chip, Paper, SelectChangeEvent, Stack} from "@mui/material";
import {DataGrid, GridColDef} from "@mui/x-data-grid";
import {deepOrange} from "@mui/material/colors";
import Button from "@mui/material/Button";
import {stringToColour} from "../../services/util";

const columns: GridColDef[] = [
    {
        field: 'id',
        headerName: 'ID',
        width: 70,
    },
    {
        field: 'avatar',
        headerName: 'Avatar',
        width: 70,
        sortable: false,
        filterable: false,
        hideable: false,
        disableColumnMenu: true,
        renderCell: (params) => <Avatar sx={{ bgcolor: deepOrange[400] }} alt={params.row.displayName} >{params.row.displayName[0]}</Avatar>,
    },
    { field: 'displayName', headerName: 'Display Name', width: 200 },
    { field: 'email', headerName: 'Email Address', width: 300 },
    {
        field: 'roles',
        headerName: 'Roles',
        width: 200,
        renderCell: (params) => {
            return <Stack direction="row" spacing={1}  className="users users-list users-list-table users-list-table-rolestack">{params.row.roles.map((role: string) => <Chip sx={{ bgcolor: stringToColour(role), color: 'white' }} size="small"  className="users users-list users-list-table users-list-table-rolechip" label={role} key={role} />)}</Stack>;
        }
    },
    {
        field: 'actions',
        headerName: 'Actions',
        width: 200,
        filterable: false,
        hideable: false,
        sortable: false,
        disableColumnMenu: true,
        resizable: false,
        renderCell: (params) => {
            return <div>
                <Button sx={{marginRight: "20px"}} variant="contained" color="info">Edit</Button>
                <Button variant="contained" color="error">Delete</Button>
            </div>;
        },
        valueGetter: (value, row) => `${row.email}`,
    },
];

const rows = [
    { id: 1, displayName: 'John Snow', email: 'john.snow@example.com', roles: ['Admin', 'User'] },
    { id: 2, displayName: 'Arya Stark', email: 'arya.stark@example.com', roles: ['User'] },
    { id: 3, displayName: 'Tyrion Lannister', email: 'tyrion.lannister@example.com', roles: ['Admin'] },
    { id: 4, displayName: 'Daenerys Targaryen', email: 'daenerys.targaryen@example.com', roles: ['User', 'Manager'] },
    { id: 5, displayName: 'Cersei Lannister', email: 'cersei.lannister@example.com', roles: ['Manager'] },
    { id: 6, displayName: 'Bran Stark', email: 'bran.stark@example.com', roles: ['User'] },
    { id: 7, displayName: 'Sansa Stark', email: 'sansa.stark@example.com', roles: ['User', 'Admin'] },
    { id: 8, displayName: 'Jaime Lannister', email: 'jaime.lannister@example.com', roles: ['User'] },
    { id: 9, displayName: 'Samwell Tarly', email: 'samwell.tarly@example.com', roles: ['User'] },
];

const paginationModel = { page: 0, pageSize: 10 };

export default function UsersList() {
    const [age, setAge] = React.useState('');

    const handleChange = (event: SelectChangeEvent) => {
        setAge(event.target.value);
    };

    return <Box className="users users-list users-list-wrapper">
        <Typography className="users users-list users-list-header" variant="h2" component="h1">Users</Typography>
        <Paper className="users users-list users-list-table users-list-table-wrapper" elevation={2}>
            <DataGrid
                className="users users-list users-list-table users-list-table-grid"
                rows={rows}
                columns={columns}
                initialState={{ pagination: { paginationModel } }}
                pageSizeOptions={[10, 25, 50]}
            />
        </Paper>
    </Box>;
}