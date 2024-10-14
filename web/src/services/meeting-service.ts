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

type MeetingParams = {
    url: string
    token: string
}

export class Meeting {
    room: Room
    url: string
    token: string
    container: HTMLDivElement | null

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
        this.container = null
        this.url = params.url
        this.token = params.token
    }

    async connect() {
        await this.room.connect(this.url, this.token)
    }

    async disconnect() {
        await this.room.disconnect(true)
    }

    setListener() {
        this.room
            .on(RoomEvent.TrackSubscribed, this.handleTrackSubscribed)
            .on(RoomEvent.TrackUnsubscribed, this.handleTrackUnsubscribed)
            .on(RoomEvent.ActiveSpeakersChanged, this.handleActiveSpeakerChange)
            .on(RoomEvent.Disconnected, this.handleDisconnect)
            .on(RoomEvent.LocalTrackUnpublished, this.handleLocalTrackUnpublished)
    }

    handleTrackSubscribed(
        track: RemoteTrack,
        publication: RemoteTrackPublication,
        participant: RemoteParticipant
    ) {
        console.log(`called`)
        if (track.kind === Track.Kind.Video || track.kind === Track.Kind.Audio) {
            // attach it to a new HTMLVideoElement or HTMLAudioElement
            const element = track.attach()
            console.log(this.container)
            this.container?.appendChild(element)
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
