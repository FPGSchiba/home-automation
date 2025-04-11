import {Permission, useUserStore} from "../store/user";

export enum FrontendPermissions {
    settingsRead = 'settings.read',
    settingsWrite = 'settings.write',
    usersRead = 'users.read',
    usersWrite = 'users.write',
    rolesRead = 'roles.read',
    rolesWrite = 'roles.write',
    financeRead = 'finance.read',
    financeWrite = 'finance.write',
    mealRead = 'meal.read',
    mealWrite = 'meal.write',
    backupsRead = 'backups.read',
    backupsWrite = 'backups.write',
}

const writeMethods = ['POST', 'PUT', 'DELETE', 'PATCH'];

export function hasPermission(permission: FrontendPermissions, permissions: FrontendPermissions[]): boolean {
    return  permissions.includes(permission);
}

export function remapPermissions(permissions: Permission[]): FrontendPermissions[] {
    let mappedPermissions = [];
    permissions.forEach(permission => {
        permission.routes.forEach(route => {
            if (route.path.startsWith("/")) {
                // Further group permissions by read and write
                const groupName = route.path.split('/')[1];
                switch (groupName) {
                    case 'settings':
                        if (route.methods.includes('GET')) {
                            mappedPermissions.push(FrontendPermissions.settingsRead);
                        }
                        if (route.methods.some(method => writeMethods.includes(method))) {
                            mappedPermissions.push(FrontendPermissions.settingsWrite);
                        }
                        break;
                    case 'users':
                        if (route.methods.includes('GET')) {
                            mappedPermissions.push(FrontendPermissions.usersRead);
                        }
                        if (route.methods.some(method => writeMethods.includes(method))) {
                            mappedPermissions.push(FrontendPermissions.usersWrite);
                        }
                        break;
                    case 'roles':
                        if (route.methods.includes('GET')) {
                            mappedPermissions.push(FrontendPermissions.rolesRead);
                        }
                        if (route.methods.some(method => writeMethods.includes(method))) {
                            mappedPermissions.push(FrontendPermissions.rolesWrite);
                        }
                        break;
                    case 'finance':
                        if (route.methods.includes('GET')) {
                            mappedPermissions.push(FrontendPermissions.financeRead);
                        }
                        if (route.methods.some(method => writeMethods.includes(method))) {
                            mappedPermissions.push(FrontendPermissions.financeWrite);
                        }
                        break;
                    case 'meal':
                        if (route.methods.includes('GET')) {
                            mappedPermissions.push(FrontendPermissions.mealRead);
                        }
                        if (route.methods.some(method => writeMethods.includes(method))) {
                            mappedPermissions.push(FrontendPermissions.mealWrite);
                        }
                        break;
                    case 'backups':
                        if (route.methods.includes('GET')) {
                            mappedPermissions.push(FrontendPermissions.backupsRead);
                        }
                        if (route.methods.some(method => writeMethods.includes(method))) {
                            mappedPermissions.push(FrontendPermissions.backupsWrite);
                        }
                        break;
                    default:
                        console.error(`Unknown permission group: ${groupName}`);
                        return [];
                }

            }
            else {
                mappedPermissions = [ // All Permissions | Should be improved
                    FrontendPermissions.settingsRead,
                    FrontendPermissions.settingsWrite,
                    FrontendPermissions.usersRead,
                    FrontendPermissions.usersWrite,
                    FrontendPermissions.rolesRead,
                    FrontendPermissions.rolesWrite,
                    FrontendPermissions.financeRead,
                    FrontendPermissions.financeWrite,
                    FrontendPermissions.mealRead,
                    FrontendPermissions.mealWrite,
                    FrontendPermissions.backupsRead,
                    FrontendPermissions.backupsWrite
                ]
            }
        })
    })
    return mappedPermissions;
}