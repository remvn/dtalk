<script setup lang="ts">
import { meetingFetch } from '@/logic/meeting-fetch'
import { useMeetingData } from '@/stores/meeting-store'
import ParticipantAvatar from '../ParticipantAvatar.vue'
import type { User } from '@/types/user'
import { ref } from 'vue'
import { errorToast } from '@/logic/toast'
import { Button } from '@/components/ui/button'

const { user } = defineProps<{
    user: User
}>()

const meetingData = useMeetingData()
const isLoading = ref(false)

async function handleAccept() {
    isLoading.value = true
    try {
        meetingFetch.accept({
            room_id: meetingData.data.id!,
            accepted: true,
            requester_id: user.id
        })
    } catch (e: any) {
        errorToast(e.message)
    } finally {
        isLoading.value = false
    }
}
</script>

<template>
    <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
            <ParticipantAvatar :name="user.name"></ParticipantAvatar>
            <span>{{ user.name }}</span>
        </div>
        <Button @click="handleAccept" :disabled="isLoading" size="sm" variant="outline">
            Accept
        </Button>
    </div>
</template>
