import { meetingFetch } from '@/logic/meeting-fetch'
import { useMeetingData } from '@/stores/meeting-store'
import { useQuery } from '@tanstack/vue-query'

const keys = {
    joinRequest: ['meeting', 'join-requests'],
    participants: ['meeting', 'participants']
}

function useJoinRequests() {
    const meetingData = useMeetingData()
    const query = useQuery({
        queryKey: [...keys.joinRequest, meetingData.data.id],
        retry: false,
        queryFn: () => {
            return meetingFetch.listJoinRequesters({
                room_id: meetingData.data.id!
            })
        }
    })
    return query
}

function useParticipants() {
    const meetingData = useMeetingData()
    const query = useQuery({
        queryKey: [...keys.participants, meetingData.data.id],
        queryFn: () => {
            return meetingFetch.listParticipants({
                room_id: meetingData.data.id!
            })
        }
    })
    return query
}

export const meetingQuery = {
    keys,
    useJoinRequests,
    useParticipants
}
