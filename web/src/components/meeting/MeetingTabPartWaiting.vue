<script setup lang="ts">
import { meetingQuery } from '@/queries/meeting-query'
import { computed } from 'vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import MdiAccountClockOutline from '~icons/mdi/account-clock-outline'
import { Badge } from '@/components/ui/badge'
import MeetingTabPartWaitingItem from './MeetingTabPartWaitingItem.vue'

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
            <div v-if="isSuccess" class="space-y-2">
                <MeetingTabPartWaitingItem
                    v-for="item in data"
                    :key="item.id"
                    :user="item"
                ></MeetingTabPartWaitingItem>
            </div>
            <ErrorAlert v-if="isError" :message="error?.message"></ErrorAlert>
        </div>
    </div>
</template>
