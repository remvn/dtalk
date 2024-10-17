<script setup lang="ts">
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger
} from '@/components/ui/dialog'
import { useMeetingData } from '@/stores/meeting-store'
import { useBrowserLocation, useClipboard } from '@vueuse/core'
import { computed } from 'vue'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

const meetingData = useMeetingData()
const location = useBrowserLocation()
const inviteURL = computed(() => {
    return `${location.value.origin}/meeting/join?id=${meetingData.data.id}`
})
const { copy, copied } = useClipboard({ source: inviteURL })
</script>

<template>
    <Dialog>
        <DialogTrigger as-child>
            <slot></slot>
        </DialogTrigger>
        <DialogContent class="max-w-sm">
            <DialogHeader>
                <DialogTitle>Invite</DialogTitle>
                <DialogDescription> Use the link below to invite other people </DialogDescription>
            </DialogHeader>
            <div class="flex gap-2">
                <Input v-model="inviteURL" readonly></Input>
                <Button @click="copy(inviteURL)">
                    <span v-if="!copied">Copy</span>
                    <span v-else> Copied! </span>
                </Button>
            </div>
        </DialogContent>
    </Dialog>
</template>
