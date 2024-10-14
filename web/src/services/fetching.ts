import { getAPIBaseURL } from '@/lib/config'
import { useUserInfo } from '@/stores/user-store'

export const defaultHeaders = {
    'Content-Type': 'application/json'
}

export function getAuthHeader(headers = defaultHeaders) {
    const userInfo = useUserInfo()
    return {
        ...headers,
        Authorization: `Bearer ${userInfo.info.token}`
    }
}

export function getURL(url: string) {
    const fullURL = new URL(url, getAPIBaseURL()).href
    return fullURL
}

export async function getJSON(promise: Promise<Response>) {
    const res = await promise
    const json = await res.json()
    return json
}

export async function getResMessage(res: Response) {
    let text = ''
    try {
        text = await res.text()
        const json = JSON.parse(text)
        if (typeof json.message !== 'string' || json.message == null) {
            throw new Error('invalid or missing message')
        }
        return json.message
    } catch (e) {
        return text || res.statusText
    }
}
