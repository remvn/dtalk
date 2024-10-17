<script setup lang="ts">
import { getLkServerURL } from '@/config/config'
import { Meeting, type MeetingRenderMap } from '@/logic/meeting-service'
import { useMeetingData } from '@/stores/meeting-store'
import { onBeforeUnmount, onMounted, provide, ref, shallowRef } from 'vue'
import { useThrottleFn } from '@vueuse/core'
import MdiPhoneHangup from '~icons/mdi/phone-hangup'
import { Button } from '@/components/ui/button'
import MediaToggleButton from '@/components/MediaToggleButton.vue'
import MdiMicrophoneOff from '~icons/mdi/microphone-off'
import MdiMicrophoneOutline from '~icons/mdi/microphone-outline'
import MdiCameraOff from '~icons/mdi/camera-off'
import MdiCameraOutline from '~icons/mdi/camera-outline'
import MeetingTab from '@/components/meeting/MeetingTab.vue'
import { MeetingTabComposableKey } from '@/types/meeting'
import { useTwBreakpoints } from '@/hooks/use-tw-breakpoints'
import { useMeetingTab } from '@/hooks/use-meeting-tab'
import MeetingTabToggleBar from '@/components/meeting/MeetingTabToggleBar.vue'
import MeetingRoomHeader from '@/components/meeting/MeetingRoomHeader.vue'

const meetingData = useMeetingData()
const renderMap = shallowRef<MeetingRenderMap>(new Map())

const meeting = new Meeting({
    token: meetingData.data.token!,
    url: getLkServerURL(),
    renderMap: renderMap,
    setGridSize: setGridSize
})

const breakpoints = useTwBreakpoints()
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

const meetingTab = useMeetingTab()
provide(MeetingTabComposableKey, meetingTab)

onMounted(async () => {
    meeting.setListener()
    await meeting.connect()

    meeting.rerenderGrid()

    window.addEventListener('resize', handleWindowResize)
})

onBeforeUnmount(() => {
    window.removeEventListener('resize', handleWindowResize)
})

const iconClass = 'size-6'
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
                <div
                    v-if="meetingTab.state.value.isDrawerOpen"
                    class="w-[358px] h-full p-4 flex-shrink-0"
                >
                    <MeetingTab></MeetingTab>
                </div>
            </div>
            <div class="row-span-1 flex justify-between items-center px-6 py-4">
                <MeetingRoomHeader></MeetingRoomHeader>
                <div class="flex gap-3">
                    <MediaToggleButton @toggle="handleMicroToggle" :is-enabled="isMicroEnabled">
                        <template #disabled>
                            <MdiMicrophoneOutline :class="iconClass"></MdiMicrophoneOutline>
                        </template>
                        <template #enabled>
                            <MdiMicrophoneOff :class="iconClass"></MdiMicrophoneOff>
                        </template>
                    </MediaToggleButton>
                    <MediaToggleButton @toggle="handleCameraToggle" :is-enabled="isCameraEnabled">
                        <template #disabled>
                            <MdiCameraOutline :class="iconClass"></MdiCameraOutline>
                        </template>
                        <template #enabled>
                            <MdiCameraOff :class="iconClass"></MdiCameraOff>
                        </template>
                    </MediaToggleButton>
                    <Button
                        @click="handleDisconnect"
                        size="icon"
                        class="h-12 w-20 rounded-full"
                        variant="destructive"
                    >
                        <MdiPhoneHangup :class="iconClass"></MdiPhoneHangup>
                    </Button>
                </div>
                <MeetingTabToggleBar :icon-class="iconClass"></MeetingTabToggleBar>
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
