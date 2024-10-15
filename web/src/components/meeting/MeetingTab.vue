<script setup lang="ts">
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import MeetingTabParticipant from '@/components/meeting/MeetingTabParticipant.vue'
import { inject, type Ref } from 'vue'
import { MeetingTabStateKey, type MeetingTabState } from '@/types/meeting'

const tabState = inject(MeetingTabStateKey) as Ref<MeetingTabState>
</script>

<template>
    <Tabs v-model="tabState.selectedTab" class="w-[400px]">
        <TabsList class="grid w-full grid-cols-2">
            <TabsTrigger value="participant"> Participants </TabsTrigger>
            <TabsTrigger value="chat"> Chat </TabsTrigger>
        </TabsList>
        <MeetingTabParticipant></MeetingTabParticipant>
    </Tabs>
</template>
