import * as React from "react";
import Typography from "@mui/material/Typography";
import {DataGrid, GridColDef} from "@mui/x-data-grid";
import {CircularProgress, Paper} from "@mui/material";
import Button from "@mui/material/Button";
import {useBackupStore} from "../../store/backup";
import {useEffect} from "react";
import {useNotificationStore} from "../../store/notification";
import IconButton from "@mui/material/IconButton";
import ReplayIcon from '@mui/icons-material/Replay';

const columns: GridColDef[] = [
    {
        field: 'id',
        headerName: 'ID',
        width: 200,
        hideable: true,
        filterable: false,
    },
    {
        field: 'name',
        headerName: 'Name',
        width: 200,
        sortable: true,
        filterable: true,
        hideable: false,
    },
    {
        field: 'identifier',
        headerName: 'Backup Type',
        width: 120,
        sortable: true,
        filterable: true,
        hideable: true,
    },
    {
        field: 'schedule',
        headerName: 'Schedule',
        width: 120,
        sortable: true,
        hideable: true,
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
            const editOnClick = () => {
                window.location.assign(`/backups/jobs/${params.id}`);
            }
            const deleteOnClick = () => {
                // TODO: Open modal here to confirm delete
                console.log(`wanting to delete job ${params.id}`);
            }
            return <div>
                <Button sx={{marginRight: "20px"}} variant="contained" color="info" onClick={editOnClick}>Edit</Button>
                <Button variant="contained" color="error" onClick={deleteOnClick}>Delete</Button>
            </div>;
        },
        valueGetter: (value, row) => `${row.id}`,
    },
];

const paginationModel = { page: 0, pageSize: 10 };

export default function BackupJobList() {
    const backupJobs = useBackupStore((state) => state.backupJobs);
    const fetchBackupJobs = useBackupStore((state) => state.fetchBackupJobs);
    const notify = useNotificationStore((state) => state.notify);
    const [fetching, setFetching] = React.useState(false);

    const getBackupJobs = () => {
        setFetching(true);
        fetchBackupJobs().then((res) => {
            if (!res.success) {
                notify({
                    message: res.message,
                    level: "error",
                    title: "Could not fetch backup jobs",
                })
            }
        });
        setFetching(false);
    }

    useEffect(() => {

        if (backupJobs.length === 0) {
            getBackupJobs();
        }
    })

    return (
        <div className="backup backup-job-list backup-job-list-wrapper">
            <div className="backup backup-job-list backup-job-list-header">
                <Typography className="backup backup-job-list backup-job-list-title" variant="h4" component="h1">Backup Job List</Typography>
                { fetching ? <CircularProgress className="backup backup-job-list backup-job-list-progress" /> : <div className="backup backup-job-list backup-job-list-progress"></div>}
                <IconButton className="backup backup-job-list backup-job-list-reload" onClick={getBackupJobs} disabled={fetching} >
                    <ReplayIcon />
                </IconButton>
            </div>
            <Paper className="users users-list users-list-table users-list-table-wrapper" elevation={2}>
                <DataGrid

                    className="users users-list users-list-table users-list-table-grid"
                    rows={backupJobs}
                    columns={columns}
                    initialState={{ pagination: { paginationModel } }}
                    pageSizeOptions={[10, 25, 50]}
                />
            </Paper>
        </div>
    );
}