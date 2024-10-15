<script setup lang="ts">
import { getLkServerURL } from '@/lib/config'
import { Meeting, type MeetingRender } from '@/services/meeting-service'
import { useMeetingData } from '@/stores/meeting-store'
import { onMounted, shallowRef } from 'vue'

const meetingData = useMeetingData()
const renderArr = shallowRef<MeetingRender[]>([])

const meeting = new Meeting({
    token: meetingData.data.token,
    url: getLkServerURL(),
    renderArr: renderArr
})

onMounted(async () => {
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
        <div>
            <div class="circle bg-black" v-for="item in renderArr" :key="item.participantID">
                <!-- {{ item.participantID }} -->
                <video
                    v-if="item.videoElement != null"
                    class="video rounded-full"
                    :srcObject.prop="item.videoElement.srcObject"
                    autoplay
                ></video>
            </div>
        </div>
    </div>
</template>

<style scoped>
.circle {
    width: 200px;
    height: 200px;
    border-radius: 200px;
}
.video {
    width: 200px;
    height: 200px;
}
</style>
