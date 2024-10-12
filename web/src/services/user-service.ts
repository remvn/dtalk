export async function requestToken(body: { name: string }) {
    const res = await fetch('/api/auth/request-token', {
        method: 'POST',
        body: JSON.stringify(body)
    })
    if (res.status != 200) {
        return null
    }
    const json = await res.json()
    return json
}

export async function createMeeting(body: { room_name: string }) {
    const res = await fetch('/api/meeting/create', {
        method: 'POST',
        body: JSON.stringify(body)
    })
    const json = await res.json()
    return json
}
