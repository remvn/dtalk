<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { TabsContent } from '@/components/ui/tabs'
import MeetingTabPartJoinRequest from '@/components/meeting/MeetingTabPartJoinRequest.vue'
import MeetingTabPartAttending from '@/components/meeting/MeetingTabPartAttending.vue'
import MeetingInviteDialog from './MeetingInviteDialog.vue'
import { Button } from '@/components/ui/button'
import MdiAccountPlusOutline from '~icons/mdi/account-plus-outline'
</script>

<template>
    <TabsContent value="participant" class="h-full">
        <Card class="h-full">
            <CardHeader>
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
            <CardContent class="space-y-3">
                <MeetingTabPartJoinRequest></MeetingTabPartJoinRequest>
                <MeetingTabPartAttending></MeetingTabPartAttending>
            </CardContent>
        </Card>
    </TabsContent>
</template>
