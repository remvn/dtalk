export function getAppMode() {
    return import.meta.env.MODE
}

export function getAPIBaseURL() {
    const mode = getAppMode()
    if (mode === 'production') {
        const loc = window.location
        return `${loc.protocol}//${loc.hostname}`
    } else {
        return 'http://localhost:8080'
    }
}

export function getLkServerURL() {
    const mode = getAppMode()
    if (mode === 'production') {
        const loc = window.location
        return `ws://livekit.${loc.hostname}`
    } else {
        return 'ws://localhost:7880'
    }
}
