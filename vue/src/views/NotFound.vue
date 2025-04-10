<template>
  <div class="app-container">
    <h1>404 - 页面未找到</h1>
    <p>您访问的页面不存在。</p>
    <el-button type="primary" @click="backToIndex">返回首页</el-button>
  </div>
</template>

<script>
export default {
  methods: {
    backToIndex() {
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
  text-align: center;
}
</style>
