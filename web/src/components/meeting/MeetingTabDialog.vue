<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent } from '@/components/ui/dialog'
import MeetingTab from '@/components/meeting/MeetingTab.vue'
import { inject } from 'vue'
import { MeetingTabComposableKey } from '@/types/meeting'

const meetingTab = inject(MeetingTabComposableKey)!
function handleClose() {
    meetingTab.state.value.isDialogOpen = false
}
</script>

<template>
    <Dialog v-model:open="meetingTab.state.value.isDialogOpen">
        <DialogContent class="flex flex-col sm:max-w-[425px] p-0 h-[70dvh]">
            <div class="flex-grow">
                <MeetingTab></MeetingTab>
            </div>
            <div class="flex-shrink-0 flex justify-end p-6">
                <Button @click="handleClose" variant="secondary">Close</Button>
            </div>
        </DialogContent>
    </Dialog>
</template>
