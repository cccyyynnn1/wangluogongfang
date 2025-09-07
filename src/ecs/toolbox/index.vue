<script lang="ts">
  export default {
    name: 'Toolboxs',
  }
</script>
<script setup lang="ts">
  import { useSettingsStore } from '@/store/modules/settings'
  import EncodingDecoding from './components/encoding-decoding.vue'
  import Unixtime from './components/unixtime.vue'
  import Hexconvert from './components/hexconvert.vue'
  import ResetPassword from './components/reset-password.vue'
  import Subnetmask from './components/subnetmask.vue'
  import Regexp from './components/regexp.vue'

  import { ToolType } from '/#/store'
  import { Component } from 'vue'
  const settingsStore = useSettingsStore()
  interface Props {
    modelValue: boolean
    toolType: ToolType
  }
  const toollDom: { [key in ToolType]?: Component } = {
    'encoding-decoding': EncodingDecoding,
    'reset-password': ResetPassword,
    hexconvert: Hexconvert,
    unixtime: Unixtime,
    subnetmask: Subnetmask,
    regexp: Regexp,
  }
  const titles: { [key in ToolType]?: string } = {
    'reset-password': '重置密码',
    'encoding-decoding': '编解码工具',
    subnetmask: '子网掩码计算',
    regexp: '正则表达式',
    all: '工具箱',
    unixtime: '时间戳转换工具',
    hexconvert: '进制转换工具',
  } as const
  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
  }>()
  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
  })
  const visible = useVModel(props, 'modelValue', emits)
  const title = computed(() => titles[props.toolType])
</script>

<template>
  <template v-if="props.toolType === 'all'">
    <el-dialog v-model="visible" title="工具箱" width="1180px">
      <div class="toolbox-container">
        <div
          class="toolbox-item toolbox-item-1"
          @click="() => settingsStore.changeToolboxVisible(true, 'encoding-decoding')"
        >
          <h3>编解码小工具</h3>
          <p>Base64、URL编解码、unicode中文互转等常用编解码工具</p>
        </div>
        <div class="toolbox-item toolbox-item-2" @click="() => settingsStore.changeToolboxVisible(true, 'hexconvert')">
          <h3>进制转换</h3>
          <p>二进制、10进制、16进制等常用进制转换</p>
        </div>
        <div class="toolbox-item toolbox-item-3" @click="() => settingsStore.changeToolboxVisible(true, 'unixtime')">
          <h3>时间戳转换</h3>
          <p>时间戳、时间互转</p>
        </div>
        <div class="toolbox-item toolbox-item-4" @click="() => settingsStore.changeToolboxVisible(true, 'subnetmask')">
          <h3>子网掩码计算</h3>
          <p>子网掩码计算、子网范围计算、子网可用IP</p>
        </div>
        <div
          v-if="false"
          class="toolbox-item toolbox-item-5"
          @click="() => settingsStore.changeToolboxVisible(true, 'regexp')"
        >
          <h3>正则表达式</h3>
          <p>正则表达式测试工具</p>
        </div>
      </div>
    </el-dialog>
  </template>
  <template v-else>
    <el-dialog v-model="visible" :title="title" :width="props.toolType === 'reset-password' ? '700px' : '1180px'">
      <component :is="toollDom[props.toolType]" />
    </el-dialog>
  </template>
</template>

<style scoped lang="scss">
  .toolbox-container {
    display: flex;
    padding: 38px 70px;
    flex-wrap: wrap;
    .toolbox-item {
      width: 300px;
      height: 150px;
      margin-bottom: 48px;
      border-radius: 10px;
      cursor: pointer;
      font-family: PingFangSC-Medium, PingFang SC;
      h3 {
        font-size: 18px;
        color: #1b242d;
        margin: 37px 0 0 21px;
      }
      p {
        width: 173px;
        font-size: 12px;
        font-weight: 400;
        color: #858c93;
        margin: 12px 0 0 21px;
        line-height: 20px;
      }
      &:not(:nth-of-type(3n + 1)) {
        margin-left: 48px;
      }
      &.toolbox-item-1 {
        background: url(@/assets/toolbox_imgaes/box1.png);
        &:hover {
          background: url(@/assets/toolbox_imgaes/box1hover.png);
        }
      }
      &.toolbox-item-2 {
        background: url(@/assets/toolbox_imgaes/box2.png);
        &:hover {
          background: url(@/assets/toolbox_imgaes/box2hover.png);
        }
      }
      &.toolbox-item-3 {
        background: url(@/assets/toolbox_imgaes/box3.png);
        &:hover {
          background: url(@/assets/toolbox_imgaes/box3hover.png);
        }
      }
      &.toolbox-item-4 {
        background: url(@/assets/toolbox_imgaes/box4.png);
        &:hover {
          background: url(@/assets/toolbox_imgaes/box4hover.png);
        }
      }
      &.toolbox-item-5 {
        background: url(@/assets/toolbox_imgaes/box5.png);
        &:hover {
          background: url(@/assets/toolbox_imgaes/box5hover.png);
        }
      }
    }
  }
</style>
