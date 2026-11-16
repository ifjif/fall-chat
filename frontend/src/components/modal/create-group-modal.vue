<template>
  <fl-modal :beforeChange="handleConfirm" v-model="modelValue" title="创建群聊" fixedScreen>
    <template #content>
      <fl-scrollbar>
        <fl-form ref="formRef" :rules="rules" :model="groupInfoModel">
          <fl-form-item prop="avatar" label="头像">
            <div style="min-width:80px">
              <ImageUpload size="50" @change="imageChange"></ImageUpload>
            </div>
          </fl-form-item>
          <fl-form-item prop="name" label="群名">
            <fl-input v-model="groupInfoModel.name" maxlength="100" count></fl-input>
          </fl-form-item>
          <fl-form-item prop="max_members" label="成员数">
            <fl-input v-model="groupInfoModel.max_members"></fl-input>
          </fl-form-item>
          <fl-form-item prop="join_mode" label="加入模式">
            <fl-radio-group v-model="groupInfoModel.join_mode" size="small">
              <fl-radio label="自由加入" :value="0"></fl-radio>
              <fl-radio label="需审批" :value="1"></fl-radio>
              <fl-radio label="禁止加入" :value="2"></fl-radio>
            </fl-radio-group>
          </fl-form-item>
          <fl-form-item prop="search_flag" label="搜索行为">
            <fl-radio-group v-model="groupInfoModel.search_flag" size="small">
              <fl-radio label="所有人可搜" :value="1"></fl-radio>
              <fl-radio label="仅通过id" :value="0"></fl-radio>
            </fl-radio-group>
          </fl-form-item>
          <fl-form-item label="描述">
            <fl-textarea v-model="groupInfoModel.description" maxlength="255" count></fl-textarea>
          </fl-form-item>
        </fl-form>
      </fl-scrollbar>
    </template>
  </fl-modal>
</template>

<script setup>
import { createGroupApi } from '@/api/group'
import ImageUpload from '../common/image-upload.vue'
import { ref } from 'vue'
import { uploadFile } from '@/utils/file_upload'

const modelValue = defineModel({ type: Boolean })

const formRef = ref(null)

let file_ = null

const imageChange = (file) => {
  file_ = file
  console.log(file)
}

const groupInfoModel = ref({
  avatar: "",
  name: "",
  description: "",
  max_members: "500",
  search_flag: 1,
  join_mode: 0
})

const rules = {
  'avatar': [{ type: 'string', required: true, message: '请上传头像', trigge: 'change' }],
  'name': [{ type: 'string', required: true, message: '请输入群名', trigge: 'blur' }],
  'max_members': [{ required: true, message: '请指定最大成员数', trigge: 'blur' }],
  'search_flag': [{ type: 'number', required: true, message: '请设置搜索标志', trigge: 'change' }],
  'join_mode': [{ type: 'number', required: true, message: '请设置加入模式', trigge: 'change' }],
}

const handleConfirm = async () => {
  let fileUrl = groupInfoModel.value.avatar
  if (!file_ && !fileUrl) {
    try {
      await formRef.value.validate()
    } catch (e) {
      return
    }
  }

  if (!fileUrl) {
    // 文件存在, 进行上传
    uploadFile(file_, ({ fileUrl }) => {
      groupInfoModel.value.avatar = fileUrl
      doGroupCreate()
    })

    return
  }

  doGroupCreate()

}

const doGroupCreate = async () => {
  try {
    await formRef.value.validate()
  } catch (e) {
    return
  }
  const req = { ...groupInfoModel.value }
  req.max_members = parseInt(req.max_members)

  try {
    await createGroupApi(req)
    modelValue.value = false
  } catch (err) {
  }
}
</script>
