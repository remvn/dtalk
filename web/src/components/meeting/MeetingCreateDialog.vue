<script setup lang="ts">
import { ref } from 'vue'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
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
import { errorToast } from '@/logic/toast'
import { meetingFetch } from '@/logic/meeting/meeting-fetch'
import MdiVideoPlusOutline from '~icons/mdi/video-plus-outline'
import MdiArrowRight from '~icons/mdi/arrow-right'
import { ReloadIcon } from '@radix-icons/vue'

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
        const json = await meetingFetch.create({
            room_name: values.room_name
        })
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
            <Button size="lg" variant="outline">
                <MdiVideoPlusOutline class="size-5 mr-2"></MdiVideoPlusOutline>
                New meeting
            </Button>
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
                    <template v-if="loading">
                        <ReloadIcon class="size-4 mr-2 animate-spin" />
                        Creating...
                    </template>
                    <template v-else>
                        Create
                        <MdiArrowRight class="size-4 ml-2"></MdiArrowRight>
                    </template>
                </Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>
