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

export type MeetingRender = {
    participantID: string
    videoElement?: {
        srcObject: MediaProvider | null
    }
    audioElement?: HTMLAudioElement
}

type MeetingParams = {
    url: string
    token: string
    renderArr: ShallowRef<MeetingRender[], MeetingRender[]>
    setGridSize: (numParticipants: number) => void
}

export class Meeting {
    room: Room
    url: string
    token: string
    renderArr: ShallowRef<MeetingRender[], MeetingRender[]>
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
        this.renderArr = params.renderArr
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

    getRender(id: string): MeetingRender | null {
        for (const item of this.renderArr.value) {
            if (item.participantID === id) {
                return item
            }
        }
        return null
    }

    handleLocalTrackPublished(pub: LocalTrackPublication, participant: LocalParticipant) {
        const track = pub.track
        console.log('handleLocalTrackPublished', track?.kind)
        if (pub.kind !== Track.Kind.Video) return
        if (track == null || track.mediaStream == null) return

        const arr = this.renderArr.value
        let render = this.getRender(participant.identity)
        if (render == null) {
            render = {
                participantID: participant.identity,
                videoElement: {
                    srcObject: track.mediaStream
                }
            }
            arr.push(render)
        } else {
            render.videoElement = {
                srcObject: track.mediaStream
            }
        }
        this.renderArr.value = arr
        this.setGridSize(this.room.numParticipants)
        triggerRef(this.renderArr)
    }

    handleTrackSubscribed(
        track: RemoteTrack,
        pub: RemoteTrackPublication,
        participant: RemoteParticipant
    ) {
        console.log(`called`)
        if (track.kind === Track.Kind.Video || track.kind === Track.Kind.Audio) {
            // attach it to a new HTMLVideoElement or HTMLAudioElement
            const element = track.attach()
        }
    }

    handleTrackUnsubscribed(
        track: RemoteTrack,
        publication: RemoteTrackPublication,
        participant: RemoteParticipant
    ) {
        // remove tracks from all attached elements
        track.detach()
    }

    handleLocalTrackUnpublished(publication: LocalTrackPublication, participant: LocalParticipant) {
        // when local tracks are ended, update UI to remove them from rendering
        publication.track?.detach()
    }

    handleActiveSpeakerChange(speakers: Participant[]) {
        // show UI indicators when participant is speaking
    }

    handleDisconnect() {
        console.log('disconnected from room')
    }
}
