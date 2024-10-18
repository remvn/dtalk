import type { QueryClient } from '@tanstack/vue-query'
import type { Room } from 'livekit-client'
import type { MeetingRenderer } from '@/logic/meeting/meeting-renderer'
import { meetingQuery } from '@/queries/meeting-query'

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
                this.handleNewJoinRequest(json)
                break
            }
        }
    }

    handleNewJoinRequest(json: { total_count: number }) {
        this.queryClient.invalidateQueries({
            queryKey: [...meetingQuery.keys.joinRequest, this.room.name]
        })
    }
}
