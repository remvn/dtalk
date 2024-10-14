import { defineStore } from 'pinia'
import { reactive } from 'vue'

type MeetingData = {
    roomId: string
    roomName: string
    token: string
}

export const useMeetingData = defineStore('meeting_data', () => {
    const data = reactive<MeetingData>({
        roomId: '',
        roomName: '',
        token: ''
    })

    return { data }
})
