import FileEntry from "./FileEntry";

function FilesTable(props) {
  if (props.files.length === 0) {
    return null;
  }
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
