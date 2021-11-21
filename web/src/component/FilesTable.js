import FileEntry from "./FileEntry";

function FilesTable(props) {
  return (
    <table>
      <thead>
        <tr>
          <th>Filename</th>
          <th>Folder Name</th>
          <th>Size</th>
        </tr>
      </thead>
      <tbody>
        {props.files.map((file) => (
          <FileEntry
            setPlayingFile={props.setPlayingFile}
            key={file.id}
            file={file}
          />
        ))}
      </tbody>
    </table>
  );
}

export default FilesTable;
