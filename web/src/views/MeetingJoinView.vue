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
import { useUserInfo } from '@/stores/user-store'
import { onMounted, ref } from 'vue'
import { useMeetingData } from '@/stores/meeting-store'
import { useRouter } from 'vue-router'
import ErrorAlert from '@/components/ErrorAlert.vue'
import { ReloadIcon } from '@radix-icons/vue'
import { meetingFetch } from '@/logic/meeting-fetch'
import { HTTPError } from 'ky'
import { userFetch } from '@/logic/user-fetch'

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
const meetingData = useMeetingData()
const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')

const onSubmit = form.handleSubmit(async (values) => {
    loading.value = true
    try {
        const tokenJson = await userFetch.requestToken({ name: values.name })
        userInfo.data = {
            name: values.name,
            token: tokenJson.access_token
        }

        const json = await meetingFetch.join({ room_id: props.room_id! })
        meetingData.data = {
            id: json.id,
            name: json.name,
            token: json.token,
            createDate: new Date(json.create_date)
        }
        router.push('/meeting/room')
    } catch (e: any) {
        if (e instanceof HTTPError) errorMessage.value = e.message
    } finally {
        loading.value = false
    }
})

const meetingName = ref('TODO')
onMounted(() => {})
</script>

<template>
    <NavBar></NavBar>
    <div class="flex-grow flex justify-center items-center">
        <div class="px-4">
            <h1 class="scroll-m-20 text-4xl tracking-tight lg:text-5xl">
                Join meeting: {{ meetingName }}
            </h1>
            <form class="space-y-6 pt-6" @submit="onSubmit">
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
                <Button type="submit" :disabled="loading">
                    <template v-if="loading">
                        <ReloadIcon class="w-4 h-4 mr-2 animate-spin" />
                        Waiting for the host approval...
                    </template>
                    <template v-else> Join Meeting </template>
                </Button>
                <ErrorAlert :message="errorMessage"></ErrorAlert>
            </form>
        </div>
    </div>
</template>
