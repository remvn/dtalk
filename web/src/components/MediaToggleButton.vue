<script setup lang="ts">
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
        :variant="props.isEnabled ? 'destructive' : 'secondary'"
    >
        <slot name="enabled" v-if="props.isEnabled"></slot>
        <slot name="disabled" v-if="!props.isEnabled"></slot>
    </Button>
</template>
