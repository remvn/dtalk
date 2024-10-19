<script setup lang="ts">
import { MeetingTabComposableKey } from '@/types/meeting'
import { inject } from 'vue'
import MeetingTab from '@/components/meeting/MeetingTab.vue'

const meetingTab = inject(MeetingTabComposableKey)!
</script>

<template>
    <div
        v-if="meetingTab.state.value.isDrawerOpen"
        class="w-[358px] h-full pt-4 px-4 flex-shrink-0"
    >
        <MeetingTab class="border border-secondary rounded-md"></MeetingTab>
    </div>
</template>
