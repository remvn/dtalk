<script setup lang="ts">
import { computed } from 'vue'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'

const { name } = defineProps({
    name: String
})
const char = computed(() => {
    if (name == null) return 'A'
    return name[0].toUpperCase()
})
</script>

<template>
    <Avatar>
        <AvatarFallback class="text-lg">{{ char }}</AvatarFallback>
    </Avatar>
</template>
