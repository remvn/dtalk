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
import { useForm } from 'vee-validate'
import { useRouter } from 'vue-router'
import { getResMessage } from '@/logic/fetching'
import { errorToast } from '@/logic/toast'
import { meetingFetch } from '@/logic/meeting-fetch'

const isOpen = ref(false)

const formSchema = toTypedSchema(
    z.object({
        room_name: z.string().min(2).max(50)
    })
)

const form = useForm({
    validationSchema: formSchema
})

const router = useRouter()
const loading = ref(false)

const onSubmit = form.handleSubmit(async (values) => {
    loading.value = true
    try {
        const res = await meetingFetch.create({
            room_name: values.room_name
        })
        if (!res.ok) {
            errorToast(await getResMessage(res))
            return
        }
        const json = await res.json()
        router.push(`/meeting/join?id=${json.room_id}`)
    } catch (e: any) {
        errorToast(e.message)
    } finally {
        loading.value = false
        isOpen.value = false
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
                                placeholder="Enter room name..."
                                v-bind="componentField"
                            />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                </FormField>
            </form>

            <DialogFooter>
                <Button @click="onSubmit" :disabled="loading">
                    <template v-if="loading">Creating...</template>
                    <template v-else>Create</template>
                </Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>
