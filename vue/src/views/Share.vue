<template>
  <div class="app-container">
    <h1>{{ currentTitle }}</h1>
    <el-button type="primary" @click="backToIndex">返回首页</el-button>
    <el-form :label-position="'top'" :model="form" label-width="100px">
      <el-form-item label="房间信息">
        <el-input v-model="displayRoomInfo" readonly></el-input>
      </el-form-item>
    </el-form>
    <el-tabs v-model="activeTab" @tab-click="updateTitle">
      <el-tab-pane label="剪贴板共享" name="clipboard">
        <el-form :label-position="'top'" :model="clipboardForm" label-width="100px">
          <el-form-item label="粘贴文本">
            <el-input type="textarea" v-model="clipboardForm.content" placeholder="请输入文本内容"></el-input>
            <br>
            <el-button type="primary" @click="pasteContent">粘贴文本</el-button>
          </el-form-item>
          <el-form-item label="获取文本">
            
            <div id="clipboardContents" class="clipboard-contents">
              <div v-for="(content, index) in clipboardContents" :key="index">{{ index + 1 }}. {{ content }}</div>
            </div>
            <el-button type="primary" @click="getClipboardContents">获取文本</el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="danger" @click="clearClipboard">清空剪贴板内容</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <el-tab-pane label="文件共享" name="files">
        <el-form :label-position="'top'" :model="fileForm" label-width="100px">
          <el-form-item label="上传文件">
            <el-upload
              class="upload-demo"
              :action="`${apiUrl}/rooms/${roomID}/upload`"
              :on-change="handleChange"
              :file-list="fileList"
              :auto-upload="false"
              :data="uploadData"
              ref="upload"
              multiple> <!-- 添加 multiple 属性 -->
              <!-- <el-button slot="trigger" type="primary">选取文件</el-button> -->
              <!-- <el-button type="primary">上传<i class="el-icon-upload el-icon--right"></i></el-button>
              <el-button type="success" @click="submitUpload">上传到服务器</el-button> -->

              <el-button style="width: 200px;" slot="trigger" size="small" type="primary">选取文件</el-button>
              <el-button style="width: 200px;" size="small" type="success" @click="submitUpload">上传到服务器<i class="el-icon-upload el-icon--right"></i></el-button>
            </el-upload>
          </el-form-item>
          <el-form-item label="文件列表">
            <div id="fileList" class="file-list">
              <div v-for="file in files" :key="file.name" class="file-item">
                <div>{{ file.name }}</div>
                <img v-if="file.isImage" :src="file.url" class="file-preview" @click="previewImage(file.url)" />
                <el-button type="primary" @click="downloadFile(file.url, file.name)">下载</el-button>
              </div>
            </div>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
