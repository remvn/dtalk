<script setup lang="ts">
import MeetingTabToggle from '@/components/meeting/MeetingTabToggleButton.vue'
import MdiAccountMultiple from '~icons/mdi/account-multiple'
import MdiAccountMultipleOutline from '~icons/mdi/account-multiple-outline'
import MdiMessageTextOutline from '~icons/mdi/message-text-outline'
import MdiMessageText from '~icons/mdi/message-text'
import { inject } from 'vue'
import { MeetingTabComposableKey } from '@/types/meeting'

const { iconClass } = defineProps({
    iconClass: String
})

const meetingTab = inject(MeetingTabComposableKey)!
</script>

<template>
    <div class="flex gap-3">
        <MeetingTabToggle name="participant" @toggle="meetingTab.toggleTab">
            <template #enabled>
                <MdiAccountMultiple :class="iconClass"></MdiAccountMultiple>
            </template>
            <template #disabled>
                <MdiAccountMultipleOutline :class="iconClass"></MdiAccountMultipleOutline>
            </template>
        </MeetingTabToggle>
        <MeetingTabToggle name="chat" @toggle="meetingTab.toggleTab">
            <template #enabled>
                <MdiMessageText :class="iconClass"></MdiMessageText>
            </template>
            <template #disabled>
                <MdiMessageTextOutline :class="iconClass"></MdiMessageTextOutline>
            </template>
        </MeetingTabToggle>
    </div>
</template>
