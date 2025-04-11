import {IBackupJob} from "./types";
import {create} from "zustand";
import automationBackend, {ApiStatus} from "../services/automation.backend";

type BackupState = {
    backupJobs: IBackupJob[]
}

type BackupActions = {
    fetchBackupJobs: () => Promise<{success: boolean, message: string}>
}

const defaultValues = {
    backupJobs: []
}

const useBackupStore = create<BackupState & BackupActions>((set) => ({
    backupJobs: defaultValues.backupJobs,
    fetchBackupJobs: async () => {
        const response = await automationBackend.getBackupJobs();
        if (response.status == ApiStatus.SUCCESS) {
            set({ backupJobs: response.jobs });
            return {success: true, message: response.message };
        }
        return {success: false, message: response.message };
    },
}));

export {useBackupStore};