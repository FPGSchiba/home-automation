import Cookies from 'js-cookie';
import { IUser } from "../store/types";

export function setUserInfoToCookies(data: { user: IUser, token: string }): string | undefined {
    return Cookies.set('userInfo', JSON.stringify(data), { secure: true });
}

export function getUserInfoFromCookies(): { user: IUser, token: string } | undefined {
    const data = JSON.parse(Cookies.get('userInfo') || '{}');
    return Object.keys(data).length ? data : undefined;
}

export function checkUserInfo(): boolean {
    const userInfo = getUserInfoFromCookies();
    if (userInfo) {
        if (userInfo?.user?.id && userInfo?.token) {
            return isTokenExpired(userInfo.token);
        }
        return false;
    } else {
        return false;
    }
}

function prettyFyDate(timeDifferenceInMs: number): string {
    const seconds = Math.floor((timeDifferenceInMs / 1000) % 60);
    const minutes = Math.floor((timeDifferenceInMs / (1000 * 60)) % 60);
    const hours = Math.floor((timeDifferenceInMs / (1000 * 60 * 60)) % 24);
    const days = Math.floor(timeDifferenceInMs / (1000 * 60 * 60 * 24));

    // Construct a relative time string
    let relativeTime = '';
    if (days > 0) {
        relativeTime += `${days} day${days > 1 ? 's' : ''} `;
    }
    if (hours > 0) {
        relativeTime += `${hours} hour${hours > 1 ? 's' : ''} `;
    }
    if (minutes > 0) {
        relativeTime += `${minutes} minute${minutes > 1 ? 's' : ''} `;
    }
    if (seconds > 0 || relativeTime === '') {
        relativeTime += `${seconds} second${seconds > 1 ? 's' : ''}`;
    }
    return relativeTime;
}

function isTokenExpired(token: string): boolean {
    const tokenData = JSON.parse(atob(token.split('.')[1]));
    console.log("Token expires in: ", prettyFyDate(tokenData.exp * 1000 - Date.now()));
    return tokenData.exp * 1000 > Date.now();
}

export function eraseCookies(): void {
    Cookies.remove('userInfo');
}
