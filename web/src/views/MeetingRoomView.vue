<script setup lang="ts">
import { getLkServerURL } from '@/lib/config'
import { Meeting, type MeetingRenderMap } from '@/services/meeting-service'
import { useMeetingData } from '@/stores/meeting-store'
import { onBeforeUnmount, onMounted, provide, ref, shallowRef } from 'vue'
import { breakpointsTailwind, useBreakpoints, useThrottleFn } from '@vueuse/core'
import MdiPhoneHangup from '~icons/mdi/phone-hangup'
import { Button } from '@/components/ui/button'
import MediaToggleButton from '@/components/MediaToggleButton.vue'
import MdiMicrophoneOff from '~icons/mdi/microphone-off'
import MdiMicrophoneOutline from '~icons/mdi/microphone-outline'
import MdiCameraOff from '~icons/mdi/camera-off'
import MdiCameraOutline from '~icons/mdi/camera-outline'
import MeetingTab from '@/components/MeetingTab.vue'
import { type MeetingTabState, MeetingTabStateKey } from '@/types/meeting'

const meetingData = useMeetingData()
const renderMap = shallowRef<MeetingRenderMap>(new Map())

const meeting = new Meeting({
    token: meetingData.data.token,
    url: getLkServerURL(),
    renderMap: renderMap,
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
    await participant.setCameraEnabled(!participant.isCameraEnabled)
    isCameraEnabled.value = participant.isCameraEnabled
    meeting.rerenderGrid()
}

const isMicroEnabled = ref(false)
async function handleMicroToggle() {
    const participant = meeting.room.localParticipant
    await participant.setMicrophoneEnabled(!participant.isMicrophoneEnabled)
    isMicroEnabled.value = participant.isMicrophoneEnabled
}

async function handleDisconnect() {}

const tabState = ref<MeetingTabState>({
    isDrawerOpen: false,
    isDialogOpen: false,
    selectedTab: 'participant'
})
provide(MeetingTabStateKey, tabState)

onMounted(async () => {
    meeting.setListener()
    await meeting.connect()

    meeting.rerenderGrid()

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
                        v-for="item in renderMap.values()"
                        :key="item.participantID"
                    >
                        <video
                            v-if="item.videoSrc != null"
                            class="video"
                            :srcObject.prop="item.videoSrc"
                            autoplay
                        ></video>
                        <audio
                            v-if="item.audioSrc != null"
                            :srcObject.prop="item.audioSrc"
                            autoplay
                        ></audio>
                    </div>
                </div>
                <div v-if="tabState.isDrawerOpen" class="w-[400px]">
                    <MeetingTab v-model="tabState.selectedTab"></MeetingTab>
                </div>
            </div>
            <div class="row-span-1 flex justify-between items-center px-6 py-4">
                <span class="text-lg">{{ meetingData.data.roomName }}</span>
                <div class="flex gap-3">
                    <MediaToggleButton @toggle="handleMicroToggle" :is-enabled="isMicroEnabled">
                        <template #disabled>
                            <MdiMicrophoneOutline class="size-6"></MdiMicrophoneOutline>
                        </template>
                        <template #enabled>
                            <MdiMicrophoneOff class="size-6"></MdiMicrophoneOff>
                        </template>
                    </MediaToggleButton>
                    <MediaToggleButton @toggle="handleCameraToggle" :is-enabled="isCameraEnabled">
                        <template #disabled>
                            <MdiCameraOutline class="size-6"></MdiCameraOutline>
                        </template>
                        <template #enabled>
                            <MdiCameraOff class="size-6"></MdiCameraOff>
                        </template>
                    </MediaToggleButton>
                    <Button
                        @click="handleDisconnect"
                        size="icon"
                        class="h-12 w-20 rounded-full"
                        variant="destructive"
                    >
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
