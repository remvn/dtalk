import { defineStore } from 'pinia'
import { reactive } from 'vue'

type UserInfo = {
    name: string
    token: string
}

export const useUserInfo = defineStore('user_info', () => {
    const info = reactive<UserInfo>({
        name: '',
        token: ''
    })

    function setInfo(newInfo: UserInfo) {
        info.name = newInfo.name
        info.token = newInfo.token
    }

    return { info, setInfo }
})
