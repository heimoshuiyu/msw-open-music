const component_search_folders = {
	emits: ['play_audio', 'set_token'],
	data() {
		return {
			search_foldernames: "",
			folders: [],
			folder: {},
			offset: 0,
			limit: 10,
			folder_offset: 0,
			folder_limit: 10,
			files_in_folder: [],
			playing_audio_file: {},
		}
	},
	template: `
<div class="search_toolbar">
<input type="text" v-model="search_foldernames" placeholder="Enter folder name" />
<button @click="first_search_folders">Search Folders</Button>
<button @click="last_page">Last Page</button>
<span>{{ offset }}~{{ offset + folders.length }}</span>
<button @click="next_page">Next Page</button>
</div>

<table>
<thead>
<tr>
	<th>Folder Name</th>
	<th>Action</th>
</tr>
</thead>
<tbody>
<tr v-for="folder in folders">
	<td>{{ folder.foldername }}</td>
	<td><button @click="view_folder(folder)">View</button></td>
</tr>
</tbody>
</table>

<div class="search_toolbar">
<button @click="folder_last_page">Last Page</button>
<span>{{ folder_offset }}~{{ folder_offset + files_in_folder.length }}</span>
<button @click="folder_next_page">Next Page</button>
</div>
<table>
	<thead>
		<tr>
			<th>Filename</th>
			<th>Folder Name</th>
			<th>Size</th>
			<th>Action</th>
		</tr>
	</thead>
	<tbody>
		<tr v-for="file in files_in_folder">
			<component-file :file=file @play_audio="$emit('play_audio', $event)"></component-file>
		</tr>
	</tbody>
</table>
`,
	mounted() {
		if (this.$route.query.folder_id) {
			this.folder.id = parseInt(this.$route.query.folder_id)
			this.folder.foldername = this.$route.query.foldername
			this.get_files_in_folder()
		}
	},
	methods: {
		folder_last_page() {
			this.folder_offset = this.folder_offset - this.folder_limit
			if (this.folder_offset < 0) {
				this.folder_offset = 0
				return
			}
			this.get_files_in_folder()
		},
		folder_next_page() {
			this.folder_offset = this.folder_offset + this.folder_limit
			this.get_files_in_folder()
		},
		view_folder(folder) {
			this.folder = folder
			this.get_files_in_folder()
		},
		get_files_in_folder() {
			axios.post('/api/v1/get_files_in_folder', {
				folder_id: this.folder.id,
				limit: this.folder_limit,
				offset: this.folder_offset,
			}).then((response) => {
				var files = response.data.files
				for (var key in files) {
					files[key].foldername = this.folder.foldername
				}
				this.files_in_folder = files
			})
		},
		last_page() {
			this.offset = this.offset - this.limit
			if (this.offset < 0) {
				this.offset = 0
				return
			}
			this.search_folders()
		},
		next_page() {
			this.offset = this.offset + this.limit
			this.search_folders()
		},
		first_search_folders() {
			this.offset = 0
			this.search_folders()
		},
		search_folders() {
			axios.post('/api/v1/search_folders', {
				foldername: this.search_foldernames,
				limit: this.limit,
				offset: this.offset,
			}).then((response) => {
				this.folders = response.data.folders
			})
		},
	},
}

const component_token = {
	progs: ['token'],
	emits: ['set_token'],
	data() {
		return {
			token_tmp: "",
		}
	},
	template: `
<table><tbody><tr>
<td>Token</td>
<td><input type="text" v-model="token_tmp" @change="emit_set_token" placeholder="token" /></td>
</tr></tbody></table>
`,
	methods: {
		emit_set_token() {
			this.$emit('set_token', this.token_tmp)
		},
	},
}

