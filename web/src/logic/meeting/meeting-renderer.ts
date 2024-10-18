import { Track, type Participant, type Room } from 'livekit-client'
import { triggerRef, type ShallowRef } from 'vue'

export type MeetingRenderMap = Map<string, MeetingRender>

export type MeetingRender = {
    participantID: string
    videoSrc?: MediaStream
    audioSrc?: MediaStream
}

type Params = {
    room: Room
    renderMap: ShallowRef<MeetingRenderMap, MeetingRenderMap>
    setGridSize: (numParticipants: number) => void
}

export class MeetingRenderer {
    room: Room
    renderMap: ShallowRef<MeetingRenderMap, MeetingRenderMap>
    setGridSize: (numParticipants: number) => void

    constructor(params: Params) {
        this.room = params.room
        this.renderMap = params.renderMap
        this.setGridSize = params.setGridSize
    }

    /* render functions */

    updateRenderMap() {
        const map: MeetingRenderMap = new Map()
        const localParticipant = this.room.localParticipant
        map.set(localParticipant.identity, {
            videoSrc: this.getVideoStream(localParticipant),
            // dont include local audio
            participantID: localParticipant.identity
        })
        for (const [key, value] of this.room.remoteParticipants) {
            map.set(key, {
                participantID: key,
                videoSrc: this.getVideoStream(value),
                audioSrc: this.getAudioStream(value)
            })
        }
        console.log(`rendered ${map.size} participants`)
        this.renderMap.value = map
    }

    getVideoStream(participant: Participant): MediaStream | undefined {
        if (!participant.isCameraEnabled) return undefined
        return participant.getTrackPublication(Track.Source.Camera)?.videoTrack?.mediaStream
    }

    getAudioStream(participant: Participant): MediaStream | undefined {
        if (!participant.isMicrophoneEnabled) return undefined
        return participant.getTrackPublication(Track.Source.Microphone)?.audioTrack?.mediaStream
    }

    renderGrid() {
        this.updateRenderMap()
        this.setGridSize(this.renderMap.value.size)
        triggerRef(this.renderMap)
    }
}
