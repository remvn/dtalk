import { defineStore } from 'pinia'
import { reactive } from 'vue'

type UserInfo = {
    name?: string
    token?: string
}

export const useUserInfo = defineStore('user-info', () => {
    const data = reactive<UserInfo>({
        name: undefined,
        token: undefined
    })

    function $reset() {
        data.name = undefined
        data.token = undefined
    }

    return { data, $reset }
})
