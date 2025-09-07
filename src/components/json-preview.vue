<script lang="ts">
  export default {
    name: 'JsonPreview',
  }
</script>

<script setup lang="ts">
  // @ts-ignore
  import { json } from '@codemirror/lang-json'
  import { Codemirror } from 'vue-codemirror'
  import { EditorState } from '@codemirror/state'

  defineProps({
    jsonValue: {
      type: String,
      default: '',
    },
  })
</script>

<template>
  <codemirror
    :model-value="jsonValue"
    :style="{
      width: '100%',
      height: '100%',
    }"
    v-bind="{
      ...{
        tabSize: 1,
        autoDestroy: true,
        extensions: [json(), EditorState.readOnly.of(true)],
      },
    }"
  />
</template>

<style scoped lang="scss">
  :deep() {
    .jv-code {
      background: #f8f7ff;
      box-shadow: null;
    }
    .v-codemirror,
    .cm-editor,
    .cm-content,
    .cm-scroller {
      position: relative;
      z-index: 1;
      // width: 100%;
      background: #f8f7ff;
      box-shadow: null;
      &::-webkit-scrollbar {
        height: 8px !important;
      }
    }
    .cm-line {
      // width: 100%;
      width: 1280px;
      padding: 2px 10px;
      white-space: normal;
      word-break: break-all;
      color: inherit !important;
      &::selection,
      & ::selection {
        background-color: var(--el-color-primary) !important;
        color: #fff;
      }
    }
    .el-textarea {
      margin: 16px 0 26px;
      color: #303133;
      .el-textarea__inner {
        background: #f8f7ff;
        padding: 20px;
        &:focus,
        &:hover {
          box-shadow: 0 0 0 1px var(--el-border-color) inset;
        }
      }
    }
    .cm-editor {
      width: max-content !important;
      height: fit-content;
    }
  }
</style>
