<script setup lang="ts">
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
import { useQueryClient } from '@tanstack/vue-query'
import { Meeting } from '@/logic/meeting/meeting-service'
import type { MeetingRenderMap } from '@/logic/meeting/meeting-renderer'
import { useRouter } from 'vue-router'
import { useUserInfo } from '@/stores/user-store'
import MeetingRoomParticipant from '@/components/meeting/MeetingRoomParticipant.vue'
import OverlayScroll from '@/components/OverlayScroll.vue'
import { cn } from '@/lib/utils'
import MeetingTabDialog from '@/components/meeting/MeetingTabDialog.vue'
import MeetingTabDrawer from '@/components/meeting/MeetingTabDrawer.vue'

const router = useRouter()
const userInfo = useUserInfo()
const meetingData = useMeetingData()
const renderMap = shallowRef<MeetingRenderMap>(new Map())
const queryClient = useQueryClient()

const meeting = new Meeting({
    token: meetingData.data.token!,
    url: meetingData.data.wsUrl!,
    renderMap: renderMap,
    queryClient: queryClient,
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
}

const isMicroEnabled = ref(false)
async function handleMicroToggle() {
    const participant = meeting.room.localParticipant
    await participant.setMicrophoneEnabled(!participant.isMicrophoneEnabled)
    isMicroEnabled.value = participant.isMicrophoneEnabled
}

async function handleDisconnect() {
    meeting.room.disconnect()
    console.log(`resetting states`)
    meetingData.$reset()
    userInfo.$reset()
    router.push('/')
}

const meetingTab = useMeetingTab()
provide(MeetingTabComposableKey, meetingTab)

onMounted(async () => {
    if (meetingData.data.id == null) {
        console.log('meeting data is undefined')
        router.push('/')
        return
    }
    meeting.setListener()
    try {
        await meeting.connect()
    } catch (e) {
        console.log('unable to connect')
        console.log(e)
        router.push('/')
    }

    meeting.renderer.renderGrid()

    window.addEventListener('resize', handleWindowResize)
})

onBeforeUnmount(() => {
    window.removeEventListener('resize', handleWindowResize)
})

const iconClass = 'size-6'
</script>

<template>
    <!-- These css below drove me crazy. need to add basis-0  -->
    <!-- and min-h-0 to prevent it overflow and stay within h-full -->
    <!-- The same also apply to overlayscroll plugin  -->
    <div class="flex-grow basis-0 min-h-0">
        <div class="grid grid-rows-10 h-full">
            <div class="flex row-span-9">
                <div class="flex-grow grid video-grid lg:px-4 lg:pt-4 gap-2 lg:gap-4">
                    <MeetingRoomParticipant
                        v-for="item in renderMap.values()"
                        :key="item.participantID"
                        :item="item"
                    >
                    </MeetingRoomParticipant>
                </div>
                <MeetingTabDrawer></MeetingTabDrawer>
            </div>
            <!-- horizontal scroll for toggler bar  -->
            <OverlayScroll class="row-span-1">
                <div :class="cn('flex gap-6 items-center h-full px-6', 'justify-between')">
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
                        <MediaToggleButton
                            @toggle="handleCameraToggle"
                            :is-enabled="isCameraEnabled"
                        >
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
            </OverlayScroll>
        </div>
        <MeetingTabDialog></MeetingTabDialog>
    </div>
</template>

<style scoped>
.video-grid {
    grid-template-columns: repeat(v-bind('gridSize.col'), minmax(0, 1fr));
    grid-template-rows: repeat(v-bind('gridSize.row'), minmax(0, 1fr));
}
</style>
