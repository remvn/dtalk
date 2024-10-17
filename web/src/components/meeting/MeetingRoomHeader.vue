<script setup lang="ts">
import { useMeetingData } from '@/stores/meeting-store'
import { zeroPad } from '@/lib/left-pad'
import { useIntervalFn } from '@vueuse/core'
import { ref } from 'vue'

const meetingData = useMeetingData()
const timeElapsed = ref('00:00')
useIntervalFn(() => {
    if (meetingData.data.createDate == null) return
    const createDate = meetingData.data.createDate
    let seconds = (Date.now() - createDate.getTime()) / 1000
    seconds = Math.floor(seconds)
    const minutes = Math.floor(seconds / 60)
    seconds = seconds % 60
    timeElapsed.value = `${zeroPad(minutes, 2)}:${zeroPad(seconds, 2)}`
}, 1000)
</script>

<template>
    <div class="flex items-center text-lg font-mono gap-2">
        <span> {{ timeElapsed }}</span>
        <span class="text-gray-500">|</span>
        <span>
            {{ meetingData.data.name }}
        </span>
    </div>
</template>
