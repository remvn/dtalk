import type { QueryClient } from '@tanstack/vue-query'
import type { Room } from 'livekit-client'
import type { MeetingRenderer } from './meeting-renderer'

type Params = {
    room: Room
    queryClient: QueryClient
    renderer: MeetingRenderer
}

export class MeetingMessenger {
    room: Room
    queryClient: QueryClient
    renderer: MeetingRenderer

    constructor(params: Params) {
        this.room = params.room
        this.queryClient = params.queryClient
        this.renderer = params.renderer
    }

    handleMessage(json: any) {
        if (json.kind == null) return
        switch (json.kind) {
            case 'new_join_request': {
                break
            }
        }
    }
}
