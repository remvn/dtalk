<script setup lang="ts">
import { Tabs } from '@/components/ui/tabs'
import MeetingTabParticipant from '@/components/meeting/MeetingTabParticipant.vue'
import { inject } from 'vue'
import { MeetingTabComposableKey } from '@/types/meeting'
import MeetingTabChat from '@/components/meeting/MeetingTabChat.vue'

const meetingTab = inject(MeetingTabComposableKey)!
</script>

<template>
    <Tabs v-model="meetingTab.state.value.selectedTab" class="w-full h-full">
        <MeetingTabParticipant></MeetingTabParticipant>
        <MeetingTabChat></MeetingTabChat>
    </Tabs>
</template>
