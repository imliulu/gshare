<template>
  <div class="app-container">
    <h1>Local Transfer --  房间管理</h1>
    <el-form :model="form" label-width="100px">
      <el-form-item label="创建房间">
        <el-input v-model="form.roomID" placeholder="请输入房间ID（可选）"></el-input>
        <el-input v-model="form.password" placeholder="请输入房间密码（可选）"></el-input>
        <el-button type="primary" @click="createRoom">创建房间</el-button>
      </el-form-item>
      <el-form-item label="在线房间">
        <div v-if="rooms.length === 0" class="no-rooms">当前无在线的房间</div>
        <div v-else>
          <div v-for="room in rooms" :key="room.ID" class="room-item">
            <div class="room-info">房间ID: {{ room.ID }} (创建于: {{ room.CreatedAt }})</div>
            <div class="room-actions">
              <el-button type="primary" @click="joinRoom(room.ID)">加入</el-button>
              <el-button type="danger" @click="deleteRoom(room.ID)">删除</el-button>
            </div>
          </div>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: 'IndexPage',
  data() {
    return {
      form: {
        roomID: '',
        password: ''
      },
      rooms: []
    }
  },
  created() {
    this.fetchOnlineRooms()
  },
  methods: {
    createRoom() {
      let roomID = this.form.roomID || Date.now().toString()
      let password = this.form.password || '123456' // 设置默认密码
      this.$api.post('/rooms', { id: roomID, password: password })
        .then(response => {
          this.$message.success(response.data.message)
          this.form.roomID = ''
          this.form.password = ''
          this.$router.push({ path: `/share/${roomID}`, query: { password: password } })
        })
        .catch(error => {
          this.$message.error('创建房间失败: ' + error.message)
        })
    },
    fetchOnlineRooms() {
      this.$api.get('/rooms')
        .then(response => {
          this.rooms = response.data.rooms
        })
        .catch(error => {
          this.$message.error('获取在线房间失败: ' + error.message)
        })
    },
    joinRoom(roomID) {
      this.$prompt('请输入房间密码', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'password'
      }).then(({ value }) => {
        this.$api.post('/rooms/join', { ID: roomID, Password: value })
          .then(response => {
            if (response.data.error) {
              this.$message.error(response.data.error)
            } else {
              this.$message.success(response.data.message)
              this.$router.push({ path: `/share/${roomID}`, query: { password: value } })
            }
          })
          .catch(error => {
            this.$message.error('加入房间失败: ' + error.message)
          })
      }).catch(() => {})
    },
    deleteRoom(roomID) {
      this.$prompt('请输入房间密码', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'password'
      }).then(({ value }) => {
        this.$api.delete(`/rooms/${roomID}`, { data: { password: value } })
          .then(response => {
            if (response.data.error) {
              this.$message.error(response.data.error)
            } else {
              this.$message.success(response.data.message)
              this.fetchOnlineRooms()
            }
          })
          .catch(error => {
            this.$message.error('删除房间失败: ' + error.message)
          })
      }).catch(() => {})
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}
.room-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border: 1px solid #ccc;
  margin-bottom: 5px;
  background-color: #f9f9f9;
}
.room-item:hover {
  background-color: #e9e9e9;
}
.no-rooms {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100px;
  text-align: center;
  font-size: 14px;
  color: #666;
  background-color: #f9f9f9;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-top: 20px;
}
</style>
