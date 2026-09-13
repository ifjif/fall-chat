import { convMembListBriefBatchApi } from "@/api/conversation_member";

export class FindConvMembBriefJob {
  id;
  value;
  values = new Set();
  messages = [];
  result;
  conv;
  completeCallback;

  constructor(conv, message, completeCallback) {
    this.conv = conv
    this.id = conv.id
    this.value = message.sender_id
    this.values.add(this.value)
    this.messages.push(message)
    this.completeCallback = completeCallback;
  }

  key() {
    return this.id
  }

  delay() {
    return 500
  }

  merge(job) {
    this.values.add(job.value)
  }

  async execute() {
    const members = Array.from(this.values)
    const data = {
      conv_id: this.id,
      member_ids: members
    }

    // 进行请求
    this.result = await convMembListBriefBatchApi(data)
  }
  onComplete() {
    if (this.completeCallback) {
      this.completeCallback(this.result, this.conv, this.messages)
    }
  }
}