export default {
  name: 'SharePage',
  data() {
    return {
      form: {
        roomID: '',
        password: ''
      },
      activeTab: 'clipboard',
      currentTitle: 'AirTransfer--剪贴板共享',
      clipboardForm: {
        content: ''
      },
      clipboardContents: [],
      fileForm: {},
      fileList: [],
      files: [],
      closeOnEsc: null, // 用于存储事件监听器
      roomID: '' // 添加 roomID 属性
    }
  },
  created() {
    this.form.roomID = this.$route.params.roomID;
    this.roomID = this.$route.params.roomID; // 初始化 roomID
    this.form.password = this.$route.query.password || '123456'; // 从路由参数中获取密码
    this.getClipboardContents();
    this.fetchFiles();
    this.updateTitle(); // 初始化标题
  },
  watch: {
    activeTab() {
      this.updateTitle();
    }
  },
  computed: {
    apiUrl() {
      return process.env.VUE_APP_API_URL;
    },
    uploadData() {
      return { roomID: this.roomID };
    },
    displayRoomInfo() {
      return `房间ID: ${this.form.roomID}, 密码: ${this.form.password}`;
    }
  },
  methods: {
    backToIndex() {
      this.$router.push('/');
    },
    pasteContent() {
      if (this.clipboardForm.content) {
        const payload = { content: this.clipboardForm.content.toString() };
        this.$api.post(`/rooms/${this.roomID}/clipboard`, payload, {
          headers: {
            'Content-Type': 'application/json'
          }
        })
          .then(response => {
            this.$message.success(response.data.message);
            this.clipboardForm.content = '';
            this.getClipboardContents();
          })
          .catch(error => {
            this.$message.error('粘贴失败: ' + error.message);
          });
      } else {
        this.$message.error('请输入文本内容');
      }
    },
    getClipboardContents() {
      this.$api.get(`/rooms/${this.roomID}/clipboard`)
        .then(response => {
          this.clipboardContents = response.data.contents;
        })
        .catch(error => {
          this.$message.error('获取失败: ' + error.message);
        });
    },
    clearClipboard() {
      this.$api.delete(`/rooms/${this.roomID}/clipboard`)
        .then(response => {
          this.$message.success(response.data.message);
          this.getClipboardContents();
        })
        .catch(error => {
          this.$message.error('清空失败: ' + error.message);
        });
    },
    handleChange(file, fileList) {
      this.fileList = fileList;
      console.log('Selected Files:', fileList); // 打印文件列表
    },
    submitUpload() {
      const formData = new FormData();
      this.fileList.forEach(file => {
        formData.append('file', file.raw);
      });
      console.log('FormData:', formData); // 打印 FormData
      console.log('FormData Entries:', Array.from(formData.entries())); // 打印 FormData 的所有条目

      this.$api.post(`/rooms/${this.roomID}/upload`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
        .then(response => {
          console.log('Upload Response:', response);
          this.$message.success('文件上传成功');
          this.fetchFiles();
          this.$refs.upload.clearFiles();
        })
        .catch(error => {
          this.$message.error('上传失败: ' + error.message);
        });
    },
    fetchFiles() {
      this.$api.get(`/rooms/${this.roomID}/files`)
        .then(response => {
          this.files = response.data.files.map(file => ({
            ...file,
            url: `http://10.20.10.241:8088${file.url}`
          }));
          console.log('Fetched Files:', this.files); // 打印文件列表
        })
        .catch(error => {
          this.$message.error('获取文件列表失败: ' + error.message);
        });
    },
    previewImage(url) {
      const closeAlert = () => {
        this.$msgbox.close();
        document.removeEventListener('keydown', this.closeOnEsc);
      };

      this.closeOnEsc = (event) => {
        if (event.key === 'Escape') {
          closeAlert();
        }
      };

      this.$alert(`<img src="${url}" style="width: 100%; height: auto;" />`, '预览', {
        dangerouslyUseHTMLString: true,
        showConfirmButton: false,
        showCancelButton: true,
        cancelButtonText: '关闭',
        callback: action => {
          if (action === 'cancel') {
            closeAlert();
          }
        }
      });

      document.addEventListener('keydown', this.closeOnEsc);
    },
    downloadFile(url, name) {
      const link = document.createElement('a');
      link.href = url;
      link.download = name;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },
    updateTitle() {
      if (this.activeTab === 'clipboard') {
        this.currentTitle = 'AirTransfer--剪贴板共享';
      } else if (this.activeTab === 'files') {
        this.currentTitle = 'AirTransfer--文件共享';
      }
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
  max-width: 100%;
  box-sizing: border-box;
}

.clipboard-contents {
  height: 200px;
  overflow-y: auto;
  border: 1px solid #ccc;
  padding: 10px;
  margin-bottom: 10px;
}

.file-list {
  display: flex;
  flex-direction: column;
}

.file-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.file-preview {
  width: 50px;
  height: 50px;
  margin-right: 10px;
  cursor: pointer;
}

@media (max-width: 480px) {
  .app-container {
    padding: 5px;
  }

  h1 {
    font-size: 20px;
  }

  .el-button {
    font-size: 12px;
  }

  /* 确保房间信息输入框自适应宽度 */
  .el-form-item__content {
    width: 100%;
  }

  /* 按钮水平排列 */
  .el-form-item__content .el-button {
    width: 48%;
    margin: 0 1%;
  }

  /* 文件共享部分的按钮布局 */
  .el-upload .el-button {
    width: 48%;
    margin: 0 1%;
  }
}
</style>