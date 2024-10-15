<script setup lang="ts">
import { getLkServerURL } from '@/lib/config'
import { Meeting, type MeetingRender } from '@/services/meeting-service'
import { useMeetingData } from '@/stores/meeting-store'
import { onBeforeUnmount, onMounted, ref, shallowRef } from 'vue'
import { breakpointsTailwind, useBreakpoints, useThrottleFn } from '@vueuse/core'

const meetingData = useMeetingData()
const renderArr = shallowRef<MeetingRender[]>([])

const meeting = new Meeting({
    token: meetingData.data.token,
    url: getLkServerURL(),
    renderArr: renderArr,
    setGridSize: setGridSize
})

const breakpoints = useBreakpoints(breakpointsTailwind)
const gridSize = ref({
    col: 1,
    row: 1
})

const handleWindowResize = useThrottleFn(() => {
    setGridSize(meeting.room.numParticipants)
}, 100)

function setGridSize(numParticipants: number) {
    let row = 1
    let col = 1
    let maxCol = 2
    if (breakpoints.isGreaterOrEqual('lg')) maxCol = 3
    if (numParticipants >= 2) col = Math.min(maxCol, 2)
    if (numParticipants >= 5) col = Math.min(maxCol, 3)

    row = Math.ceil((numParticipants * 1.0) / col)
    gridSize.value = {
        row: row,
        col: col
    }
}

onMounted(async () => {
    meeting.setListener()
    await meeting.connect()
    const p = meeting.room.localParticipant
    // turn on the local user's camera and mic, this may trigger a browser prompt
    // to ensure permissions are granted
    await p.setCameraEnabled(true)
    await p.setMicrophoneEnabled(true)

    window.addEventListener('resize', handleWindowResize)
})

onBeforeUnmount(() => {
    window.removeEventListener('resize', handleWindowResize)
})
</script>

<template>
    <div class="flex-grow">
        <div class="grid grid-rows-10 h-full">
            <div class="row-span-9 flex">
                <div class="flex-grow h-full grid video-grid">
                    <div
                        class="video-container"
                        v-for="item in renderArr"
                        :key="item.participantID"
                    >
                        <video
                            v-if="item.videoElement != null"
                            class="video"
                            :srcObject.prop="item.videoElement.srcObject"
                            autoplay
                        ></video>
                    </div>
                </div>
                <div class="w-[300px]"></div>
            </div>
            <div class="row-span-1 flex justify-center items-center">
                <span class="text-lg">TODO Tools Bar</span>
            </div>
        </div>
    </div>
</template>

<style scoped>
.video-grid {
    grid-template-columns: repeat(v-bind('gridSize.col'), minmax(0, 1fr));
    grid-template-rows: repeat(v-bind('gridSize.row'), minmax(0, 1fr));
}
.video-container {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}
.video {
    width: auto;
    height: 100%;
}
</style>
