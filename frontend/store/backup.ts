import {IBackupJob, IBackupJobCreate, IErrorDetail} from "./types";
import {create} from "zustand";
import automationBackend, {ApiStatus} from "../services/automation.backend";

type BackupState = {
    backupJobs: IBackupJob[]
    backupJobDeleteModalOpen: boolean
    backupJobDeleteId: string
    backupJobCreateModalOpen: boolean
}

type BackupActions = {
    fetchBackupJobs: () => Promise<{success: boolean, message: string, errors?: IErrorDetail[]}>
    deleteBackupJob: (id: string) => Promise<{success: boolean, message: string}>
    setBackupJobDeleteModalOpen: (open: boolean, jobId: string) => void
    createBackupJob: (job: IBackupJob) => Promise<{success: boolean, message: string}>
    setBackupJobCreateModalOpen: (open: boolean) => void
}

const defaultValues = {
    backupJobs: [],
    backupDeleteModalOpen: false,
    backupJobDeleteId: '',
    backupJobCreateModalOpen: false,
}

const useBackupStore = create<BackupState & BackupActions>((set) => ({
    backupJobs: defaultValues.backupJobs,
    backupJobDeleteModalOpen: defaultValues.backupDeleteModalOpen,
    backupJobDeleteId: defaultValues.backupJobDeleteId,
    backupJobCreateModalOpen: defaultValues.backupJobCreateModalOpen,
    fetchBackupJobs: async () => {
        const response = await automationBackend.getBackupJobs();
        if (response.status == ApiStatus.SUCCESS) {
            set({ backupJobs: response.jobs });
            return {success: true, message: response.message };
        }
        return {success: false, message: response.message };
    },
    deleteBackupJob: async (id: string) => {
        const response = await automationBackend.deleteBackupJob(id);
        if (response.status == ApiStatus.SUCCESS) {
            set((state) => ({
                ...state,
                backupJobs: state.backupJobs.filter((job) => job.id !== id),
                backupJobDeleteModalOpen: false,
                backupJobDeleteId: '',
            }));
            return {success: true, message: response.message };
        }
        return {success: false, message: response.message };
    },
    setBackupJobDeleteModalOpen: (open: boolean, jobId: string) => {
        set((state) => ({
            ...state,
            backupJobDeleteModalOpen: open,
            backupJobDeleteId: jobId,
        }));
    },
    createBackupJob: async (job: IBackupJobCreate) => {
        const response = await automationBackend.createBackupJob(job);
        if (response.status == ApiStatus.SUCCESS) {
            set((state) => ({
                ...state,
                backupJobCreateModalOpen: false,
                backupJobs: [...state.backupJobs, response.job],
            }));
            return {success: true, message: response.message };
        }
        return {success: false, message: response.message, errors: response.errors };
    },
    setBackupJobCreateModalOpen: (open: boolean) => {
        set((state) => ({
            ...state,
            backupJobCreateModalOpen: open,
        }));
    },
}));

export {useBackupStore};