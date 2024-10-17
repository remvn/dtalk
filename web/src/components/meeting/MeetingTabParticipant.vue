<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle
} from '@/components/ui/card'
import { TabsContent } from '@/components/ui/tabs'
import { meetingFetch } from '@/logic/meeting-fetch'
import { useMeetingData } from '@/stores/meeting-store'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import ErrorAlert from '@/components/ErrorAlert.vue'

const meetingData = useMeetingData()

const requesterQuery = useQuery({
    queryKey: ['join-requesters', meetingData.data.roomId],
    retry: false,
    queryFn: () => {
        return meetingFetch.listJoinRequesters({
            room_id: meetingData.data.roomId!
        })
    }
})

const joinRequestCount = computed(() => {
    if (requesterQuery.data.value == null) return 0
    return requesterQuery.data.value.length
})

const participantsQuery = useQuery({
    queryKey: ['participants', meetingData.data.roomId],
    queryFn: () => {
        return meetingFetch.listParticipants({
            room_id: meetingData.data.roomId!
        })
    }
})
</script>

<template>
    <TabsContent value="participant" class="h-full">
        <Card class="h-full">
            <CardHeader>
                <CardTitle class="text-lg">Participants</CardTitle>
                <!-- <CardDescription> -->
                <!--     Make changes to your account here. Click save when you're done. -->
                <!-- </CardDescription> -->
            </CardHeader>
            <CardContent class="space-y-3">
                <div class="space-y-3" v-if="joinRequestCount > 0">
                    <span>Waiting to join</span>
                    <div>
                        <div v-if="requesterQuery.isSuccess.value">
                            <div
                                v-for="item in requesterQuery.data.value"
                                :key="item.id"
                                class="flex items-center justify-between"
                            >
                                <div class="flex items-center gap-2">
                                    <Avatar>
                                        <AvatarFallback>{{ item.name[0] }}</AvatarFallback>
                                    </Avatar>
                                    <span>{{ item.name }}</span>
                                </div>
                                <Button size="sm" variant="outline">Accept</Button>
                            </div>
                        </div>
                        <ErrorAlert
                            v-if="requesterQuery.status.value === 'error'"
                            :message="requesterQuery.error.value?.message"
                        ></ErrorAlert>
                    </div>
                </div>
                <div class="space-y-3">
                    <span>Attending</span>
                    <div v-if="participantsQuery.isSuccess.value">
                        <div
                            v-for="item in participantsQuery.data.value"
                            :key="item.id"
                            class="flex items-center gap-2"
                        >
                            <Avatar>
                                <AvatarFallback>{{ item.name[0] }}</AvatarFallback>
                            </Avatar>
                            <span>{{ item.name }}</span>
                        </div>
                    </div>
                    <ErrorAlert
                        v-if="participantsQuery.status.value === 'error'"
                        :message="participantsQuery.error.value?.message"
                    ></ErrorAlert>
                </div>
            </CardContent>
        </Card>
    </TabsContent>
</template>
