import type { InjectionKey, Ref } from 'vue'

export type MeetingTabState = {
    isDrawerOpen: boolean
    isDialogOpen: boolean
    selectedTab: string
}
export const MeetingTabStateKey = Symbol() as InjectionKey<Ref<MeetingTabState>>
