import { meetingFetch } from '@/logic/meeting-fetch'
import { useMeetingData } from '@/stores/meeting-store'
import { useQuery } from '@tanstack/vue-query'

function useJoinRequests() {
    const data = useMeetingData()
    const query = useQuery({
        queryKey: ['join-requesters', data.data.id],
        retry: false,
        queryFn: () => {
            return meetingFetch.listJoinRequesters({
                room_id: data.data.id!
            })
        }
    })
    return query
}

function useParticipants() {
    const data = useMeetingData()
    const query = useQuery({
        queryKey: ['participants', data.data.id],
        queryFn: () => {
            return meetingFetch.listParticipants({
                room_id: data.data.id!
            })
        }
    })
    return query
}

export const meetingQuery = {
    useJoinRequests,
    useParticipants
}
