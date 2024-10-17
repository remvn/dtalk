<script setup lang="ts">
import { meetingQuery } from '@/queries/meeting-query'
import ParticipantAvatar from '@/components/ParticipantAvatar.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import MdiAccountSupervisorOutline from '~icons/mdi/account-supervisor-outline'
import { computed } from 'vue'
import { Badge } from '@/components/ui/badge'

const { data, isSuccess, isError, error } = meetingQuery.useParticipants()
const count = computed(() => {
    if (data.value == null) return 0
    return data.value.length
})
</script>

<template>
    <div class="space-y-3">
        <div class="flex items-center">
            <MdiAccountSupervisorOutline class="size-6"></MdiAccountSupervisorOutline>
            <span class="font-light ml-2"> Attending </span>
            <Badge class="ml-2">{{ count }}</Badge>
        </div>
        <div v-if="isSuccess" class="space-y-2">
            <div v-for="item in data" :key="item.id" class="flex items-center gap-2">
                <ParticipantAvatar :name="item.name"></ParticipantAvatar>
                <span>{{ item.name }}</span>
            </div>
        </div>
        <ErrorAlert v-if="isError" :message="error?.message"></ErrorAlert>
    </div>
</template>
