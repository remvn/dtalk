import { defaultHeaders, getAuthHeader, getURL } from './fetching'

export function requestToken(body: { name: string }) {
    return fetch(getURL('/api/auth/request-token'), {
        method: 'POST',
        body: JSON.stringify(body),
        headers: defaultHeaders
    })
}

export function createMeeting(body: { room_name: string }) {
    return fetch(getURL('/api/meeting/create'), {
        method: 'POST',
        body: JSON.stringify(body),
        headers: defaultHeaders
    })
}

export function joinMeeting(body: { room_id: string }) {
    return fetch(getURL('/api/meeting/join'), {
        method: 'POST',
        body: JSON.stringify(body),
        headers: getAuthHeader()
    })
}
