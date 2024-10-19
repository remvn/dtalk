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
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { TabsContent } from '@/components/ui/tabs'
</script>

<template>
    <TabsContent value="chat" class="h-full">
        <Card class="h-full border-0">
            <CardHeader>
                <CardTitle>Chat</CardTitle>
                <CardDescription> You can send messages to other people here. </CardDescription>
            </CardHeader>
            <CardContent class="space-y-2"> Content </CardContent>
            <CardFooter> Test </CardFooter>
        </Card>
    </TabsContent>
</template>
