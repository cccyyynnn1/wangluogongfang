<script lang="ts">
  export default {
    name: 'Hexconvert',
  }
</script>

<script setup lang="ts">
  const NUMERALS = '0123456789'
  const LETTERS_LOWERCASE = 'abcdefghijklmnopqrstuvwxyz'
  const LETTERS_UPPERCASE = LETTERS_LOWERCASE.toUpperCase()

  const binary_key: { [key: number]: string } = {
    52: LETTERS_LOWERCASE + LETTERS_UPPERCASE,
    58: (NUMERALS + LETTERS_LOWERCASE + LETTERS_UPPERCASE).replace(/[0OlI]/g, ''),
    62: NUMERALS + LETTERS_LOWERCASE + LETTERS_UPPERCASE,
  }

  const binarys = [
    {
      key: '2进制',
      value: 2,
    },
    {
      key: '8进制',
      value: 8,
    },
    {
      key: '10进制',
      value: 10,
    },
    {
      key: '16进制',
      value: 16,
    },
    {
      key: '26进制',
      value: 26,
    },
    {
      key: '32进制',
      value: 32,
    },
    {
      key: '36进制',
      value: 36,
    },
    {
      key: '52进制',
      value: 52,
    },
    {
      key: '58进制',
      value: 58,
    },
    {
      key: '62进制',
      value: 62,
    },
  ]
  const binaryType = ref(2)
  const binaryVal = ref('')
  const binaryNumList = ref<{ value: number; result: string }[]>(binarys.map(({ value }) => ({ value, result: '' })))
  const to10 = (baseNum: string, binary: number) => {
    if (binary <= 36) return parseInt(baseNum, binary)
    let val = 0
    baseNum = baseNum.split('').reverse().join('')
    const curKey = binary_key[binary]
    for (let i = 0; i < baseNum.length; i++) {
      const c = baseNum[i]
      val += curKey.indexOf(c) * Math.pow(curKey.length, i)
    }
    return val
  }

  const toTargetBinary = (baseNum: number, binary: number) => {
    if (binary <= 36) return baseNum.toString(binary)
    const arr = []
    const curKey = binary_key[binary]
    while (baseNum > 0) {
      arr.push(curKey[baseNum % curKey.length])
      baseNum = Math.floor(baseNum / curKey.length)
    }
    //数组反转，因为个位在索引0的位置，应反过来显示
    return arr.reverse().join('')
  }

  const formatBinary = (num: number, binary: number) => {
    return toTargetBinary(num, binary)
  }

  const getBinaryNum = () => {
    const num = to10(binaryVal.value, binaryType.value)
    const arr = binarys.map(({ value }) => ({ value, result: formatBinary(num, value) }))
    binaryNumList.value = arr
  }
</script>

<template>
  <div class="unixtime">
    <el-row>
      <el-col :span="3">
        <el-select v-model="binaryType" :style="{ width: '120px', marginRight: '20px' }">
          <el-option v-for="item in binarys" :key="item.value" :label="item.key" :value="item.value" />
        </el-select>
      </el-col>
      <el-col :span="19">
        <el-input v-model="binaryVal" style="width: 100%" />
      </el-col>
      <el-col :span="2">
        <div style="text-align: end">
          <el-button type="primary" @click="getBinaryNum">转换</el-button>
        </div>
      </el-col>
    </el-row>

    <el-table border :data="binaryNumList" style="margin-top: 20px">
      <el-table-column align="center" label="进制" prop="value" show-overflow-tooltip width="80" />
      <el-table-column align="center" label="结果" prop="result" show-overflow-tooltip />
    </el-table>
  </div>
</template>

<style scoped lang="scss"></style>
