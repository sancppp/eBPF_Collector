<template>
  <div id="app">
    <h1>mobisys容器云团队-容器异常检测系统示例</h1>

    <div class="sections">
      <!-- <h4>被观测系统</h4> -->
      <iframe src="被观测系统.drawio.html" frameborder="0"></iframe>


      <div class="section">
        <h2>PUT 功能</h2>
        <input v-model="putKey" placeholder="Key">
        <input v-model="putValue" placeholder="Value">
        <button @click="putData">PUT</button>
      </div>

      <div class="section">
        <h2>GET 功能</h2>
        <input v-model="getKey" placeholder="Key">
        <button @click="getData">GET</button>
      </div>

      <div class="section cmd-section">
        <h2>CMD 功能 <i class="fas fa-user-secret"></i></h2>
        <input v-model="cmd" placeholder="CMD">
        <button @click="cmdData">CMD</button>
      </div>

      <div class="section db-section">
        <h2>DB 功能</h2>
        <button @click="dbData">DB</button>
      </div>
    </div>

    <div class="response" v-if="response">
      <h2>Response:</h2>
      <pre>{{ formattedResponse }}</pre>
    </div>

    <div class="kafka-events">
      <h2>容器异常行为事件</h2>
      <ul>
        <li v-for="(event, index) in kafkaEvents" :key="index" class="event-item">
          {{ event }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      putKey: '',
      putValue: '',
      getKey: '',
      cmd: '',
      response: null,
      kafkaEvents: []
    };
  },
  created() {
    this.setupWebSocket();
  },
  methods: {
    async putData() {
      try {
        const res = await axios.get('http://192.168.252.131:8080/put', {
          params: { key: this.putKey, value: this.putValue },
        });
        this.response = res.data;
      } catch (error) {
        this.response = error.message;
      }
    },
    async getData() {
      try {
        const res = await axios.get('http://192.168.252.131:8080/get', {
          params: { key: this.getKey },
        });
        this.response = res.data;
      } catch (error) {
        this.response = error.message;
      }
    },
    async cmdData() {
      try {
        const res = await axios.get('http://192.168.252.131:8080/cmd', {
          params: { cmd: this.cmd },
        });
        this.response = res.data;
      } catch (error) {
        this.response = error.message;
      }
    },
    async dbData() {
      try {
        const res = await axios.get('http://192.168.252.131:8080/db');
        this.response = res.data;
      } catch (error) {
        this.response = error.message;
      }
    },
    setupWebSocket() {
      const ws = new WebSocket('ws://192.168.252.128:8090');
      console.log('WebSocket client created');

      ws.onopen = () => {
        console.log('WebSocket connection opened');
      };

      ws.onmessage = (event) => {
        console.log('WebSocket message received:', event.data);
        const data = JSON.parse(event.data);
        this.kafkaEvents.unshift(data); // 将新事件插入到数组的头部
        this.animateEvent();
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      ws.onclose = () => {
        console.log('WebSocket connection closed');
      };
    },
    animateEvent() {
      this.$nextTick(() => {
        const eventItems = document.querySelectorAll('.event-item');
        const firstEventItem = eventItems[0]; // 获取第一个事件项
        if (firstEventItem) {
          firstEventItem.classList.add('highlight');
          setTimeout(() => {
            firstEventItem.classList.remove('highlight');
          }, 2000);
        }
      });
    }
  },
  computed: {
    formattedResponse() {
      return typeof this.response === 'object'
        ? JSON.stringify(this.response, null, 2)
        : this.response;
    }
  }
};
</script>

<style>
@import '~@fortawesome/fontawesome-free/css/all.css';

body {
  font-family: Arial, sans-serif;
  text-align: center;
}

h1 {
  margin-top: 50px;
}

.sections {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.section {
  margin: 10px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.section input {
  display: block;
  margin: 10px 0;
  padding: 10px;
  font-size: 16px;
  width: calc(100% - 24px);
  box-sizing: border-box;
}

.section button {
  padding: 10px 20px;
  font-size: 16px;
}

.iframe-section {
  flex: 1 1 1%;
  /* 占据一整行 */
  text-align: left;
  width: 1 cap;
}

.iframe-section iframe {
  width: 100%;
  height: 30vh;
  /* 根据需要调整高度 */
  border: 1px solid #ccc;
  border-radius: 5px;
}

.response {
  margin-top: 20px;
  text-align: left;
  width: 80%;
  margin-left: 10%;
  max-height: calc(8 * 1.5em);
  /* 20 行，高度根据字体大小调整 */
  overflow-y: auto;
  /* 使内容可滚动 */
  border: 1px solid #ccc;
  /* 添加边框以区分 */
  padding: 10px;
  /* 添加内边距 */
  background-color: #f8f8f8;
  /* 添加背景色 */
  border-radius: 5px;
  /* 圆角边框 */
}

pre {
  background-color: #f8f8f8;
  padding: 10px;
  border-radius: 5px;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.cmd-section {
  background-color: #f0f0f0;
  border: 2px solid #000;
  box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.75);
}

.cmd-section h2 {
  display: flex;
  align-items: center;
  justify-content: center;
}

.cmd-section h2 i {
  margin-left: 10px;
  color: #000;
}

.db-section {
  background-color: #e0f0e0;
  border: 2px solid #008000;
  box-shadow: 0px 0px 10px 0px rgba(0, 128, 0, 0.75);
}

.db-section h2 {
  display: flex;
  align-items: center;
  justify-content: center;
}

.db-section h2 i {
  margin-left: 10px;
  color: #008000;
}

.kafka-events {
  margin-top: 20px;
  text-align: left;
  width: 80%;
  margin-left: 10%;
}

.kafka-events ul {
  list-style: none;
  padding: 0;
}

.kafka-events .event-item {
  font-size: 16px;
  padding: 10px;
  border-bottom: 1px solid #ccc;
}

.kafka-events .event-item.highlight {
  font-weight: bold;
  color: red;
  animation: highlight-animation 2s ease-in-out;
}

@keyframes highlight-animation {
  0% {
    transform: scale(1);
  }

  50% {
    transform: scale(1.1);
  }

  100% {
    transform: scale(1);
  }
}
</style>
