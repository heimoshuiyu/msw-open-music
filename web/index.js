const component_share = {
	emits: ['play_audio', 'set_token'],
	props: ['token'],
	template: `
<div class="page">
<h3>Share with others!</h3>
<p v-if="error_status">{{ error_status }}</p>
<p>Share link: <a :href="computed_share_link">{{ computed_share_link }}</a> , or share this page directly.</p>
<table>
<thead>
<tr>
	<th>Filename</th>
	<th>Folder Name</th>
	<th>Size</th>
</tr>
</thead>
<tbody>
<tr>
	<component-file :file=file @play_audio="$emit('play_audio', $event)"></component-file>
</tr>
</tbody>
</table>
</div>
`,
	computed: {
		computed_share_link() {
			return window.location.href
		},
	},
	data() {
		return {
			file: {},
			error_status: "",
		}
	},
	mounted() {
		if (this.$route.query.id) {
			this.get_file_info()
		}
	},
	methods: {
		get_file_info() {
			axios.post('/api/v1/get_file_info', {
				id: parseInt(this.$route.query.id),
			}).then((response) => {
				this.file = response.data
			}).catch((error) => {
				if (error.response) {
					this.error_status = error.response.data.status
				} else {
					this.error_status = 'Network error'
				}
			})
		},
	},
}

