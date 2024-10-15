<script setup lang="ts">
import { getLkServerURL } from '@/lib/config'
import { Meeting, type MeetingRender } from '@/services/meeting-service'
import { useMeetingData } from '@/stores/meeting-store'
import { onBeforeUnmount, onMounted, ref, shallowRef } from 'vue'
import { breakpointsTailwind, useBreakpoints, useThrottleFn } from '@vueuse/core'
import HeroiconsVideoCamera from '~icons/heroicons/video-camera'
import HeroiconsMicrophone from '~icons/heroicons/microphone'
import MdiPhoneHangup from '~icons/mdi/phone-hangup'
import { Button } from '@/components/ui/button'

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

const isCameraEnabled = ref(false)
async function handleCameraToggle() {
    const participant = meeting.room.localParticipant
    await participant.setCameraEnabled(participant.isCameraEnabled)
    isCameraEnabled.value = participant.isCameraEnabled
}

onMounted(async () => {
    meeting.setListener()
    await meeting.connect()

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
            <div class="row-span-1 flex justify-between items-center px-6 py-4">
                <span class="text-lg">{{ meetingData.data.roomName }}</span>
                <div class="flex gap-3">
                    <Button size="icon" class="size-12 rounded-full" variant="secondary">
                        <HeroiconsMicrophone class="size-6"></HeroiconsMicrophone>
                    </Button>
                    <Button size="icon" class="size-12 rounded-full" variant="secondary">
                        <HeroiconsVideoCamera class="size-6"></HeroiconsVideoCamera>
                    </Button>
                    <Button size="icon" class="h-12 w-20 rounded-full" variant="destructive">
                        <MdiPhoneHangup class="size-6"></MdiPhoneHangup>
                    </Button>
                </div>
                <span class="text-lg">Settings</span>
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
