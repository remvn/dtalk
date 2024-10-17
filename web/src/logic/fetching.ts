import { getAPIBaseURL } from '@/lib/config'
import { useUserInfo } from '@/stores/user-store'
import type { HTTPError, Hooks } from 'ky'

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

export function getURL(url: string, params?: any) {
    if (params != null) {
        const paramStr = new URLSearchParams(params).toString()
        url += `?${paramStr}`
    }
    const fullURL = new URL(url, getAPIBaseURL())
    return fullURL.href
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
        return json.message as string
    } catch (e) {
        return text || res.statusText
    }
}

export async function convertKyError(error: HTTPError): Promise<HTTPError> {
    const { response } = error
    if (response && response.body) {
        error.message = await getResMessage(response)
    }
    return error
}

export const defaultKyHooks: Hooks = {
    beforeError: [convertKyError]
}
