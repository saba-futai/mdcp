<template>
  <!-- 绑定 dark-theme 类 -->
  <div class="app-container" :class="{ 'dark-theme': isDarkMode, 'light-theme': !isDarkMode }">
    <!-- 顶部工具栏 -->
    <header class="header">
      <h1 class="title">
        <div class="icon-box">🎵</div>
        RHYTHM GAME PICKER
      </h1>
      <div class="header-actions">
        <!-- 主题切换按钮 -->
        <button @click="toggleTheme" class="btn btn-secondary theme-toggle" :title="isDarkMode ? '切换浅色' : '切换深色'">
          {{ isDarkMode ? '☀️' : '🌙' }}
        </button>
        <button @click="openAddModal" class="btn btn-primary">
          <span>➕</span> <span class="hide-mobile">添加歌曲</span>
        </button>
        <button @click="openRandomModal" class="btn btn-accent">
          <span>🎲</span> <span class="hide-mobile">随机抽取</span>
        </button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-bar-wrapper">
      <div class="filter-bar">
        <div class="filter-group search-group">
          <div class="input-wrapper">
            <span class="search-icon">🔍</span>
            <input type="text" v-model="searchTerm" placeholder="SEARCH SONGS..." class="neo-input">
          </div>
        </div>
        <div class="filter-group">
          <label class="filter-label">SORT BY:</label>
          <div class="select-wrapper">
            <select v-model="sortType" class="neo-select">
              <option value="title_asc">A-Z 歌名</option>
              <option value="title_desc">Z-A 歌名</option>
              <option value="level_desc">Lv. 高 -> 低</option>
              <option value="level_asc">Lv. 低 -> 高</option>
            </select>
          </div>
        </div>
        <div class="filter-group">
          <label class="checkbox-label neo-checkbox">
            <input type="checkbox" v-model="showFavoritesOnly">
            <span class="checkmark"></span>
            ONLY FAVORITES
          </label>
        </div>
      </div>
    </div>

    <!-- 歌曲网格 -->
    <main class="main-content">
      <transition-group name="list-anim" tag="div" class="songs-grid">
        <div
            v-for="song in visibleSongs"
            :key="song.id"
            class="song-card"
            :style="{ backgroundColor: song.color }"
        >
          <!-- 卡片内容 -->
          <div class="card-header">
             <span class="level-badge">LV.{{ song.level || '?' }}</span>
             <button @click.stop="toggleFavorite(song.id)" class="favorite-btn" :class="{ active: song.isFavorite }">
              {{ song.isFavorite ? '★' : '☆' }}
            </button>
          </div>
          
          <div class="card-body">
            <h3 class="song-title">{{ song.title }}</h3>
            <p class="song-artist">{{ song.artist }}</p>
          </div>

          <div class="card-footer">
            <div class="action-buttons">
               <button @click.stop="toggleBlacklist(song.id)" class="neo-icon-btn danger" :class="{ active: song.isBlacklisted }" title="黑名单">
                 {{ song.isBlacklisted ? '🚫' : '⭕' }}
               </button>
               <button @click.stop="editSong(song)" class="neo-icon-btn" title="编辑">✏️</button>
               <button @click.stop="deleteSong(song)" class="neo-icon-btn" title="删除">🗑️</button>
            </div>
          </div>
        </div>
      </transition-group>

      <div ref="loadMoreTrigger" v-show="hasMoreSongs" class="load-more-trigger"></div>

      <div v-if="visibleSongs.length === 0 && !loading" class="empty-state">
        <div class="empty-icon-box"></div>
        <p>NO SONGS FOUND.</p>
      </div>
    </main>

    <!-- 添加/编辑歌曲模态框 -->
    <transition name="neo-modal">
      <div v-if="showAddModal" class="modal-overlay" @click.self="closeAddModal">
        <div class="modal-content neo-box">
          <div class="modal-header">
            <h2>{{ editingItem ? 'EDIT SONG' : 'NEW SONG' }}</h2>
            <button @click="closeAddModal" class="close-btn">✖</button>
          </div>
          <form @submit.prevent="saveSong" class="form">
            <div class="form-group">
              <label>TITLE *</label>
              <input v-model="formData.title" type="text" required class="neo-input">
            </div>
            <div class="form-group">
              <label>ARTIST / PACK</label>
              <input v-model="formData.artist" type="text" class="neo-input">
            </div>
            <div class="form-group">
              <label>LEVEL</label>
              <input v-model.number="formData.level" type="number" min="0" max="20" step="0.1" class="neo-input">
            </div>
            <div class="form-group">
              <label>COLOR</label>
              <div class="color-picker">
                <div
                    v-for="color in macaronColors"
                    :key="color"
                    class="color-swatch"
                    :style="{ backgroundColor: color }"
                    :class="{ selected: formData.color === color }"
                    @click="formData.color = color"
                ></div>
              </div>
            </div>
             <div class="form-group">
              <label class="checkbox-label neo-checkbox">
                <input type="checkbox" v-model="formData.isBlacklisted">
                <span class="checkmark"></span>
                添加到黑名单 (SKIP IN RANDOM)
              </label>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary full-width">SAVE RECORD</button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- 随机抽取模态框 -->
    <transition name="neo-modal">
      <div v-if="showRandomModal" class="modal-overlay" @click.self="showRandomModal = false">
        <div class="modal-content neo-box">
           <div class="modal-header">
            <h2>🎲 RANDOMIZER</h2>
            <button @click="showRandomModal = false" class="close-btn">✖</button>
          </div>

          <div v-if="searchTerm.trim() !== ''" class="info-tag">
            FILTERING FROM <strong>{{ filteredSongs.length }}</strong> SONGS
          </div>

          <form @submit.prevent="performRandomPick" class="form">
            <div class="form-group">
              <label>COUNT</label>
              <div class="stepper-input">
                <button type="button" @click="changeRandomCount(-1)" class="stepper-btn">-</button>
                <input v-model.number="randomOptions.count"  min="1" max="10" required class="neo-input no-border-radius">
                <button type="button" @click="changeRandomCount(1)" class="stepper-btn">+</button>
              </div>
            </div>
            <div class="form-group">
              <label>LEVEL RANGE</label>
              <div class="range-group">
                <input v-model.number="randomOptions.minLevel" type="number" min="0" placeholder="MIN" class="neo-input">
                <span class="divider">TO</span>
                <input v-model.number="randomOptions.maxLevel" type="number" max="20" placeholder="MAX" class="neo-input">
              </div>
            </div>
            <div class="form-group">
              <label class="checkbox-label neo-checkbox">
                <input type="checkbox" v-model="randomOptions.onlyFavorites">
                <span class="checkmark"></span>
                ONLY FAVORITES
              </label>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-accent full-width">LET'S ROLL</button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- 抽取结果模态框 -->
    <transition name="neo-modal">
      <div v-if="showResultModal" class="modal-overlay" @click.self="showResultModal = false">
        <div class="modal-content neo-box result-modal">
          <div class="modal-header center-header">
            <h2>🎉 RESULTS 🎉</h2>
          </div>
          <div class="result-grid">
            <div
                v-for="(song, index) in randomResults"
                :key="song.id"
                class="result-card"
                :style="{ backgroundColor: song.color }"
            >
              <div class="result-number">#{{ index + 1 }}</div>
              <h3>{{ song.title }}</h3>
              <p>{{ song.artist }}</p>
              <div class="result-badge">LV.{{ song.level }}</div>
            </div>
          </div>

          <div v-if="randomResults.length === 0" class="empty-state-small">
            <p style="margin-bottom: 3em;">NOTHING FOUND.</p>
          </div>

          <div class="form-actions grid-2">
            <button @click="showResultModal = false" class="btn btn-secondary">CLOSE</button>
            <button @click="performRandomPick" class="btn btn-accent">REROLL</button>
          </div>
        </div>
      </div>
    </transition>

    <!-- 删除确认模态框 -->
    <transition name="neo-modal">
      <div v-if="showDeleteConfirmModal" class="modal-overlay" @click.self="closeDeleteConfirmModal">
        <div class="modal-content neo-box warning-box">
          <h2>⚠️ DELETE?</h2>
          <p>Are you sure you want to delete <strong>"{{ songToDeleteTitle }}"</strong>?</p>
          <div class="form-actions grid-2">
            <button type="button" @click="closeDeleteConfirmModal" class="btn btn-secondary">NO</button>
            <button type="button" @click="confirmDelete" class="btn btn-danger">YES, DELETE</button>
          </div>
        </div>
      </div>
    </transition>

    <!-- 深色模式全局遮罩 -->
     <div v-if="isDarkMode" class="dark-mode-overlay"></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, onUnmounted } from 'vue'
