const component_search_folders = {
	emits: ['play_audio'],
	data() {
		return {
			search_foldernames: "",
			folders: [],
			offset: 0,
			limit: 10,
			folder_offset: 0,
			folder_limit: 10,
			files_in_folder: [],
			playing_audio_file: {},
		}
	},
	template: `
<input type="text" v-model="search_foldernames" />
<button @click="first_search_folders">Search Folders</Button>

<button @click="last_page">Last Page</button>
<span>{{ offset }}~{{ offset + folders.length }}</span>
<button @click="next_page">Next Page</button>

<table border="1">
<thead>
<tr>
	<th>ID</th>
	<th>Folder Name</th>
	<th>Action</th>
</tr>
</thead>
<tbody>
<tr v-for="folder in folders">
	<td>{{ folder.id }}</td>
	<td>{{ folder.foldername }}</td>
	<td><button @click="get_files_in_folder(folder)">View</button></td>
</tr>
</tbody>
</table>

<table border="1">
	<thead>
		<tr>
			<th>ID</th>
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
		get_files_in_folder(folder) {
			axios.post('/api/v1/get_files_in_folder', {
				folder_id: folder.id,
				limit: this.folder_limit,
				offset: this.folder_offset,
			}).then((response) => {
				this.files_in_folder = response.data.files
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
<table border="1">
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
	<td><strong>Pattern List</strong></td>
</tr>
<tr v-for="p in pattern">
	<td>{{ p }}</td>
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

const component_file = {
	props: ['file'],
	emits: ['play_audio'],
	template: `
<td>{{ file.id }}</td>
<td>{{ file.filename }}</td>
<td>{{ file.foldername }}</td>
<td>{{ computed_readable_size }}</td>
<td>
	<button @click="download_file(file)" :disabled="disabled">{{ computed_download_status }}</button>
	<button @click="emit_play_audio">Play</button>
</td>
`,
	data() {
		return {
			download_loaded: 0,
			disabled: false,
		}
	},
	methods: {
		emit_play_audio() {
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
			return '/api/v1/get_file_direct?id=' + this.file.id
		},
		computed_show() {
			return this.file.id ? true : false
		},
	},
}

const component_search_files = {
	emits: ['play_audio'],
	template: `
<input type="text" name="filename" v-model="search_filenames" />
<button @click="first_search_files">Search</button>
<button @click="last_page">Last Page</button>
<span>{{ offset }}~{{ offset + files.length }}</span>
<button @click="next_page">Next Page</button>
<table border="1">
	<thead>
		<tr>
			<th>ID</th>
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

const routes = [
	{ path: '/', component: component_search_files },
	{ path: '/search_folders', component: component_search_folders, },
	{ path: '/update_database', component: component_update_database },
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
app.component('component-search_files', component_search_files)

app.use(router)

app.mount('#app')
