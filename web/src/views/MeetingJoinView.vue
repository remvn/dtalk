<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import NavBar from '@/components/NavBar.vue'
import { Button } from '@/components/ui/button'
import {
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { joinMeeting, requestToken } from '@/services/user-service'
import { useUserInfo } from '@/stores/user-info'

const props = defineProps({
    room_id: String
})

const formSchema = toTypedSchema(
    z.object({
        name: z.string().min(2).max(50)
    })
)

const form = useForm({
    validationSchema: formSchema
})

const userInfo = useUserInfo()

const onSubmit = form.handleSubmit(async (values) => {
    try {
        let json = await requestToken({ name: values.name })
        userInfo.setInfo({ name: values.name, token: json.access_token })
        json = await joinMeeting({ room_id: props.room_id! })
    } catch (e) {
        console.log(e)
    }
})
</script>

<template>
    <NavBar></NavBar>
    <div class="flex-grow">
        <h1 class="scroll-m-20 text-4xl tracking-tight lg:text-5xl">
            Join meeting: TODO: meeting name
        </h1>
        <form class="w-2/3 space-y-6" @submit="onSubmit">
            <FormField v-slot="{ componentField }" name="name">
                <FormItem>
                    <FormLabel>Join under name: </FormLabel>
                    <FormControl>
                        <Input
                            type="text"
                            placeholder="Enter your display name"
                            v-bind="componentField"
                        />
                    </FormControl>
                    <FormDescription>Other people will see this as your name.</FormDescription>
                    <FormMessage />
                </FormItem>
            </FormField>
            <Button type="submit"> Join Meeting </Button>
        </form>
    </div>
</template>