const component_manage= {
	props: ['token'],
	emits: ['set_token'],
	data() {
		return {
			root: "",
			pattern: [".flac", ".mp3"],
			pattern_tmp: "",
			s: "",
		}
	},
	template: `
<div>
<h4>关于本站</h4>
<p>本站是 MSW Project 的一个应用，希望以个人之力分享被隐藏在历史中的音乐。</p>
<p>自己是V家厨，喜欢的p主包括 wonder-k, buzzG, *luna 等，但却因为种种原因淹没在主流音乐APP的曲库中。本站的初衷是为了让那些知名度低的 VOCALOID / ACG / 东方曲，能够被更多有缘人听到，同时有一个跨平台的工具，能够在低网速的条件下享受硬盘中的无损音乐。</p>
<component-token :token="token" @set_token="$emit('set_token', $event)"></component-token>
<table>
<tbody>
<tr>
	<td>Root</td>
	<td><input type="text" v-model="root" placeholder="/path/to/root" /></td>
</tr>
<tr>
	<td><button @click="add_pattern">Add Pattern</button></td>
	<td><input type="text" v-model="pattern_tmp" placeholder=".wav" /></td>
</tr>
<tr>
	<td colspan="2"><strong>Pattern List</strong></td>
</tr>
<tr v-for="p in pattern">
	<td colspan="2">{{ p }}</td>
</tr>
<tr>
	<td><button @click="update_database">Update</button></td>
	<td><button @click="reset_database">Reset</button></td>
</tr>
<tr>
	<td>Status</td>
	<td>{{ s }}</td>
</tr>
</tbody>
</table>
</div>
`,
	methods: {
		add_pattern() {
			this.pattern.push(this.pattern_tmp)
			this.pattern_tmp = ""
		},
		reset_database() {
			axios.post('/api/v1/reset', {
				token: this.token,
			}).then((response) => {
				this.s = response.data.status
			}).catch((err) => {
				this.s = err.response.data.status
			})
		},
		update_database() {
			this.s = "Updating..."
			axios.post('/api/v1/walk', {
				token: this.token,
				root: this.root,
				pattern: this.pattern,
			}).then((response) => {
				this.s = response.data.status
			}).catch((err) => {
				this.s = err.response.data.status
			})
		}
	},
}

const component_file_dialog = {
	props: ['file', 'show_dialog'],
	emits: ['play_audio', 'close_dialog'],
	template: `
<dialog open v-if="show_dialog">
	<p>{{ file.filename }}</p>
	<p>Download 使用 Axios 异步下载<br />Play 调用网页播放器播放源文件<br />Stream 将串流播放稍低码率的文件</p>
	<button @click="download_file(file)" :disabled="disabled">{{ computed_download_status }}</button>
	<button @click="emit_play_audio">Play</button>
	<button @click="emit_stream_audio">Stream</button>
	<button @click="emit_close_dialog">Close</button>
</dialog>
	`,
	data() {
		return {
			download_loaded: 0,
			disabled: false,
		}
	},
	methods: {
		emit_close_dialog() {
			this.$emit('close_dialog')
		},
		emit_stream_audio() {
			this.file.play_back_type = 'stream',
			this.$emit("play_audio", this.file)
			this.emit_close_dialog()
		},
		emit_play_audio() {
			console.log("pressed button")
			this.file.play_back_type = 'raw'
			this.$emit("play_audio", this.file)
			this.emit_close_dialog()
		},
		download_file(file) {
			this.disabled = true
			axios({
				url: '/api/v1/get_file',
				method: 'POST',
				responseType: 'blob', // important
				data: {
					id: file.id,
				},
				onDownloadProgress: ProgressEvent => {
					this.download_loaded = ProgressEvent.loaded
				}
			}).then((response) => {
				const url = window.URL.createObjectURL(new Blob([response.data]));
				const link = document.createElement('a');
				link.href = url;
				link.setAttribute('download', file.filename);
				document.body.appendChild(link);
				link.click();
				this.download_loaded = 0
				this.disabled = false
				this.emit_close_dialog()
			})
		},
	},
	computed: {
		computed_download_status() {
			if (this.download_loaded === 0) {
				return 'Download'
			} else {
				return Math.round(this.download_loaded / this.file.filesize * 100) + '%'
			}
		},
	},
}

const component_file = {
	props: ['file'],
	emits: ['play_audio'],
	template: `
<td @click="dialog">{{ file.filename }}</td>
<td @click="show_folder">{{ file.foldername }}</td>
<td>{{ computed_readable_size }}</td>
<td>
	<button @click="dialog">Dialog</button>
	<component-file-dialog
		@close_dialog="close_dialog"
		@play_audio="$emit('play_audio', $event)"
		:show_dialog="show_dialog"
		:file="file"
	></component-file-dialog>
</td>
`,
	data() {
		return {
			download_loaded: 0,
			disabled: false,
			show_dialog: false,
		}
	},
	methods: {
		show_folder() {
			this.$router.push({
				path: '/search_folders',
				query: {
					folder_id: this.file.folder_id,
					foldername: this.file.foldername,
				}
			})
		},
		close_dialog() {
			this.show_dialog = false
		},
		dialog() {
			this.show_dialog = this.show_dialog ? false : true
		},
	},
	computed: {
		computed_readable_size() {
			let filesize = this.file.filesize
			if (filesize < 1024) {
				return filesize
			}
			if (filesize < 1024 * 1024) {
				return Math.round(filesize / 1024) + 'K'
			}
			if (filesize < 1024 * 1024 * 1024) {
				return Math.round(filesize / 1024 / 1024) + 'M'
			}
			if (filesize < 1024 * 1024 * 1024 * 1024) {
				return Math.round(filesize / 1024 / 1024 / 1024) + 'G'
			}
		},
	},
}

