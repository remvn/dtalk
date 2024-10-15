import {
    Participant,
    RemoteParticipant,
    RemoteTrack,
    RemoteTrackPublication,
    Room,
    RoomEvent,
    VideoPresets,
    Track,
    LocalTrackPublication,
    LocalParticipant
} from 'livekit-client'
import { triggerRef, type ShallowRef } from 'vue'

export type MeetingRenderMap = Map<string, MeetingRender>

export type MeetingRender = {
    participantID: string
    videoSrc?: MediaStream
    audioSrc?: MediaStream
}

type MeetingParams = {
    url: string
    token: string
    renderMap: ShallowRef<MeetingRenderMap, MeetingRenderMap>
    setGridSize: (numParticipants: number) => void
}

export class Meeting {
    room: Room
    url: string
    token: string
    renderMap: ShallowRef<MeetingRenderMap, MeetingRenderMap>
    setGridSize: (numParticipants: number) => void

    constructor(params: MeetingParams) {
        this.room = new Room({
            // automatically manage subscribed video quality
            adaptiveStream: true,

            // optimize publishing bandwidth and CPU for published tracks
            dynacast: true,

            // default capture settings
            videoCaptureDefaults: {
                resolution: VideoPresets.h720.resolution
            }
        })
        this.renderMap = params.renderMap
        this.url = params.url
        this.token = params.token
        this.setGridSize = params.setGridSize
    }

    async connect() {
        await this.room.connect(this.url, this.token)
    }

    async disconnect() {
        await this.room.disconnect(true)
    }

    setListener() {
        this.room
            .on(RoomEvent.LocalTrackPublished, this.handleLocalTrackPublished.bind(this))
            .on(RoomEvent.LocalTrackUnpublished, this.handleLocalTrackUnpublished.bind(this))
            .on(RoomEvent.TrackSubscribed, this.handleTrackSubscribed.bind(this))
            .on(RoomEvent.TrackUnsubscribed, this.handleTrackUnsubscribed.bind(this))
            .on(RoomEvent.ActiveSpeakersChanged, this.handleActiveSpeakerChange.bind(this))
            .on(RoomEvent.Disconnected, this.handleDisconnect.bind(this))
    }

    handleLocalTrackPublished(pub: LocalTrackPublication, participant: LocalParticipant) {
        const track = pub.track
        if (pub.kind !== Track.Kind.Video) return
        if (track == null || track.mediaStream == null) return

        this.rerenderGrid()
    }

    handleLocalTrackUnpublished(pub: LocalTrackPublication, participant: LocalParticipant) {
        // when local tracks are ended, update UI to remove them from rendering
        pub.track?.detach()
        this.rerenderGrid()
    }

    handleTrackSubscribed(
        track: RemoteTrack,
        pub: RemoteTrackPublication,
        participant: RemoteParticipant
    ) {
        this.rerenderGrid()
    }

    handleTrackUnsubscribed(
        track: RemoteTrack,
        publication: RemoteTrackPublication,
        participant: RemoteParticipant
    ) {
        // remove tracks from all attached elements
        track.detach()
        this.rerenderGrid()
    }

    handleActiveSpeakerChange(speakers: Participant[]) {
        // show UI indicators when participant is speaking
    }

    handleDisconnect() {
        console.log('disconnected from room')
    }

    updateRenderMap() {
        const map: MeetingRenderMap = new Map()
        const localParticipant = this.room.localParticipant
        map.set(localParticipant.identity, {
            videoSrc: this.getVideoStream(localParticipant),
            participantID: localParticipant.identity
        })
        for (const [key, value] of this.room.remoteParticipants) {
            map.set(key, {
                participantID: key,
                videoSrc: this.getVideoStream(value),
                audioSrc: this.getAudioStream(value)
            })
        }
        console.log(map)
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

    rerenderGrid() {
        this.updateRenderMap()
        this.setGridSize(this.renderMap.value.size)
        triggerRef(this.renderMap)
    }
}
