import ky from 'ky'
import { convertKyError, defaultHeaders, defaultKyHooks, getAuthHeader, getURL } from './fetching'
import type { User } from '@/types/user'

function create(body: { room_name: string }) {
    return fetch(getURL('/api/meeting/create'), {
        method: 'POST',
        body: JSON.stringify(body),
        headers: defaultHeaders
    })
}

function publicData(params: { room_id: string }) {
    return fetch(getURL('/api/meeting/public-data', params), {
        method: 'GET'
    })
}

function join(body: { room_id: string }) {
    return fetch(getURL('/api/meeting/join'), {
        method: 'POST',
        body: JSON.stringify(body),
        headers: getAuthHeader()
    })
}

function listParticipants(params: { room_id: string }) {
    return ky
        .get<User[]>(getURL('/api/meeting/participants'), {
            searchParams: params,
            hooks: defaultKyHooks,
            headers: getAuthHeader()
        })
        .json()
}

function listJoinRequesters(params: { room_id: string }) {
    return ky
        .get<User[]>(getURL('/api/meeting/join-requesters'), {
            searchParams: params,
            hooks: defaultKyHooks,
            headers: getAuthHeader()
        })
        .json()
}

export const meetingFetch = {
    create,
    publicData,
    join,
    listParticipants,
    listJoinRequesters
}
