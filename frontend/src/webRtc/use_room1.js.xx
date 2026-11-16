import { ref } from 'vue'
export function useRoom() {
  const roomId = ref('')
  const selfId = ref('')
  const remotePeers = ref([])

  function join(id, uid) {
    roomId.value = id
    selfId.value = uid
    remotePeers.value = []
  }

  function addRemote(peerId) {
    if (!remotePeers.value.includes(peerId)) {
      remotePeers.value.push(peerId)
    }
  }

  function removeRemote(peerId) {
    remotePeers.value = remotePeers.value.filter(id => id !== peerId)
  }

  function leave() {
    roomId.value = ''
    selfId.value = ''
    remotePeers.value = []
  }

  return {
    roomId,
    selfId,
    remotePeers,
    join,
    leave,
    addRemote,
    removeRemote
  }
}
