<script setup lang="ts">
import { Tabs } from '@/components/ui/tabs'
import MeetingTabParticipant from '@/components/meeting/MeetingTabParticipant.vue'
import { inject } from 'vue'
import { MeetingTabComposableKey } from '@/types/meeting'

const meetingTab = inject(MeetingTabComposableKey)!
</script>

<template>
    <Tabs v-model="meetingTab.state.value.selectedTab" class="w-full h-full">
        <MeetingTabParticipant></MeetingTabParticipant>
    </Tabs>
</template>
