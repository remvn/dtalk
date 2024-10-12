<script setup lang="ts">
import { ref } from 'vue'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { Button } from '@/components/ui/button'
import {
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage
} from '@/components/ui/form'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { toast } from '@/components/ui/toast'
import { useForm } from 'vee-validate'
import { useUserInfo } from '@/stores/user-info'
import { createMeeting, requestToken } from '@/services/user-service'

const isOpen = ref(false)

const formSchema = toTypedSchema(
    z.object({
        name: z.string().min(2).max(50),
        room_name: z.string().min(2).max(50)
    })
)

const form = useForm({
    validationSchema: formSchema
})

const userInfo = useUserInfo()

const onSubmit = form.handleSubmit(async (values) => {
    try {
        isOpen.value = false
        let json = await requestToken({ name: values.name })
        userInfo.setInfo({
            name: values.name,
            token: json.access_token
        })
        json = await createMeeting({
            room_name: values.room_name
        })
    } catch (e) {
        toast({
            title: 'An error happened',
            description: 'Please try again later.',
            variant: 'destructive'
        })
        console.log(e)
    }
})
</script>

<template>
    <Dialog v-model:open="isOpen">
        <DialogTrigger as-child>
            <Button size="lg" variant="outline">Create meeting</Button>
        </DialogTrigger>
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>Create a new meeting</DialogTitle>
                <DialogDescription> </DialogDescription>
            </DialogHeader>

            <form class="space-y-4" @submit="onSubmit">
                <FormField v-slot="{ componentField }" name="room_name">
                    <FormItem>
                        <FormLabel>Room name:</FormLabel>
                        <FormControl>
                            <Input
                                type="text"
                                placeholder="Enter room name.."
                                v-bind="componentField"
                            />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="name">
                    <FormItem>
                        <FormLabel>Your display name:</FormLabel>
                        <FormControl>
                            <Input
                                type="text"
                                placeholder="Enter your display name.."
                                v-bind="componentField"
                            />
                        </FormControl>
                        <FormDescription> Other people will see this as your name </FormDescription>
                        <FormMessage />
                    </FormItem>
                </FormField>
            </form>

            <DialogFooter>
                <Button @click="onSubmit">Create & Join</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>
