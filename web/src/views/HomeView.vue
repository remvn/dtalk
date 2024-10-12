<script setup lang="ts">
import MeetingCreateDialog from '@/components/MeetingCreateDialog.vue'
import { Button } from '@/components/ui/button'
</script>

<template>
    <div>TODO: logo, navbar here</div>
    <div class="flex-grow grid grid-cols-1 md:grid-cols-2">
        <div class="flex justify-center items-center p-6">
            <div class="xl:w-4/6">
                <h1 class="scroll-m-20 text-4xl tracking-tight lg:text-5xl">
                    Meeting and live-streaming for everyone 🎉
                </h1>
                <h2 class="text-xl lg:text-xl pt-4">
                    Connect, live-stream everywhere instantly without an account
                </h2>
                <div class="flex gap-3 pt-4">
                    <MeetingCreateDialog></MeetingCreateDialog>
                    <Button size="lg">Broadcast</Button>
                </div>
            </div>
        </div>
        <div class="flex justify-center items-center p-6">
            <div class="xl:w-4/6">
                <span>TODO: images, info here</span>
            </div>
        </div>
    </div>
</template>
