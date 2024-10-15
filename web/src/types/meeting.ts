import type { InjectionKey, Ref } from 'vue'

export type MeetingTabState = {
    isDrawerOpen: boolean
    isDialogOpen: boolean
    selectedTab: string
}
export type MeetingTabComposable = {
    state: Ref<MeetingTabState>
    toggleTab: (name: string) => void
}
export const MeetingTabComposableKey = Symbol() as InjectionKey<MeetingTabComposable>
