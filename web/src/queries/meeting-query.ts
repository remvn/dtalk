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
    const meetingData = useMeetingData()
    const query = useQuery({
        queryKey: ['participants', meetingData.data.id],
        queryFn: () => {
            return meetingFetch.listParticipants({
                room_id: meetingData.data.id!
            })
        }
    })
    return query
}

export const meetingQuery = {
    useJoinRequests,
    useParticipants
}
