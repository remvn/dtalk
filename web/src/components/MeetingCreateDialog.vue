<script setup lang="ts">
import { h, ref } from 'vue'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { Button } from '@/components/ui/button'
import {
    Form,
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

const isOpen = ref(false)

const formSchema = toTypedSchema(
    z.object({
        name: z.string().min(2).max(50)
    })
)

const form = useForm({
    validationSchema: formSchema
})

const onSubmit = form.handleSubmit((values) => {
    isOpen.value = false
    toast({
        title: 'You submitted the following values:',
        description: h(
            'pre',
            { class: 'mt-2 w-[340px] rounded-md bg-slate-950 p-4' },
            h('code', { class: 'text-white' }, JSON.stringify(values, null, 2))
        )
    })
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

            <form @submit="onSubmit">
                <FormField v-slot="{ componentField }" name="name">
                    <FormItem>
                        <FormLabel>Your display name:</FormLabel>
                        <FormControl>
                            <Input
                                type="text"
                                placeholder="Type your display name here"
                                v-bind="componentField"
                            />
                        </FormControl>
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
