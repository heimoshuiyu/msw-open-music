const component_search_folders = {
	emits: ['play_audio'],
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
<input type="text" v-model="search_foldernames" />
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
			<component-file :file=file @play_audio="$emit('play_audio', file)"></component-file>
		</tr>
	</tbody>
</table>
`,
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

const component_update_database = {
	data() {
		return {
			token: "",
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
	<td>Token</td>
	<td><input type="text" v-model="token" /></td>
</tr>
<tr>
	<td>Root</td>
	<td><input type="text" v-model="root" /></td>
</tr>
<tr>
	<td><button @click="add_pattern">Add Pattern</button></td>
	<td><input type="text" v-model="pattern_tmp" /></td>
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
		},
		emit_play_audio() {
			console.log("pressed button")
			this.file.play_back_type = 'raw'
			this.$emit("play_audio", this.file)
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
<td>{{ file.filename }}</td>
<td>{{ file.foldername }}</td>
<td>{{ computed_readable_size }}</td>
<td>
	<button @click="dialog">Dialog</button>
	<component-file-dialog
		@close_dialog="close_dialog"
		@play_audio="$emit('play_audio', this.file)"
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
		close_dialog() {
			this.show_dialog = false
		},
		dialog() {
			this.show_dialog = true
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
<video v-if="computed_show" :src="computed_playing_audio_file_url" controls autoplay>
</video>
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
<input type="text" name="filename" v-model="search_filenames" />
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
			<component-file :file=file @play_audio="$emit('play_audio', file)"></component-file>
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
	emits: ['play_audio'],
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
			<component-file :file=file @play_audio="$emit('play_audio', file)"></component-file>
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
	{ path: '/update_database', component: component_update_database},
]
const router = VueRouter.createRouter({
	history: VueRouter.createWebHashHistory(),
	routes,
})

const app = Vue.createApp({
	data() {
		return {
			playing_audio_file: {},
		}
	},
	methods: {
		play_audio(file) {
			console.log(file)
			this.playing_audio_file = file
		},
	},
})

app.component('component-search-folders', component_search_folders)
app.component('component-update-database', component_update_database)
app.component('component-file', component_file)
app.component('component-audio-player', component_audio_player)
app.component('component-search-files', component_search_files)
app.component('component-get-random-files', component_get_random_files)
app.component('component-file-dialog', component_file_dialog)

app.use(router)

app.mount('#app')
