import Cookies from 'js-cookie';
import { IUserInfo } from "../store/types";

export function setUserInfoToCookies(data: { user: IUserInfo }): string | undefined {
    return Cookies.set('userInfo', JSON.stringify(data), { secure: true });
}

export function getUserInfoFromCookies(): { user: IUserInfo } | undefined {
    const data = JSON.parse(Cookies.get('userInfo') || '{}');
    return Object.keys(data).length ? data : undefined;
}

export function checkUserInfo(): boolean {
    const userInfo = getUserInfoFromCookies();
    if (userInfo) {
        if (userInfo?.user?.id) {
            console.log("User is authenticated");
            return true;
        }
        return false;
    } else {
        return false;
    }
}

export function eraseCookies(): void {
    Cookies.remove('userInfo');
}
