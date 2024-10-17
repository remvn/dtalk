import './assets/index.css'
import './assets/base.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { VueQueryPlugin, type VueQueryPluginOptions } from '@tanstack/vue-query'

const app = createApp(App)

app.use(createPinia())
app.use(router)

const vueQueryOptions: VueQueryPluginOptions = {
    queryClientConfig: {
        defaultOptions: {
            queries: {
                retry: false,
                refetchOnWindowFocus: false
            }
        }
    }
}
app.use(VueQueryPlugin, vueQueryOptions)

app.mount('#app')
