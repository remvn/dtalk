import ky from 'ky'
import { defaultKyHooks, getURL } from './fetching'

function getLivekitClientURL() {
    type Res = {
        url: string
    }
    return ky
        .get<Res>(getURL('/api/public/livekit-client-url'), {
            hooks: defaultKyHooks
        })
        .json()
}

export const publicFetch = {
    getLivekitClientURL
}
