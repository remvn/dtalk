<script setup lang="ts">
import { computed } from 'vue'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { cn } from '@/lib/utils'

const { name, avalarClass, fallbackClass } = defineProps({
    name: String,
    avalarClass: String,
    fallbackClass: String
})

const char = computed(() => {
    if (name == null) return 'A'
    return name[0].toUpperCase()
})
</script>

<template>
    <Avatar :class="cn(avalarClass)">
        <AvatarFallback :class="cn('text-lg', fallbackClass)">{{ char }}</AvatarFallback>
    </Avatar>
</template>