// 模拟后端方法，实际使用时请确保 Wails 绑定正确
import { GetAllSongs, AddSong, UpdateSong, DeleteSong as BackendDeleteSong, ToggleFavorite, ToggleBlacklist, RandomPick } from '../wailsjs/go/main/App'

// --- Mock Data (如需测试，请取消注释并注释掉上方的 Import) ---
// const GetAllSongs = async () => Array.from({ length: 12 }, (_, i) => ({ id: i, title: `Dark Anthem ${i}`, artist: 'Void Composer', level: 10 + i%5, color: '#FFB3BA', isFavorite: false }));
// const AddSong = async () => {}; const UpdateSong = async () => {}; const BackendDeleteSong = async () => {}; 
// const ToggleFavorite = async () => {}; const ToggleBlacklist = async () => {}; const RandomPick = async () => [];
// ---------------------------------------------

// Dark Mode State
const isDarkMode = ref(false);

// 从本地存储读取主题偏好
onMounted(() => {
  const savedTheme = localStorage.getItem('theme');
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDarkMode.value = true;
  }
});

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value;
  localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light');
};

// ... (原有逻辑保持不变) ...
const searchTerm = ref('');
const allSongs = ref([])
const visibleSongs = ref([])
const loadMoreTrigger = ref(null)
const loading = ref(true);
const PAGE_SIZE = 20;
let observer = null;
const sortType = ref('title_asc')
const showFavoritesOnly = ref(false)
const showAddModal = ref(false)
const showRandomModal = ref(false)
const showResultModal = ref(false)
const editingItem = ref(null)
const randomResults = ref([])
const showDeleteConfirmModal = ref(false)
const songToDeleteId = ref(null)
const songToDeleteTitle = ref('')
const macaronColors = [
  '#FFB3BA', '#FFDFBA', '#FFFFBA', '#BAFFC9', '#BAE1FF',
  '#E0BBE4', '#FFB7D5', '#C7CEEA', '#B4E7F5', '#FED9B7',
  '#A0E7E5', '#F7C6C7', '#C9E4DE', '#FFD1DC', '#E4C1F9',
];

