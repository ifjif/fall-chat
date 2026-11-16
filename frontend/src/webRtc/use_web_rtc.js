import { FlMessageBox } from "fall-ui";
import { useMediaStream } from "./use_media_stream";
import { usePeerConnection } from "./use_peer_connection";
import { useRoom } from "./use_room";
import { useSignaling } from "./use_signaling";
import { onUnmounted } from "vue";

export function useWebRtc() {
  const media = useMediaStream()
  const signaling = useSignaling()
  const peers = usePeerConnection(signaling, () => media.localStream.value)
  const room = useRoom()

  // {video,audio, mediaKind, roomId, userId}
  async function joinRoom(config) {
    await media.start(config.video, config.audio, config.mediaKind)
    signaling.send({
      type: 'join',
      data: {
        room_id: config.roomId,
        peer_id: config.userId
      }
    })
    room.create(config.roomId, config.userId)
  }

  function leaveRoom() {
    signaling.send({
      type: 'leave',
      data: {
        room_id: room.roomId.value,
        peer_id: room.selfId.value
      }
    })
    media.stop()
    peers.closeAll()
    room.leave()
  }

  function call(peerId) {
    signaling.send({
      type: 'invite',
      data: {
        room_id: room.roomId.value,
        peer_id: peerId
      }
    })
  }

  async function switchVideoDevice(deviceId) {
    const res = await media.switchDevice('user', 'video', deviceId)
    if (!res || !res.addedTrack) return
    peers.replaceLocalTrack(res.addedTrack)
  }

  async function switchAudioDevice(deviceId) {
    const res = await media.switchDevice('user', 'audio', deviceId)
    if (!res || !res.addedTrack) return
    peers.replaceLocalTrack(res.addedTrack)
  }

  function initSignaling() {
    signaling.onOffer(msg => {
      if (!msg.peer_id || !msg.data) return
      peers.handleOffer(msg.peer_id, msg.data)
      room.addRemote(msg.peer_id)
    })

    signaling.onAnswer(msg => {
      if (!msg.peer_id || !msg.data) return
      peers.handleAnswer(msg.peer_id, msg.data)
    })

    signaling.onCandidate(msg => {
      if (!msg.peer_id || !msg.data) return
      peers.addCandidate(msg.peer_id, msg.data)
    })

    signaling.on('joined', async msg => {
      // 只能一端进行 offer
      for (const id of msg.members) {
        if (id > room.selfId.value) {
          await peers.createOffer(id)
          room.addRemote(id)
        }
      }

    })
    signaling.on('peer-joined', async msg => {
      const peerId = msg.peer_id
      if (peerId && peerId > room.selfId.value) {
        await peers.createOffer(peerId)
        room.addRemote(peerId)
      }
    })
    signaling.on('peer-left', msg => {
      const peerId = msg.peer_id
      if (peerId) {
        peers.removePeer(peerId)
        room.removeRemote(peerId)
      }
    })
    signaling.on('invite', msg => {
      const fb = FlMessageBox({
        type: 'primary',
        onOk: () => {
          joinRoom({
            video: true,
            audio: true,
            mediaKind: 'user',
            roomId: msg.room_id,
            userId: msg.id
          })
          fb.close()
        },
        title: '收到邀请'
      })
    })
  }

  onUnmounted(() => {
    leaveRoom()
  })

  initSignaling()

  return {
    ...media,
    ...peers,
    ...room,
    joinRoom,
    call,
    leaveRoom,
    switchVideoDevice,
    switchAudioDevice
  }
}
