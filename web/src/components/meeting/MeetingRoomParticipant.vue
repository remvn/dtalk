<script setup lang="ts">
import type { MeetingRender } from '@/logic/meeting/meeting-renderer'

const { item } = defineProps<{
    item: MeetingRender
}>()
</script>

<template>
    <div class="video-container">
        <video
            v-if="item.videoSrc != null"
            class="video rounded-md"
            :srcObject.prop="item.videoSrc"
            autoplay
        ></video>
        <div
            v-else
            class="w-[50%] max-w-[12rem] rounded-full bg-secondary aspect-square flex items-center justify-center"
        >
            <span class="char font-light">{{ item.name[0].toUpperCase() }}</span>
        </div>
        <audio v-if="item.audioSrc != null" :srcObject.prop="item.audioSrc" autoplay></audio>
    </div>
</template>

<style scoped>
.char {
    font-size: 7cqw;
    line-height: 7cqw;
}

.video-container {
    width: 100%;
    height: 100%;
    max-width: 100%;
    max-height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow-y: hidden;
}

.video {
    width: auto;
    height: 100%;
}
</style>
