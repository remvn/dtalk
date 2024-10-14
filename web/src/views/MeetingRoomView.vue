<script setup lang="ts">
import { getLkServerURL } from '@/lib/config'
import { Meeting } from '@/services/meeting-service'
import { useMeetingData } from '@/stores/meeting-store'
import { onMounted, useTemplateRef } from 'vue'

const meetingData = useMeetingData()

const meeting = new Meeting({
    token: meetingData.data.token,
    url: getLkServerURL()
})
const videoContainer = useTemplateRef('video-container')

onMounted(async () => {
    meeting.container = videoContainer.value
    meeting.setListener()
    await meeting.connect()
    const p = meeting.room.localParticipant
    // turn on the local user's camera and mic, this may trigger a browser prompt
    // to ensure permissions are granted
    await p.setCameraEnabled(true)
    await p.setMicrophoneEnabled(true)
})
</script>

<template>
    <div class="flex-grow">
        <div ref="video-container"></div>
    </div>
</template>
