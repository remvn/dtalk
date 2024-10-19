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
import { meetingFetch } from '@/logic/meeting/meeting-fetch'
import { HTTPError } from 'ky'
import { userFetch } from '@/logic/user-fetch'
import MdiArrowRight from '~icons/mdi/arrow-right'
import MdiArrowLeft from '~icons/mdi/arrow-left'
import { publicFetch } from '@/logic/public-fetch'

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
const lkClientURL = ref<string | undefined>()

const onSubmit = form.handleSubmit(async (values) => {
    loading.value = true
    try {
        if (userInfo.data.token == null) {
            const tokenJson = await userFetch.requestToken({ name: values.name })
            userInfo.data = {
                name: values.name,
                token: tokenJson.access_token
            }
        }
        if (lkClientURL.value == null) {
            const json = await publicFetch.getLivekitClientURL()
            lkClientURL.value = json.url
        }

        const json = await meetingFetch.join({ room_id: props.room_id! })
        meetingData.data = {
            id: json.id,
            name: json.name,
            token: json.token,
            createDate: new Date(json.create_date),
            wsUrl: lkClientURL.value
        }
        router.push('/meeting/room')
    } catch (e: any) {
        if (e instanceof HTTPError) errorMessage.value = e.message
    } finally {
        loading.value = false
    }
})

function handleGoBack() {
    router.push('/')
}

const meetingName = ref('TODO')
onMounted(async () => {
    if (props.room_id == null) {
        router.push('/')
        return
    }
    try {
        const json = await meetingFetch.publicData({ room_id: props.room_id })
        meetingName.value = json.name
    } catch (e: any) {
        errorMessage.value = `Unable to retrieve meeting's information`
    }
})
</script>

<template>
    <!-- <NavBar></NavBar> -->
    <div class="flex-grow flex justify-center items-center">
        <div class="p-4 w-full sm:w-[440px]">
            <div class="p-6 border border-secondary rounded-lg">
                <div>
                    <span class="font-thin text-muted-foreground">Join meeting: </span>
                    <h1 class="scroll-m-20 text-4xl tracking-tight lg:text-5xl">
                        {{ meetingName }}
                    </h1>
                </div>
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
                            <FormDescription
                                >Other people will see this as your name.</FormDescription
                            >
                            <FormMessage />
                        </FormItem>
                    </FormField>
                    <ErrorAlert :message="errorMessage"></ErrorAlert>
                    <div class="flex items-center justify-between">
                        <Button @click="handleGoBack" variant="outline">
                            <MdiArrowLeft class="size-4 mr-2"></MdiArrowLeft>
                            Go back
                        </Button>
                        <Button type="submit" :disabled="loading">
                            <template v-if="loading">
                                <ReloadIcon class="size-4 mr-2 animate-spin" />
                                Waiting for the host approval...
                            </template>
                            <template v-else>
                                Join Meeting
                                <MdiArrowRight class="size-4 ml-2"></MdiArrowRight>
                            </template>
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>
