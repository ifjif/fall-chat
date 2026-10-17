import { onUnmounted } from "vue";
import { useMediaStream } from "./use_media_stream";
import { usePeerConnection } from "./use_peer_connection";
import { useSignaling } from "./use_signaling";
import { useRoom } from "./use_room";
import { FlMessageBox } from "fall-ui";


export function useWebRtc() {
  const media = useMediaStream()
  const signaling = useSignaling()
  const peer = usePeerConnection(signaling, () => media.localStream.value)
  const room = useRoom()

  async function joinRoom(id, uid) {
    await media.start()
    signaling.send({
      type: 'join',
      data: {
        room_id: id,
        peer_id: uid
      }
    })
    room.join(id, uid)
  }

  function setupSignalingHandlers() {
    signaling.on('joined', msg => {
      // members是房间已有的人
      for (const id of msg.members) {
        // 只能由一端进行offer
        if (id > room.selfId.value) {
          peer.createOffer(id)
          room.addRemote(id)
        }
      }
    })

    signaling.on('invite', async msg => {
      const fb = FlMessageBox({
        type: 'info', onOk: () => {
          joinRoom(msg.room_id, msg.id)
          fb.close()
        },
        title: '收到邀请'
      })
    })

    signaling.on('peer-joined', async msg => {
      const peerId = msg.peerId
      if (peerId) {
        if (peerId > room.selfId.value) {
          // 只能由一端进行offer
          await peer.createOffer(peerId)
          room.addRemote(peerId)
        }
      }
    })

    signaling.on('peer-left', msg => {
      const peerId = msg.peerId
      if (peerId) {
        peer.removePeer(peerId)
        room.removeRemote(peerId)
      }
    })

    signaling.on('offer', async msg => {
      if (!msg.peer_id || !msg.data) return
      await peer.handleOffer(msg.peer_id, msg.data)
      room.addRemote(msg.peer_id)
    })

    signaling.on('answer', msg => {
      if (!msg.peer_id || !msg.data) return
      peer.handleAnswer(msg.peer_id, msg.data)
    })

    signaling.on('candidate', msg => {
      if (!msg.peer_id || !msg.data) return
      peer.addCandidate(msg.peer_id, msg.data)
    })
  }

  async function call(peerId) {
    signaling.send({
      type: 'invite',
      data: {
        room_id: room.roomId.value,
        peer_id: peerId
      }
    })
  }

  async function leaveRoom() {
    signaling.send({
      type: 'leave',
      data: {
        room_id: room.roomId.value,
        peer_id: room.selfId.value
      }
    })
    media.stop()
    peer.closeAll()
    room.leave()
  }

  onUnmounted(() => {
    leave()
  })

  setupSignalingHandlers()

  return {
    ...media,
    ...peer,
    ...room,
    signaling,
    joinRoom,
    call,
    leaveRoom
  }
}