const component_audio_player = {
	data() {
		return {
		}
	},
	props: ["file"],
	template: `
<div v-if="computed_show">
<span>{{ file.filename }} / {{ file.foldername }}</span>
<video class="audio-player" :src="computed_playing_audio_file_url" controls autoplay>
</video>
</div>
`,
	computed: {
		computed_playing_audio_file_url() {
			if (this.file.play_back_type === 'raw') {
				return '/api/v1/get_file_direct?id=' + this.file.id
			} else if (this.file.play_back_type === 'stream') {
				return '/api/v1/get_file_stream?id=' + this.file.id
			}
		},
		computed_show() {
			return this.file.id ? true : false
		},
	},
}

const component_search_files = {
	emits: ['play_audio'],
	template: `
<div>
<input type="text" name="filename" v-model="search_filenames" placeholder="Enter filename" />
<button @click="first_search_files">Search</button>
<button @click="last_page">Last Page</button>
<span>{{ offset }}~{{ offset + files.length }}</span>
<button @click="next_page">Next Page</button>
<table>
	<thead>
		<tr>
			<th>Filename</th>
			<th>Folder Name</th>
			<th>Size</th>
			<th>Action</th>
		</tr>
	</thead>
	<tbody>
		<tr v-for="file in files">
			<component-file :file=file @play_audio="$emit('play_audio', $event)"></component-file>
		</tr>
	</tbody>
</table>
</div>
`,
	data() {
		return {
			search_filenames: '',
			files: [],
			offset: 0,
			limit: 10,
			playing_audio_file: {},
		}
	},
	methods: {
		first_search_files() {
			this.offset = 0
			this.search_files()
		},
		search_files() {
			axios.post('/api/v1/search_files', {
				filename: this.search_filenames,
				limit: this.limit,
				offset: this.offset,
			}).then((response) => {
				this.files = response.data.files
			})
		},
		last_page() {
			this.offset = this.offset - this.limit
			if (this.offset < 0) {
				this.offset = 0
				return
			}
			this.search_files()
		},
		next_page() {
			this.offset = this.offset + this.limit
			this.search_files()
		},
	},
}

const component_get_random_files = {
	emits: ['play_audio', 'set_token'],
	data() {
		return {
			files: [],
		}
	},
	template: `
<button @click="get_random_files">Refresh</button>
<table>
	<thead>
		<tr>
			<th>Filename</th>
			<th>Folder Name</th>
			<th>Size</th>
			<th>Action</th>
		</tr>
	</thead>
	<tbody>
		<tr v-for="file in files">
			<component-file :file=file @play_audio="$emit('play_audio', $event)"></component-file>
		</tr>
	</tbody>
</table>
`,
	mounted() {
		this.get_random_files()
	},
	methods: {
		get_random_files() {
			axios.get('/api/v1/get_random_files'
			).then(response => {
				this.files = response.data.files;
			})
		}
	},
}

const routes = [
	{ path: '/', component: component_get_random_files},
	{ path: '/search_files', component: component_search_files},
	{ path: '/search_folders', component: component_search_folders},
	{ path: '/manage', component: component_manage},
]
const router = VueRouter.createRouter({
	history: VueRouter.createWebHashHistory(),
	routes,
})

const app = Vue.createApp({
	data() {
		return {
			playing_audio_file: {},
			token: "default token",
		}
	},
	methods: {
		set_token(token) {
			this.token = token
		},
		play_audio(file) {
			console.log(file)
			this.playing_audio_file = file
		},
	},
})

app.component('component-search-folders', component_search_folders)
app.component('component-manage', component_manage)
app.component('component-file', component_file)
app.component('component-audio-player', component_audio_player)
app.component('component-search-files', component_search_files)
app.component('component-get-random-files', component_get_random_files)
app.component('component-file-dialog', component_file_dialog)
app.component('component-token', component_token)

app.use(router)

app.mount('#app')