const formData = ref({
  title: '', artist: '', level: 1, color: macaronColors[0], isFavorite: false, isBlacklisted: false
})
const randomOptions = ref({
  count: 1, minLevel: 0, maxLevel: 0, onlyFavorites: false
})

const filteredSongs = computed(() => {
  let result = [...allSongs.value]
  if (searchTerm.value.trim() !== '') {
    const lowerCaseSearch = searchTerm.value.toLowerCase();
    result = result.filter(s =>
        s.title.toLowerCase().includes(lowerCaseSearch) ||
        (s.artist && s.artist.toLowerCase().includes(lowerCaseSearch))
    );
  }
  if (showFavoritesOnly.value) {
    result = result.filter(s => s.isFavorite)
  }
  switch (sortType.value) {
    case 'title_asc': result.sort((a, b) => a.title.localeCompare(b.title)); break;
    case 'title_desc': result.sort((a, b) => b.title.localeCompare(a.title)); break;
    case 'level_desc': result.sort((a, b) => (b.level || 0) - (a.level || 0)); break;
    case 'level_asc': result.sort((a, b) => (a.level || 0) - (b.level || 0)); break;
  }
  return result
})
const hasMoreSongs = computed(() => {
  return visibleSongs.value.length < filteredSongs.value.length
})


