<script setup lang="ts">
import { Tabs } from '@/components/ui/tabs'
import MeetingTabPart from '@/components/meeting/MeetingTabPart.vue'
import { inject } from 'vue'
import { MeetingTabComposableKey } from '@/types/meeting'
import MeetingTabChat from '@/components/meeting/MeetingTabChat.vue'

const meetingTab = inject(MeetingTabComposableKey)!
</script>

<template>
    <Tabs v-model="meetingTab.state.value.selectedTab" class="w-full h-full">
        <MeetingTabPart></MeetingTabPart>
        <MeetingTabChat></MeetingTabChat>
    </Tabs>
</template>
