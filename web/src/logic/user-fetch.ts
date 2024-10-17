import { defaultHeaders, getURL } from './fetching'

export function requestToken(body: { name: string }) {
    return fetch(getURL('/api/auth/request-token'), {
        method: 'POST',
        body: JSON.stringify(body),
        headers: defaultHeaders
    })
}