watch([filteredSongs], () => {
  resetAndLoadVisibleSongs()
})

async function performRandomPick() {
  try {
    if (searchTerm.value.trim() !== '') {
      let potentialPicks = [...filteredSongs.value].filter(s => !s.isBlacklisted);

      if (randomOptions.value.onlyFavorites) {
        potentialPicks = potentialPicks.filter(s => s.isFavorite);
      }
      if (randomOptions.value.minLevel > 0) {
        potentialPicks = potentialPicks.filter(s => (s.level || 0) >= randomOptions.value.minLevel);
      }
      if (randomOptions.value.maxLevel > 0) {
        potentialPicks = potentialPicks.filter(s => (s.level || 0) <= randomOptions.value.maxLevel);
      }

      for (let i = potentialPicks.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [potentialPicks[i], potentialPicks[j]] = [potentialPicks[j], potentialPicks[i]];
      }

      randomResults.value = potentialPicks.slice(0, randomOptions.value.count);

    } else {
      randomResults.value = await RandomPick(randomOptions.value);
    }

    if (!showResultModal.value) {
      showRandomModal.value = false;
      showResultModal.value = true;
    }

  } catch (err) {
    console.error('抽取失败:', err);
  }
}


function loadNextPage() {
  if (!hasMoreSongs.value) return;
  const currentPage = Math.ceil(visibleSongs.value.length / PAGE_SIZE);
  const start = currentPage * PAGE_SIZE;
  const end = start + PAGE_SIZE;
  const newSongs = filteredSongs.value.slice(start, end);
  visibleSongs.value.push(...newSongs);
}

function resetAndLoadVisibleSongs() {
  visibleSongs.value = filteredSongs.value.slice(0, PAGE_SIZE);
}

function setupObserver() {
  observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          loadNextPage();
        }
      },
      { root: null, rootMargin: '0px', threshold: 0.1 }
  );
  if (loadMoreTrigger.value) {
    observer.observe(loadMoreTrigger.value);
  }
}

onMounted(async () => {
  await loadSongs()
  setupObserver()
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect();
  }
})

async function loadSongs() {
  loading.value = true;
  try {
    allSongs.value = await GetAllSongs() || []
    resetAndLoadVisibleSongs()
  } catch (err) {
    console.error('加载歌曲失败:', err)
  } finally {
    loading.value = false;
  }
}

function openAddModal() { resetForm(); showAddModal.value = true; }
function editSong(song) {
  editingItem.value = { ...song };
  formData.value = {
    id: song.id, title: song.title, artist: song.artist,
    level: song.level || 1, color: song.color || macaronColors[0],
    isFavorite: song.isFavorite, isBlacklisted: song.isBlacklisted || false,
  };
  showAddModal.value = true;
}
async function saveSong() {
  const songToSave = { ...formData.value };
  try {
    if (editingItem.value) {
      await UpdateSong(songToSave);
    } else {
      await AddSong(songToSave);
    }
    // Mock update for UI
    if (editingItem.value) {
        const idx = allSongs.value.findIndex(s => s.id === editingItem.value.id);
        if(idx !== -1) allSongs.value[idx] = songToSave;
    } else {
        songToSave.id = Date.now();
        allSongs.value.push(songToSave);
    }

    await loadSongs();
    closeAddModal();
  } catch (err) { console.error('保存失败:', err); }
}
function deleteSong(song) {
  songToDeleteId.value = song.id;
  songToDeleteTitle.value = song.title;
  showDeleteConfirmModal.value = true;
}
async function confirmDelete() {
  if (songToDeleteId.value === null) return;
  try {
    await BackendDeleteSong(songToDeleteId.value);
    allSongs.value = allSongs.value.filter(s => s.id !== songToDeleteId.value);
    await loadSongs();
  } catch (err) {
    console.error('删除失败:', err);
  } finally {
    closeDeleteConfirmModal();
  }
}
function closeDeleteConfirmModal() {
  showDeleteConfirmModal.value = false;
  songToDeleteId.value = null;
  songToDeleteTitle.value = '';
}
async function toggleFavorite(id) {
  try {
    const song = allSongs.value.find(s => s.id === id);
    if (song) { song.isFavorite = !song.isFavorite; }
    await ToggleFavorite(id);
  } catch (err) { console.error('操作失败:', err); await loadSongs(); }
}

