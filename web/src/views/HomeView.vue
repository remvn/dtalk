<script setup lang="ts">
import MeetingCreateDialog from '@/components/meeting/MeetingCreateDialog.vue'
import NavBar from '@/components/NavBar.vue'
import { Button } from '@/components/ui/button'
import UnderDevelopmentDialog from '@/components/UnderDevelopmentDialog.vue'
</script>

<template>
    <!-- <NavBar></NavBar> -->
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
                    <UnderDevelopmentDialog>
                        <Button size="lg">Broadcast</Button>
                    </UnderDevelopmentDialog>
                </div>
                <div class="pt-4 text-muted-foreground">
                    <a
                        href="https://github.com/remvn/dtalk"
                        class="underline underline-offset-4 hover:text-primary"
                        target="_blank"
                    >
                        View source on Github
                    </a>
                </div>
            </div>
        </div>
        <div class="flex justify-center items-center p-6">
            <div class="xl:w-4/6">
                <img src="../assets/connect.png" />
            </div>
        </div>
    </div>
</template>
