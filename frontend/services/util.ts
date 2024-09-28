import Cookies from 'js-cookie';
import { IUserInfo } from "../store/types";

export function setUserInfoToCookies(data: { user: IUserInfo }): string | undefined {
    return Cookies.set('userInfo', JSON.stringify(data), { secure: true });
}

export function getUserInfoFromCookies(): { user: IUserInfo } | undefined {
    const data = JSON.parse(Cookies.get('userInfo') || '{}');
    return Object.keys(data).length ? data : undefined;
}

export function checkUserInfo(callback: (currentAuth: boolean) => void) {
    const userInfo = getUserInfoFromCookies();
    if (userInfo) {
        if (userInfo?.user?.id) {
            callback(true);
        }
        callback(false);
    } else {
        callback(false);
    }
}

export function eraseCookies(): void {
    Cookies.remove('userInfo');
}
