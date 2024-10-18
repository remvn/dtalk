import { defineStore } from 'pinia'
import { reactive } from 'vue'

type MeetingData = {
    id?: string
    name?: string
    token?: string
    createDate?: Date
}

export const useMeetingData = defineStore('meeting_data', () => {
    const data = reactive<MeetingData>({})

    function $reset() {
        data.id = undefined
        data.name = undefined
        data.token = undefined
        data.createDate = undefined
    }

    return { data, $reset }
})
