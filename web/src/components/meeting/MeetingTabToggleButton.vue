<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { MeetingTabComposableKey } from '@/types/meeting'
import { computed, inject } from 'vue'

const props = defineProps({
    name: {
        type: String,
        required: true
    }
})

const meetingTab = inject(MeetingTabComposableKey)!
const state = meetingTab.state

function handleToggle() {
    meetingTab.toggleTab(props.name)
}

const isOpened = computed(() => {
    return (
        (state.value.isDialogOpen || state.value.isDrawerOpen) &&
        state.value.selectedTab === props.name
    )
})
</script>

<template>
    <Button
        @click="handleToggle"
        size="icon"
        class="size-12 rounded-full"
        :variant="isOpened ? 'outline' : 'ghost'"
    >
        <slot name="enabled" v-if="isOpened"></slot>
        <slot name="disabled" v-if="!isOpened"></slot>
    </Button>
</template>
