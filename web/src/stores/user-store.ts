import { defineStore } from 'pinia'
import { reactive } from 'vue'

type UserInfo = {
    name: string
    token: string
}

export const useUserInfo = defineStore('user-info', () => {
    const data = reactive<UserInfo>({
        name: '',
        token: ''
    })

    return { data }
})
