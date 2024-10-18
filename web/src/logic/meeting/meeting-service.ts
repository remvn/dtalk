import { meetingQuery } from '@/queries/meeting-query'
import type { QueryClient } from '@tanstack/vue-query'
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
    LocalParticipant,
    TrackPublication,
    DataPacket_Kind
} from 'livekit-client'
import { type ShallowRef } from 'vue'
import { MeetingMessenger } from '@/logic/meeting/meeting-messenger'
import { MeetingRenderer, type MeetingRenderMap } from '@/logic/meeting/meeting-renderer'

type Params = {
    url: string
    token: string
    renderMap: ShallowRef<MeetingRenderMap, MeetingRenderMap>
    setGridSize: (numParticipants: number) => void
    queryClient: QueryClient
}

export class Meeting {
    room: Room
    url: string
    token: string
    queryClient: QueryClient
    messenger: MeetingMessenger
    renderer: MeetingRenderer

    constructor(params: Params) {
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
        this.url = params.url
        this.token = params.token
        this.queryClient = params.queryClient

        this.renderer = new MeetingRenderer({
            room: this.room,
            setGridSize: params.setGridSize,
            renderMap: params.renderMap
        })
        this.messenger = new MeetingMessenger({
            room: this.room,
            queryClient: this.queryClient,
            renderer: this.renderer
        })
    }

    async connect() {
        await this.room.connect(this.url, this.token)
    }

    async disconnect() {
        await this.room.disconnect(true)
    }

    setListener() {
        this.room
            .on(RoomEvent.DataReceived, this.handleDataReceived.bind(this))
            .on(RoomEvent.ParticipantConnected, this.handleParticipantConnected.bind(this))
            .on(RoomEvent.ParticipantDisconnected, this.handleParticipantDisconnected.bind(this))
            .on(RoomEvent.LocalTrackPublished, this.handleLocalTrackPublished.bind(this))
            .on(RoomEvent.LocalTrackUnpublished, this.handleLocalTrackUnpublished.bind(this))
            .on(RoomEvent.TrackSubscribed, this.handleTrackSubscribed.bind(this))
            .on(RoomEvent.TrackUnsubscribed, this.handleTrackUnsubscribed.bind(this))
            .on(RoomEvent.TrackMuted, this.handleTrackMuted.bind(this))
            .on(RoomEvent.TrackUnmuted, this.handleTrackUnmuted.bind(this))
            .on(RoomEvent.ActiveSpeakersChanged, this.handleActiveSpeakerChange.bind(this))
            .on(RoomEvent.Disconnected, this.handleDisconnect.bind(this))
    }

    /* data message handler */

    handleDataReceived(
        payload: Uint8Array,
        participant?: RemoteParticipant | undefined,
        kind?: DataPacket_Kind | undefined,
        topic?: string | undefined
    ) {
        try {
            const text = new TextDecoder().decode(payload)
            const json = JSON.parse(text)
            this.messenger.handleMessage(json)
        } catch (e) {
            console.log(e)
        }
    }

    /* remote participant connections handler */

    handleParticipantConnected(remoteParticipant: RemoteParticipant) {
        console.log(remoteParticipant.name, 'connected')
        this.queryClient.invalidateQueries({
            queryKey: [...meetingQuery.keys.participants, this.room.name]
        })
        this.renderer.renderGrid()
    }

    handleParticipantDisconnected(remoteParticipant: RemoteParticipant) {
        console.log(remoteParticipant.name, 'disconnected')
        this.queryClient.invalidateQueries({
            queryKey: [...meetingQuery.keys.participants, this.room.name]
        })
        this.renderer.renderGrid()
    }

    /* local track publication handler */

    handleLocalTrackPublished(pub: LocalTrackPublication, participant: LocalParticipant) {
        // dont include local audio
        const track = pub.track
        if (pub.kind !== Track.Kind.Video) return
        if (track == null || track.mediaStream == null) return

        console.log(`local`, participant.name, pub.kind, 'subscribed')
        this.renderer.renderGrid()
    }

    handleLocalTrackUnpublished(pub: LocalTrackPublication, participant: LocalParticipant) {
        // when local tracks are ended, update UI to remove them from rendering
        console.log(`local`, participant.name, pub.kind, 'unsubsribed')
        pub.track?.detach()
        this.renderer.renderGrid()
    }

    /* remote track publication handler */

    handleTrackSubscribed(
        track: RemoteTrack,
        pub: RemoteTrackPublication,
        participant: RemoteParticipant
    ) {
        console.log('remote', participant.name, track.kind, 'subscribed')
        this.renderer.renderGrid()
    }

    handleTrackUnsubscribed(
        track: RemoteTrack,
        publication: RemoteTrackPublication,
        participant: RemoteParticipant
    ) {
        console.log('remote', participant.name, track.kind, 'unsubscibed')
        // remove tracks from all attached elements
        track.detach()
        this.renderer.renderGrid()
    }

    /* track mute handler */

    handleTrackMuted(pub: TrackPublication, participant: Participant) {
        this.logTrackHandle(pub, participant, 'muted')
        this.renderer.renderGrid()
    }

    handleTrackUnmuted(pub: TrackPublication, participant: Participant) {
        this.logTrackHandle(pub, participant, 'unmuted')
        this.renderer.renderGrid()
    }

    logTrackHandle(pub: TrackPublication, participant: Participant, eventType: string) {
        const location = participant.isLocal ? 'local' : 'remote'
        console.log(location, participant.name, pub.kind, `track ${eventType}`)
    }

    /* misc handler */

    handleActiveSpeakerChange(speakers: Participant[]) {
        // show UI indicators when participant is speaking
    }

    handleDisconnect() {
        console.log('disconnected from room')
    }
}
