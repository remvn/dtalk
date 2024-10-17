import ky from 'ky'
import { defaultKyHooks, getURL } from './fetching'

function requestToken(body: { name: string }) {
    type Res = {
        access_token: string
    }
    return ky
        .post<Res>(getURL('/api/auth/request-token'), {
            json: body,
            hooks: defaultKyHooks
        })
        .json()
}

export const userFetch = {
    requestToken
}
