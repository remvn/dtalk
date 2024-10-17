import { getAPIBaseURL } from '@/config/config'
import { useUserInfo } from '@/stores/user-store'
import type { HTTPError, Hooks } from 'ky'

export function getAuthHeader() {
    const userInfo = useUserInfo()
    return {
        Authorization: `Bearer ${userInfo.data.token}`
    }
}

export function getURL(url: string) {
    const urlObj = new URL(url, getAPIBaseURL())
    return urlObj.href
}

async function getResMessage(res: Response) {
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
