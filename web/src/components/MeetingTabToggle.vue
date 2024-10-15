<script setup lang="ts">
import { Button } from './ui/button'

const props = defineProps({
    isEnabled: Boolean
})

const emit = defineEmits(['toggle'])

function handleToggle() {
    emit('toggle')
}
</script>

<template>
    <Button
        @click="handleToggle"
        size="icon"
        class="size-12 rounded-full"
        :variant="props.isEnabled ? 'outline' : 'ghost'"
    >
        <slot name="enabled" v-if="props.isEnabled"></slot>
        <slot name="disabled" v-if="!props.isEnabled"></slot>
    </Button>
</template>
