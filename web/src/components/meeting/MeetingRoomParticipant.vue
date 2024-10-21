<script setup lang="ts">
import { cn } from '@/lib/utils'
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
            muted
            playsinline
        ></video>
        <div
            v-else
            :class="
                cn(
                    'char-container w-[50%] max-w-[6rem] md:max-w-[12rem]',
                    'rounded-full bg-secondary aspect-square flex items-center justify-center'
                )
            "
        >
            <span class="text-[40cqw]">{{ item.name[0].toUpperCase() }}</span>
        </div>
        <audio v-if="item.audioSrc != null" :srcObject.prop="item.audioSrc" autoplay></audio>
    </div>
</template>

<style scoped>
.video-container {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}

.video {
    width: auto;
    height: 100%;
}

.char-container {
    container-type: inline-size;
}
</style>