async function toggleBlacklist(id) {
  try {
    const song = allSongs.value.find(s => s.id === id);
    if (song) {
      song.isBlacklisted = !song.isBlacklisted;
    }
    await ToggleBlacklist(id);
  } catch (err) {
    console.error('黑名单操作失败:', err);
    await loadSongs();
  }
}

function closeAddModal() {
  showAddModal.value = false;
  editingItem.value = null;
  resetForm();
}

function resetForm() {
  formData.value = {
    title: '', artist: '', level: 1,
    color: macaronColors[0], isFavorite: false, isBlacklisted: false
  };
}

function openRandomModal() {
  randomOptions.value.onlyFavorites = showFavoritesOnly.value;
  showRandomModal.value = true;
}
function changeRandomCount(delta) {
  const current = randomOptions.value.count;
  const newValue = current + delta;
  if (newValue >= 1 && newValue <= 10) {
    randomOptions.value.count = newValue;
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;700;900&display=swap');

/* --- 全局变量与主题定义 --- */
:root {
  /* 浅色模式 (默认) */
  --bg-color: #fdfbf7;
  --dot-color: #000000;
  --text-color: #000000;
  --border-color: #000000;
  --shadow-color: #000000;
  --surface-color: #FFFFFF; /* 卡片/输入框背景 */
  --input-bg: #FFFFFF;
  --modal-overlay: rgba(0, 0, 0, 0.6);
  --btn-text-color: #000000;
  --card-text-color: #000000; /* 卡片内部始终保持黑色文本，因为卡片背景是浅色 */
  
  /* 尺寸 */
  --border-width: 3px;
  --shadow-offset: 4px;
}

.light-theme {
  /* 浅色模式 (默认) */
  --bg-color: #fdfbf7;
  --dot-color: #000000;
  --text-color: #000000;
  --border-color: #000000;
  --shadow-color: #000000;
  --surface-color: #FFFFFF; /* 卡片/输入框背景 */
  --input-bg: #FFFFFF;
  --modal-overlay: rgba(0, 0, 0, 0.6);
  --btn-text-color: #000000;
  --card-text-color: #000000; /* 卡片内部始终保持黑色文本，因为卡片背景是浅色 */
  
  /* 尺寸 */
  --border-width: 3px;
  --shadow-offset: 4px;
}

/* 深色模式 (Neo Brutalism Dark) */
.dark-theme {
  --bg-color: #121212; /* 深炭色 */
  --dot-color: #2a2a2a; /* 深灰点阵 */
  --text-color: #EEEEEE; /* 灰白字 */
  --border-color: #FFFFFF; /* 高反差白边框 */
  --shadow-color: #444444; /* 灰色硬阴影 */
  --surface-color: #1E1E1E; /* 深色组件背景 */
  --input-bg: #2C2C2C;
  --modal-overlay: rgba(255, 255, 255, 0.1); /* 浅色微遮罩 */
  --btn-text-color: #000000; /* 按钮保持鲜艳，文字保持黑色 */
}

.dark-mode-overlay{
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.2);
  pointer-events: none;
  z-index: 5;
}

* {
  box-sizing: border-box;
}

.app-container {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  color: var(--text-color);
  background-color: var(--bg-color);
  background-image: radial-gradient(var(--dot-color) 1px, transparent 0);
  background-size: 20px 20px;
  font-family: 'Space Grotesk', 'Courier New', sans-serif;
  transition: background-color 0.3s ease, color 0.3s ease;
}

/* --- 头部区域 --- */
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 30px;
  background: var(--surface-color);
  border-bottom: var(--border-width) solid var(--border-color);
  z-index: 10;
  transition: background-color 0.3s, border-color 0.3s;
}

