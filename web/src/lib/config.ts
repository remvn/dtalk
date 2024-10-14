export function getAppMode() {
    return import.meta.env.MODE
}

export function getAPIBaseURL() {
    const mode = getAppMode()
    if (mode === 'production') {
        return ''
    } else {
        return 'http://localhost:8080'
    }
}

export function getLkServerURL() {
    const mode = getAppMode()
    if (mode === 'production') {
        return 'TODO'
    } else {
        return 'ws://localhost:7880'
    }
}
