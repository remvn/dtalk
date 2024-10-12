import { useUserInfo } from '@/stores/user-info'

export function getAuthHeader() {
    const userInfo = useUserInfo()
    return {
        Authorization: `Bearer ${userInfo.info.token}`
    }
}

export function getBaseURL() {
    const mode = import.meta.env.MODE
    if (mode === 'production') {
        return ''
    } else {
        return 'http://localhost:8080'
    }
}

export function getURL(url: string) {
    const fullURL = new URL(url, getBaseURL()).href
    return fullURL
}

export async function requestToken(body: { name: string }) {
    const res = await fetch(getURL('/api/auth/request-token'), {
        method: 'POST',
        body: JSON.stringify(body)
    })
    if (res.status != 200) {
        return null
    }
    const json = await res.json()
    return json
}

export async function createMeeting(body: { room_name: string }) {
    const res = await fetch(getURL('/api/meeting/create'), {
        method: 'POST',
        body: JSON.stringify(body)
    })
    const json = await res.json()
    return json
}

export async function joinMeeting(body: { room_id: string }) {
    const res = await fetch(getURL('/api/meeting/join'), {
        method: 'POST',
        body: JSON.stringify(body),
        headers: getAuthHeader()
    })
    const json = await res.json()
    return json
}
