import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'

export function useTwBreakpoints() {
    return useBreakpoints(breakpointsTailwind)
}
