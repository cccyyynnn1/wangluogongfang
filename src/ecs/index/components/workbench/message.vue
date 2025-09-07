<script lang="ts">
  export default {
    name: 'Message', // 系统信息
  }
</script>

<script setup lang="ts">
  import { useSettingsStore } from '@/store/modules/settings'
  // import { GlobalMessageItem } from '~/types/store'
  import { formatTime } from '@/utils/time'
  import { usePubilcStore } from '@/store/modules/public'
  const settingsStore = useSettingsStore()
  const publicStore = usePubilcStore()
  const { setGlobalMessageList } = publicStore
  onMounted(() => {
    setGlobalMessageList()
  })
  const { globalMessageList } = storeToRefs(publicStore)

  const curIndex = ref()

  const handleItemClick = (id: number) => {
    curIndex.value = id
    settingsStore.changeMessageVisible(true, id)
  }
</script>

<template>
  <vab-card class="message">
    <template #header>
      <span class="title">消息通知</span>
    </template>
    <div class="content">
      <div class="content_warp">
        <div
          v-for="item in globalMessageList"
          :key="item.id"
          class="item"
          :class="{ itemActive: item.id == curIndex }"
          @click="handleItemClick(item.id)"
        >
          <div class="item_left">
            <el-image
              v-if="item.msgType == 'user'"
              class="img"
              :src="require(`@/assets/message_image/user-message.svg`)"
            />
            <el-image
              v-else-if="item.msgType == 'sys'"
              class="img"
              :src="require(`@/assets/message_image/system-message.svg`)"
            />
            <div v-if="!item.readStatus" class="mark"></div>
          </div>
          <div class="item_right">
            <div class="item_right-top">
              <div class="item_right_name">[{{ item.msgSource }}]</div>
              <div class="item_right_time">{{ formatTime(item.createTime) }}</div>
            </div>
            <div class="item_right-bottom">{{ item.notes }}</div>
          </div>
        </div>
        <el-empty v-if="globalMessageList.length == 0" description="暂无数据" style="height: 80%; width: 100%" />
      </div>

      <div class="content_bottom">
        <span style="opacity: 0.5" @click="settingsStore.changeMessageVisible(true)">查看所有消息</span>
      </div>
    </div>
  </vab-card>
</template>

<style scoped lang="scss">
  :deep() {
    .el-card__header {
      border-bottom: none;
      height: 40px;
    }
    .el-empty {
      --el-empty-image-width: 120px;
    }
    .title {
      color: #303133;
      font-size: 16px;
      font-weight: 500;
    }
    .content {
      width: 100%;
      height: 100%;
      .content_warp {
        padding: 10px 7px;
        height: calc(100% - 42px);
        overflow-y: auto;
        &::-webkit-scrollbar {
          width: 0;
          height: 0;
        }
        .itemActive {
          background-color: #fff;
        }
        .item {
          width: 100%;
          height: 76px;
          border-radius: 10px;
          display: flex;
          align-items: center;
          cursor: pointer;
          &:hover {
            background: #faf9ff;
          }
          .item_left {
            padding: 0 14px;
            display: flex;
            height: 76px;
            width: 76px;
            align-items: center;
            justify-content: center;
            overflow: hidden;
            position: relative;
            .mark {
              position: absolute;
              width: 10px;
              height: 10px;
              background: #ff6767;
              border-radius: 50%;
              right: 14px;
              top: 14px;
            }
            :deep() {
              .el-image__inner {
                width: 48px;
                height: 48px;
                background: #ebeaff;
                border-radius: 16px;
              }
            }
          }
          .item_right {
            width: calc(100% - 76px);
            height: 40px;
            display: flex;
            flex-direction: column;
            padding-right: 12px;
            justify-content: space-between;
            .item_right-top {
              display: flex;
              align-items: center;
              justify-content: space-between;
            }
            .item_right-bottom {
              font-size: 13px;
              color: #606266;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
              word-break: break-all;
              word-wrap: break-word;
            }
            .item_right_name {
              font-size: 14px;
              color: #4a4759;
              font-weight: 500;
              width: 80px;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
              word-break: break-all;
              word-wrap: break-word;
            }
            .item_right_time {
              width: calc(100% - 80px);
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
              word-break: break-all;
              word-wrap: break-word;
              font-size: 12px;
              color: #b7b5bf;
              text-align: right;
            }
          }
        }
      }
      .content_bottom {
        height: 42px;
        font-size: 13px;
        color: #4a4759;
        line-height: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-top: 1px solid #f1f0ff;
        cursor: pointer;
      }
    }
  }
</style>
