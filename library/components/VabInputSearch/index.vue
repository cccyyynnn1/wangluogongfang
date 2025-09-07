<script lang="ts">
  export default {
    name: 'VabInputSearch', // 搜索框
  }
</script>

<script setup lang="ts">
  import { ref } from 'vue'

  import { useVModel } from '@vueuse/core'

  const loading = ref(false)

  const props = defineProps<{
    modelValue: any
    clearable?: boolean
  }>()

  const emit = defineEmits<{
    (e: 'update:modelValue', str: string): void
    (e: 'on-search'): void
  }>()

  const handleClick = () => {
    emit('on-search')
    loading.value = true
    setTimeout(() => {
      loading.value = false
    }, 300)
  }
  const inputValue = useVModel(props, 'modelValue', emit)
  const handleClear = () => {
    emit('update:modelValue', '')
  }
</script>

<template>
  <el-input v-model="inputValue">
    <template #suffix>
      <el-icon style="cursor: pointer; font-size: 16px; margin-right: 4px" @click="handleClear">
        <CircleClose v-if="clearable && inputValue" />
      </el-icon>
      <svg
        class="search"
        :class="{ activeSearch: loading }"
        height="16px"
        version="1.1"
        viewBox="0 0 16 16"
        width="16px"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
        @click.stop="handleClick"
      >
        <title>搜索_search</title>
        <g id="站点" fill="none" fill-rule="evenodd" stroke="none" stroke-linejoin="round" stroke-width="1">
          <g id="729站点--API列表" stroke="#B5B2C3" stroke-width="1.7" transform="translate(-292, -34)">
            <g id="二级导航" transform="translate(76, 0)">
              <g id="资产数据量" transform="translate(13, 24)">
                <g id="搜索_search" transform="translate(204, 11)">
                  <path
                    id="路径"
                    d="M6.31184213,12.6236843 C9.79776115,12.6236843 12.6236843,9.79776115 12.6236843,6.31184213 C12.6236843,2.82592311 9.79776115,0 6.31184213,0 C2.82592311,0 0,2.82592311 0,6.31184213 C0,9.79776115 2.82592311,12.6236843 6.31184213,12.6236843 Z"
                  />
                  <line id="路径" stroke-linecap="round" x1="10.8495368" x2="14" y1="10.849574" y2="14.0000371" />
                </g>
              </g>
            </g>
          </g>
        </g>
      </svg>
    </template>
  </el-input>
</template>

<style scoped lang="scss">
  :deep() {
    .el-input__clear {
      display: none;
    }
  }
  svg {
    cursor: pointer;
    g {
      stroke: rgba(181, 178, 195, 1);
    }
  }
  .activeSearch {
    font-size: 13px;
    g {
      stroke: #6655e7;
    }
  }
</style>
