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
    <TabsContent value="participant">
        <Card>
            <CardHeader>
                <CardTitle>Participants</CardTitle>
                <CardDescription>
                    Make changes to your account here. Click save when you're done.
                </CardDescription>
            </CardHeader>
            <CardContent class="space-y-2">
                <div class="space-y-1">
                    <Label for="name">Name</Label>
                    <Input id="name" default-value="Pedro Duarte" />
                </div>
                <div class="space-y-1">
                    <Label for="username">Username</Label>
                    <Input id="username" default-value="@peduarte" />
                </div>
            </CardContent>
        </Card>
    </TabsContent>
</template>
