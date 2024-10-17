<script setup lang="ts">
import { meetingQuery } from '@/queries/meeting-query'
import { computed } from 'vue'
import ParticipantAvatar from '@/components/ParticipantAvatar.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import { Button } from '@/components/ui/button'
import MdiAccountClockOutline from '~icons/mdi/account-clock-outline'
import { Badge } from '@/components/ui/badge'

const { isSuccess, isError, data, error } = meetingQuery.useJoinRequests()

const count = computed(() => {
    if (data.value == null) return 0
    return data.value?.length
})
</script>

<template>
    <div class="space-y-3" v-if="count > 0">
        <div class="flex items-center">
            <MdiAccountClockOutline class="size-6"></MdiAccountClockOutline>
            <span class="font-light ml-2">Waiting</span>
            <Badge class="ml-2">{{ count }}</Badge>
        </div>
        <div>
            <div v-if="isSuccess">
                <div v-for="item in data" :key="item.id" class="flex items-center justify-between">
                    <div class="flex items-center gap-2">
                        <ParticipantAvatar :name="item.name"></ParticipantAvatar>
                        <span>{{ item.name }}</span>
                    </div>
                    <Button size="sm" variant="outline">Accept</Button>
                </div>
            </div>
            <ErrorAlert v-if="isError" :message="error?.message"></ErrorAlert>
        </div>
    </div>
</template>