const component_search_folders = {
	emits: ['play_audio', 'set_token'],
	props: ['token'],
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
			is_loading: false,
			error_status: "",
		}
	},
	computed: {
		computed_folders_page() {
			if (this.is_loading) {
				return 'Loading...'
			}
			if (this.error_status) {
				return this.error_status
			}
			return this.offset + ' ~ ' + (this.offset + this.folders.length)
		},
		computed_files_page() {
			if (this.is_loading) {
				return 'Loading...'
			}
			if (this.error_status) {
				return this.error_status
			}
			return this.folder_offset + ' ~ ' + (this.folder_offset + this.files_in_folder.length)
		},
	},
	template: `
<div class="page">
<h3>Search Folders</h3>
<div class="search_toolbar">
<input type="text" @keyup.enter="first_search_folders" v-model="search_foldernames" placeholder="Enter folder name" />
<button @click="first_search_folders">Search Folders</Button>
<button @click="last_page">Last Page</button>
<button disabled>{{ computed_folders_page }}</button>
<button @click="next_page">Next Page</button>
</div>

<table v-if="folders.length">
<thead>
<tr>
	<th>Folder Name</th>
	<th>Action</th>
</tr>
</thead>
<tbody>
<tr v-for="folder in folders">
	<td class="clickable" @click="view_folder(folder)">{{ folder.foldername }}</td>
	<td><button @click="view_folder(folder)">View</button></td>
</tr>
</tbody>
</table>

<h3>Files in folder</h3>
<div class="search_toolbar">
<button @click="folder_last_page">Last Page</button>
<button disabled>{{ computed_files_page }}</button>
<button @click="folder_next_page">Next Page</button>
</div>
<table v-if="files_in_folder.length">
	<thead>
		<tr>
			<th>Filename</th>
			<th>Folder Name</th>
			<th>Size</th>
		</tr>
	</thead>
	<tbody>
		<tr v-for="file in files_in_folder">
			<component-file :file=file @play_audio="$emit('play_audio', $event)"></component-file>
		</tr>
	</tbody>
</table>
</div>
`,
	mounted() {
		if (this.$route.query.folder_id) {
			this.folder.id = parseInt(this.$route.query.folder_id)
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
			this.is_loading = true
			axios.post('/api/v1/get_files_in_folder', {
				folder_id: this.folder.id,
				limit: this.folder_limit,
				offset: this.folder_offset,
			}).then((response) => {
				this.error_status = ""
				this.files_in_folder = response.data.files
			}).catch((error) => {
				if (error.response) {
					this.error_status = error.response.data.status
				} else {
					this.error_status = 'Network error'
				}
			}).finally(() => {
				this.is_loading = false
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
			this.is_loading = true
			axios.post('/api/v1/search_folders', {
				foldername: this.search_foldernames,
				limit: this.limit,
				offset: this.offset,
			}).then((response) => {
				this.error_status = ""
				this.folders = response.data.folders
			}).catch((error) => {
				if (error.response) {
					this.error_status = error.response.data.status
				} else {
					this.error_status = 'Network error'
				}
			}).finally(() => {
				this.is_loading = false
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
			feedback: "",
			feedback_status: "Submit",
			feedback_placeholder: "feedback...",
			submit_disabled: false,
			is_err: false,
			err_msg: "",
		}
	},
	template: `
<div class="page">
<div class="description">
<h4>关于本站</h4>
<p>一只随处可见的 葱厨&车万人 想听 TA 屯在硬盘里的音乐。</p>
<p>一点点说明：下方播放器的 Raw 模式即不转码直接播放源文件，支持断点续传；Prepare 模式：勾选后播放的文件将提前在服务器端转码，然后以支持断点续传的方式提供，如果你的网络不稳定，经常播放到一半就中断，可以尝试勾选 Prepare。</p>
<p>站内音乐来自公开网络，仅供个人使用，如有侵权或建议请提交反馈</p>
<div class="feedback">
	<input type="text" v-model="feedback" :disabled="submit_disabled" :placeholder="feedback_placeholder"/>
	<button @click="submit_feedback" :disabled="submit_disabled">{{ feedback_status }}</button>
	<label v-if="is_err">{{ err_msg }}</label>
</div>
</div>
<component-token :token="token" @set_token="$emit('set_token', $event)"></component-token>
<component-manage-database :token="token"></component-manage-database>
</div>
`,
	methods: {
		submit_feedback() {
			axios.post('/api/v1/feedback', {
				feedback: this.feedback,
			}).then((response) => {
				this.submit_disabled = true
				this.feedback = ""
				this.feedback_status = "Success"
				this.feedback_placeholder = "Thanks for your feedback!"
				this.is_err = false
			}).catch((err) => {
				console.log(err)
				this.is_err = true
				this.err_msg = err.response.data.status
			})
		},
	}
}
const component_manage_database = {
	props: ['token'],
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
	<p>
		Download 使用 Axios 异步下载<br />
		Play 调用网页播放器播放<br />
	</p>
	<button @click="download_file(file)" :disabled="disabled">{{ computed_download_status }}</button>
	<button @click="emit_play_audio">Play</button>
	<button @click="share">Share</button>
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
		share() {
			this.$router.push({
				path: '/share',
				query: {
					id: this.file.id,
				},
			})
			this.emit_close_dialog()
		},
		emit_close_dialog() {
			this.$emit('close_dialog')
		},
		emit_play_audio() {
			console.log("pressed button")
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
<td class="clickable" @click="click_filename">{{ file.filename }}</td>
<td class="clickable" @click="show_folder">{{ file.foldername }}</td>
<td>{{ computed_readable_size }}
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
		click_filename() {
			if (this.show_dialog) {
				this.file.play_back_type = 'stream'
				this.$emit('play_audio', this.file)
				this.show_dialog = false
			} else {
				this.show_dialog = true
			}
		},
		show_folder() {
			this.$router.push({
				path: '/search_folders',
				query: {
					folder_id: this.file.folder_id,
				},
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
	emits: ['stop', 'play_audio'],
	data() {
		return {
			loop: true,
			ffmpeg_config: {},
			show_dialog: false,
			is_preparing: false,
			prepare: false,
			raw: false,
			playing_url: "",
			prepared_filesize: 0,
			playing_file: {},
			error_status: "",
		}
	},
	props: ['file', 'token'],
	template: `
<div>
<h5>Player Status</h5>
<component-file-dialog
	@close_dialog="close_dialog"
	@play_audio="$emit('play_audio', $event)"
	:show_dialog="show_dialog"
	:file="file"
></component-file-dialog>
<span v-if="computed_show">
	<button @click="dialog">{{ file.filename }}</button>
	<button @click="show_folder">{{ file.foldername }}</button> 
	<button disabled>{{ computed_readable_size }}</button>
	<button v-if="error_status" @click="retry">Retry</button>
	<button @click="emit_stop">Stop</button>
</span>
<br />
<input type="checkbox" v-model="loop" />
<label>Loop</label>
<input type="checkbox" v-model="raw" />
<label>Raw</label>
<input v-show="!raw" type="checkbox" v-model="prepare" />
<label v-show="!raw">Prepare</label><br />
<video v-if="computed_video_show" class="audio-player" :src="playing_url" controls autoplay :loop="loop">
</video>
<component-stream-config @set_ffmpeg_config="set_ffmpeg_config"></component-stream-config>
<p>{{ token }}</p>
</div>
`,
	methods: {
		emit_stop() {
			this.$emit('stop')
		},
		dialog() {
			this.show_dialog = this.show_dialog ? false : true
		},
		close_dialog() {
			this.show_dialog = false
		},
		show_folder() {
			this.$router.push({
				path: '/search_folders',
				query: {
					folder_id: this.file.folder_id,
				}
			})
		},
		set_ffmpeg_config(ffmpeg_config) {
			this.ffmpeg_config = ffmpeg_config
		},
		prepare_func() {
			if (!this.file.id) {
				return
			}
			this.playing_file = {}
			this.is_preparing = true
			axios.post('/api/v1/prepare_file_stream_direct', {
				id: this.file.id,
				config_name: this.ffmpeg_config.name,
			}).then(response => {
				console.log(response.data)
				this.error_status = ''
				this.prepared_filesize = response.data.filesize
				var file = this.file
				this.playing_file = file
				this.set_playing_url()
				console.log('axios done', this.playing_file)
			}).catch((err) => {
				if (err.response) {
					this.error_status = err.response.data.status
				} else {
					this.error_status = "Network error"
				}
			}).finally(() => {
				this.is_preparing = false
			})
		},
		set_playing_url() {
			if (this.raw) {
				console.log('computed raw rul')
				this.playing_url = '/api/v1/get_file_direct?id=' + this.playing_file.id
			} else {
				if (this.prepare) {
					console.log('empty playing_file, start prepare')
					this.playing_url = '/api/v1/get_file_stream_direct?id=' + this.playing_file.id + '&config=' + this.ffmpeg_config.name
				} else {
					console.log('computed stream url')
					this.playing_url = '/api/v1/get_file_stream?id=' + this.playing_file.id + '&config=' + this.ffmpeg_config.name
				}
			}
		},
		setup_player() {
			// 如果没有勾选 prepare 则直接播放
			// 否则进入 prepare 流程
			this.playing_file = {}
			if (this.prepare && !this.raw) {
				this.prepare_func()
			} else {
				this.playing_file = this.file
				this.set_playing_url()
			}
		},
		retry() {
			this.setup_player()
		},
	},
	watch: {
		file() {
			this.setup_player()
		},
		raw() {
			if (this.prepare) {
				this.prepare_func()
			} else {
				this.set_playing_url()
			}
		},
		prepare() {
			this.playing_file = {}
			this.prepare_func()
		},
		ffmpeg_config() {
			this.setup_player()
		},
	},
	computed: {
		computed_can_retry() {
			return this.error_status ? true : false
		},
		computed_readable_size() {
			if (this.is_preparing) {
				return 'Preparing...'
			}
			if (this.error_status) {
				return this.error_status
			}
			let filesize = this.playing_file.filesize
			if (this.prepare) {
				filesize = this.prepared_filesize
			}
			if (this.raw) {
				filesize = this.playing_file.filesize
			}
			if (filesize < 1024 * 1024 * 1024) {
				filesize = Math.round(filesize / 1024) + 'K'
			}
			if (filesize < 1024 * 1024 * 1024 * 1024) {
				filesize = Math.round(filesize / 1024 / 1024) + 'M'
			}
			// add separater to number
			return filesize.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",")
		},
		computed_video_show() {
			if (this.playing_file.id && this.playing_url) {
				return true
			}
			return false
		},
		computed_show() {
			return this.file.id ? true : false
		},
	},
}

const component_search_files = {
	emits: ['play_audio'],
	props: ['token'],
	computed: {
		computed_files_page() {
			if (this.is_loading) {
				return 'Loading...'
			}
			if (this.error_status) {
				return this.error_status
			}
			return this.offset + ' ~ ' + (this.offset + this.files.length)
		},
	},
	template: `
<div class="page">
<h3>Search Files</h3>
<div class="search_toolbar">
<input type="text" name="filename" @keyup.enter="first_search_files" v-model="search_filenames" placeholder="Enter filename" />
<button @click="first_search_files">Search</button>
<button @click="last_page">Last Page</button>
<button disabled>{{ computed_files_page }}</button>
<button @click="next_page">Next Page</button>
</div>
<table v-if="files.length">
	<thead>
		<tr>
			<th>Filename</th>
			<th>Folder Name</th>
			<th>Size</th>
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
			is_loading: false,
			error_status: "",
		}
	},
	methods: {
		first_search_files() {
			this.offset = 0
			this.search_files()
		},
		search_files() {
			this.is_loading = true
			axios.post('/api/v1/search_files', {
				filename: this.search_filenames,
				limit: this.limit,
				offset: this.offset,
			}).then((response) => {
				this.error_status = ""
				this.files = response.data.files
			}).catch((error) => {
				if (error.response) {
					this.error_status = error.response.data.status
				} else {
					this.error_status = 'Network error'
				}
			}).finally(() => {
				this.is_loading = false
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
	props: ['token'],
	data() {
		return {
			files: [],
			is_loading: false,
			error_status: "",
		}
	},
	computed: {
		computed_refresh() {
			if (this.error_status) {
				return this.error_status
			}
			if (this.is_loading) {
				return 'Loading...'
			}
			return 'Refresh'
		},
	},
	template: `
<div class="page">
<div class="search_toolbar">
	<button class="refresh" @click="get_random_files">{{ computed_refresh }}</button>
</div>
<table>
	<thead>
		<tr>
			<th>Filename</th>
			<th>Folder Name</th>
			<th>Size</th>
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
	mounted() {
		this.get_random_files()
	},
	methods: {
		get_random_files() {
			this.is_loading = true
			axios.get('/api/v1/get_random_files'
			).then(response => {
				this.error_status = ""
				this.files = response.data.files;
			}).catch((error) => {
				if (error.response) {
					this.error_status = error.response.data.status
				} else {
					this.error_status = 'Network error'
				}
			}).finally(() => {
				this.is_loading = false
			})
		}
	},
}

const component_stream_config = {
	emits: ['set_ffmpeg_config'],
	data() {
		return {
			ffmpeg_config_list: [],
			selected_ffmpeg_config: {},
		}
	},
	template: `
<div>
<table>
<tbody>
<tr>
<td>
<select v-model="selected_ffmpeg_config">
	<option v-for="ffmpeg_config in ffmpeg_config_list" :value="ffmpeg_config">
		{{ ffmpeg_config.name }}
	</option>
</select>
</td>
<td>
<span>{{ selected_ffmpeg_config.args }}</span>
</td>
</tr>
</tbody>
</table>
</div>
`,
	mounted() {
		axios.get('/api/v1/get_ffmpeg_config_list',
		).then(response => {
			// 后端返回数据 ffmpeg_configs 是一个字典，name 作为 key，ffmpeg_config{} 作为 value
			// 为方便前端，此处将 ffmpeg_configs 转为数组，并添加 name 到每个对象中
			var ffmpeg_configs = response.data.ffmpeg_configs
			var tmp_list = []
			for (var key in ffmpeg_configs) {
				var ffmpeg_config = ffmpeg_configs[key]
				ffmpeg_config.name = key
				tmp_list.push(ffmpeg_config)
			}
			tmp_list.sort()
			this.ffmpeg_config_list = tmp_list
			this.selected_ffmpeg_config = this.ffmpeg_config_list[0]
		}).catch(err => {
			this.ffmpeg_config_list = [{name: 'No avaliable config'}]
			this.selected_ffmpeg_config = this.ffmpeg_config_list[0]
		})
	},
	watch: {
		selected_ffmpeg_config(n, o) {
			this.$emit('set_ffmpeg_config', this.selected_ffmpeg_config)
		},
	},
}

const routes = [
	{ path: '/', component: component_get_random_files},
	{ path: '/search_files', component: component_search_files},
	{ path: '/search_folders', component: component_search_folders},
	{ path: '/manage', component: component_manage},
	{ path: '/share', component: component_share},
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
		stop() {
			this.playing_audio_file = {}
		},
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
app.component('component-stream-config', component_stream_config)
app.component('component-manage-database', component_manage_database)
app.component('component-share', component_share)

app.use(router)

app.mount('#app')
