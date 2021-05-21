const app = Vue.createApp({
	data() {
		return {
			search_filenames: '',
			download_total: 0,
			download_loaded: 0,
		}
	},
	methods: {
	},
})

app.component('component-file', {
	props: ['file'],
	template: `
<td>{{ file.id }}</td>
<td>{{ file.filename }}</td>
<td>{{ computed_readable_size }}</td>
<td><button @click="download_file(file)">{{ computed_download_status }}</button></td>
`,
	data() {
		return {
			download_loaded: 0,
		}
	},
	methods: {
		download_file(file) {
			axios({
				url: '/api/v1/get_file',
				method: 'POST',
				responseType: 'blob', // important
				data: {
					id: file.id,
				},
				onDownloadProgress: ProgressEvent => {
					console.log(ProgressEvent.loaded)
					this.download_loaded = ProgressEvent.loaded
				}
			}).then((response) => {
				const url = window.URL.createObjectURL(new Blob([response.data]));
				const link = document.createElement('a');
				link.href = url;
				link.setAttribute('download', file.filename);
				document.body.appendChild(link);
				link.click();
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
})

app.component('component-search-files', {
	template: `
<input type="text" name="filename" v-model="search_filenames" />
<input type="button" value="Search" @click="search_files(this)" />
<table border="1">
	<thead>
		<tr>
			<th>ID</th>
			<th>Filename</th>
			<th>Size</th>
			<th>Action</th>
		</tr>
	</thead>
	<tbody>
		<tr v-for="file in files">
			<component-file :file=file></component-file>
		</tr>
	</tbody>
</table>
`,
	data() {
		return {
			files: [],
		}
	},
	methods: {
		search_files(app) {
			axios.post('http://localhost:8080/api/v1/search_files', {
				filename: app.search_filenames,
				limit: 10,
				offset: 0,
			}).then(function(response) {
				app.files = response.data.files
			})
		},
	},
})
app.mount('#app')
