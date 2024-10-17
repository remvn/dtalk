<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { TabsContent } from '@/components/ui/tabs'
import MeetingTabPartWaiting from '@/components/meeting/MeetingTabPartWaiting.vue'
import MeetingTabPartAttending from '@/components/meeting/MeetingTabPartAttending.vue'
import MeetingInviteDialog from './MeetingInviteDialog.vue'
import { Button } from '@/components/ui/button'
import MdiAccountPlusOutline from '~icons/mdi/account-plus-outline'
import OverlayScroll from '@/components/OverlayScroll.vue'
</script>

<template>
    <TabsContent value="participant" class="h-full">
        <Card class="h-full flex flex-col overflow-hidden">
            <CardHeader class="flex-shrink-0">
                <CardTitle class="text-lg">Participants</CardTitle>
                <CardDescription>
                    <MeetingInviteDialog>
                        <Button>
                            <MdiAccountPlusOutline class="size-4 mr-2"></MdiAccountPlusOutline>
                            Invite
                        </Button>
                    </MeetingInviteDialog>
                </CardDescription>
            </CardHeader>
            <OverlayScroll class="flex-grow basis-0">
                <CardContent class="space-y-6">
                    <MeetingTabPartWaiting></MeetingTabPartWaiting>
                    <MeetingTabPartAttending></MeetingTabPartAttending>
                    <div class="h-[800px]"></div>
                </CardContent>
            </OverlayScroll>
        </Card>
    </TabsContent>
</template>
