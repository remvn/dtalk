<script setup lang="ts">
import { ExclamationTriangleIcon } from '@radix-icons/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const props = defineProps({
    message: String
})
</script>

<template>
    <Alert v-if="props.message !== ''" variant="destructive">
        <ExclamationTriangleIcon class="w-4 h-4" />
        <AlertTitle>An error happened</AlertTitle>
        <AlertDescription>
            {{ props.message }}
        </AlertDescription>
    </Alert>
</template>
