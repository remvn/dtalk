import ky from 'ky'
import { defaultKyHooks, getAuthHeader, getURL } from './fetching'
import type { User } from '@/types/user'

function create(body: { room_name: string }) {
    type Res = {
        room_id: string
    }
    return ky
        .post<Res>(getURL('/api/meeting/create'), {
            json: body
        })
        .json()
}

function publicData(params: { room_id: string }) {
    type Res = {
        name: string
    }
    return ky
        .get<Res>(getURL('/api/meeting/public-data'), {
            searchParams: params,
            hooks: defaultKyHooks
        })
        .json()
}

function join(body: { room_id: string }) {
    type Res = {
        id: string
        name: string
        token: string
        create_date: Date
    }
    return ky
        .post<Res>(getURL('/api/meeting/join'), {
            json: body,
            hooks: defaultKyHooks,
            headers: getAuthHeader()
        })
        .json()
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
