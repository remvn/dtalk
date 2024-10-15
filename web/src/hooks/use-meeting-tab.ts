import { ref } from 'vue'
import { useTwBreakpoints } from '@/hooks/use-tw-breakpoints'
import type { MeetingTabState } from '@/types/meeting'

export function useMeetingTab() {
    const state = ref<MeetingTabState>({
        isDrawerOpen: false,
        isDialogOpen: false,
        selectedTab: 'participant'
    })
    const breakpoints = useTwBreakpoints()

    function toggleTab(name: string) {
        if (breakpoints.isGreaterOrEqual('lg')) {
            if (state.value.isDrawerOpen && state.value.selectedTab === name) {
                state.value.isDrawerOpen = false
                return
            }
            state.value.isDrawerOpen = true
        } else {
            if (state.value.isDialogOpen && state.value.selectedTab === name) {
                state.value.isDialogOpen = false
                return
            }
            state.value.isDialogOpen = true
        }
        state.value.selectedTab = name
    }

    return { state, toggleTab }
}