.title {
  font-size: 28px;
  color: var(--text-color);
  display: flex;
  align-items: center;
  gap: 15px;
  margin: 0;
  font-weight: 900;
  letter-spacing: -1px;
  text-transform: uppercase;
}

.icon-box {
  width: 50px;
  height: 50px;
  background: #FF6B6B;
  border: var(--border-width) solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--shadow-color);
  font-size: 24px;
  color: black; /* 图标内部保持黑色 */
}

.header-actions {
  display: flex;
  gap: 15px;
}

@media (max-width: 600px) {
  .title { font-size: 20px; }
  .hide-mobile { display: none; }
  .header { padding: 15px; }
}

/* --- 按钮 --- */
.btn {
  padding: 12px 24px;
  border: var(--border-width) solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  font-weight: 800;
  cursor: pointer;
  transition: all 0.1s;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--shadow-color);
  text-transform: uppercase;
  color: var(--btn-text-color);
}

.btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: calc(var(--shadow-offset) + 2px) calc(var(--shadow-offset) + 2px) 0 var(--shadow-color);
}

.btn:active {
  transform: translate(2px, 2px);
  box-shadow: 0 0 0 var(--shadow-color);
}

.btn-primary { background: #4ADE80; } /* Green */
.btn-accent { background: #FBBF24; } /* Yellow */
.btn-secondary { background: #E5E7EB; } /* Grey */
.btn-danger { background: #F87171; } /* Red */

/* 深色模式下的二级按钮调整 */
.dark-theme .btn-secondary {
  background: #4B5563; 
  color: white;
}

/* --- 筛选栏 --- */
.filter-bar-wrapper {
    padding: 20px 30px;
}

.filter-bar {
  display: flex;
  gap: 20px;
  padding: 20px;
  background: var(--surface-color);
  border: var(--border-width) solid var(--border-color);
  box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--shadow-color);
  flex-wrap: wrap;
  align-items: center;
  border-radius: 4px;
  transition: background-color 0.3s, border-color 0.3s, box-shadow 0.3s;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-label {
    font-weight: 800;
    font-size: 12px;
}

/* --- 输入框与下拉框 --- */
.input-wrapper {
    position: relative;
    width: 100%;
}

.search-icon {
    position: absolute;
    left: 10px;
    top: 50%;
    transform: translateY(-50%);
    z-index: 1;
    filter: grayscale(100%); /* 让图标适应主题 */
}

.neo-input, .neo-select {
  padding: 10px 15px 10px 15px;
  border: var(--border-width) solid var(--border-color);
  border-radius: 0;
  background: var(--input-bg);
  color: var(--text-color);
  font-family: inherit;
  font-weight: 700;
  outline: none;
  transition: box-shadow 0.2s;
  min-width: 200px;
}

.search-group .neo-input {
    padding-left: 35px;
}

.neo-input:focus, .neo-select:focus {
  background: #FEF9C3; /* Focus always yellow for contrast */
  color: black;
  box-shadow: 4px 4px 0 var(--shadow-color);
}

/* --- 复选框 --- */
.neo-checkbox {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  font-weight: 800;
  font-size: 12px;
  user-select: none;
}

.neo-checkbox input {
  display: none;
}

.checkmark {
  height: 20px;
  width: 20px;
  background-color: var(--input-bg);
  border: var(--border-width) solid var(--border-color);
  position: relative;
}

.neo-checkbox input:checked ~ .checkmark {
  background-color: var(--text-color);
}

.neo-checkbox input:checked ~ .checkmark:after {
  content: "";
  position: absolute;
  left: 5px;
  top: 1px;
  width: 5px;
  height: 10px;
  border: solid var(--bg-color); /* Checkmark color inverse of text */
  border-width: 0 3px 3px 0;
  transform: rotate(45deg);
}

/* --- 主内容区 --- */
.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 30px;
}

.songs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 25px;
}

/* --- 歌曲卡片 --- */
.song-card {
  position: relative;
  border: var(--border-width) solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  cursor: default;
  transition: all 0.2s;
  box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--shadow-color);
  display: flex;
  flex-direction: column;
  min-height: 200px;
  color: var(--card-text-color); /* 强制卡片内文本为黑色 */
}

.song-card:hover {
  transform: translate(-3px, -3px);
  box-shadow: 8px 8px 0 var(--shadow-color);
}

.card-header {
    padding: 15px;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    border-bottom: var(--border-width) solid rgba(0,0,0,0.1);
}

.card-body {
    padding: 15px;
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.card-footer {
    padding: 10px 15px;
    background: rgba(255,255,255,0.4);
    border-top: var(--border-width) solid rgba(0,0,0,0.1); /* 内部边框保持微妙 */
}

/* 在深色模式下，卡片的边框需要与背景对比，但卡片内部颜色是亮色 */
.dark-theme .song-card {
    border-color: #FFFFFF; /* 白边框 */
}
.dark-theme .card-footer {
     border-top-color: rgba(0,0,0,0.1);
}

.song-title {
  font-size: 20px;
  font-weight: 900;
  margin: 0 0 5px 0;
  color: black; /* Always black inside card */
  line-height: 1.1;
  word-break: break-word;
}

.song-artist {
  font-size: 14px;
  color: black;
  margin: 0;
  font-weight: 600;
  opacity: 0.7;
}

.level-badge {
  background: #000;
  color: #fff;
  padding: 4px 8px;
  font-size: 12px;
  font-weight: bold;
  border-radius: 2px;
}

.favorite-btn {
  background: transparent;
  border: none;
  font-size: 24px;
  color: rgba(0,0,0,0.2);
  cursor: pointer;
  line-height: 1;
  transition: transform 0.2s;
}
.favorite-btn:hover { transform: scale(1.2); }
.favorite-btn.active { color: #000; text-shadow: 2px 2px 0 #fff; }

/* --- 动作按钮 --- */
.action-buttons {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
}

.neo-icon-btn {
    width: 32px;
    height: 32px;
    background: white;
    border: 2px solid black;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.1s;
    box-shadow: 2px 2px 0 black;
    color: black;
}

.neo-icon-btn:hover { transform: translate(-1px, -1px); box-shadow: 3px 3px 0 black; }
.neo-icon-btn:active { transform: translate(1px, 1px); box-shadow: 1px 1px 0 black; }
.neo-icon-btn.danger.active { background: #F87171; }

/* --- 空状态 --- */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-color);
  font-weight: bold;
  font-size: 20px;
}

.empty-icon-box {
    font-size: 60px;
    width: 120px;
    height: 120px;
    background: var(--surface-color);
    border: var(--border-width) solid var(--border-color);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
    box-shadow: 6px 6px 0 var(--shadow-color);
}

/* --- 模态框 --- */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: var(--modal-overlay);
  backdrop-filter: grayscale(100%) contrast(120%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.neo-box {
  background: var(--surface-color);
  border: var(--border-width) solid var(--border-color);
  box-shadow: 10px 10px 0 var(--shadow-color);
  padding: 0;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  color: var(--text-color);
}

.modal-header {
    padding: 20px;
    border-bottom: var(--border-width) solid var(--border-color);
    background: #C4B5FD; /* Light Purple header default */
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: black; /* Header text always black */
}

.dark-theme .modal-header {
    background: #6D28D9; /* Darker purple for dark mode */
    color: white;
}

.modal-header h2 {
    margin: 0;
    font-weight: 900;
    text-transform: uppercase;
}

.close-btn {
    background: white;
    border: 2px solid black;
    width: 30px;
    height: 30px;
    cursor: pointer;
    font-weight: bold;
    box-shadow: 2px 2px 0 black;
    color: black;
}
.close-btn:active { transform: translate(1px, 1px); box-shadow: none; }

.form {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

.form-group label {
    font-weight: 800;
    font-size: 14px;
}

.full-width {
    width: 100%;
    justify-content: center;
}

/* --- 颜色选择器 --- */
.color-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
.color-swatch {
  width: 30px;
  height: 30px;
  border: 2px solid var(--border-color);
  cursor: pointer;
  transition: transform 0.2s;
}
.color-swatch:hover { transform: scale(1.1); }
.color-swatch.selected { position: relative; border-color: var(--text-color); }
.color-swatch.selected::after {
    content: '✔';
    position: absolute;
    top: 50%; left: 50%;
    transform: translate(-50%, -50%);
    font-size: 12px;
    color: black;
}

/* --- Stepper & Range --- */
.stepper-input { display: flex; }
.stepper-btn {
    width: 40px;
    background: #eee;
    border: var(--border-width) solid var(--border-color);
    font-weight: bold;
    font-size: 18px;
    cursor: pointer;
    color: black;
}
.dark-theme .stepper-btn { background: #444; color: white; }
.stepper-btn:first-child { border-right: none; }
.stepper-btn:last-child { border-left: none; }
.no-border-radius { border-radius: 0; text-align: center; width: 60px; }

.range-group { display: flex; align-items: center; gap: 10px; }
.divider { font-weight: bold; }

/* --- 结果页面 --- */
.center-header { justify-content: center; background: #6EE7B7; color: black; }
.dark-theme .center-header { background: #059669; color: white; }

.result-grid {
    padding: 20px;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 20px;
}

.result-card {
    border: var(--border-width) solid var(--border-color); /* 结果卡片使用主题边框 */
    padding: 20px;
    text-align: center;
    position: relative;
    box-shadow: 5px 5px 0 var(--shadow-color);
    color: black; /* 内容始终为黑，因为背景色是 macaronColors */
}
.dark-theme .result-card { border-color: #fff; }

.result-number {
    background: black;
    color: white;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    position: absolute;
    top: -10px;
    left: -10px;
    border: 2px solid white;
}

.result-badge {
    background: white;
    border: 2px solid black;
    display: inline-block;
    padding: 2px 8px;
    font-weight: bold;
    font-size: 12px;
    margin-top: 10px;
    color: black;
}

.grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }

.info-tag {
    background: #FEF3C7;
    border-bottom: 2px solid var(--border-color);
    padding: 10px;
    text-align: center;
    font-size: 14px;
    color: black;
}

.warning-box .modal-header { background: #FCA5A5; color: black; }
.dark-theme .warning-box .modal-header { background: #B91C1C; color: white; }
.warning-box { text-align: center; }

/* --- 动画 --- */
.neo-modal-enter-active, .neo-modal-leave-active { transition: opacity 0.2s; }
.neo-modal-enter-from, .neo-modal-leave-to { opacity: 0; }
.neo-modal-enter-active .modal-content { animation: modal-pop 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275); }
.neo-modal-leave-active .modal-content { animation: modal-pop 0.2s reverse; }

@keyframes modal-pop {
    0% { transform: scale(0.8); opacity: 0; }
    100% { transform: scale(1); opacity: 1; }
}

.list-anim-enter-active, .list-anim-leave-active { transition: all 0.3s ease; }
.list-anim-enter-from, .list-anim-leave-to { opacity: 0; transform: translateY(30px); }

.load-more-trigger { height: 20px; }
</style>